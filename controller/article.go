package controller

import (
  "encoding/json"
  "fmt"
  "github.com/gomarkdown/markdown"
  "github.com/gomarkdown/markdown/ast"
  "github.com/gomarkdown/markdown/html"
  "github.com/gomarkdown/markdown/parser"
  "github.com/google/uuid"
  "github.com/gorilla/mux"
  "html/template"
  "image"
  _ "image/jpeg"
  _ "image/png"
  "io"
  "main/model"
  "net/http"
  "os"
  "strconv"
  "strings"
  "time"
)

type ArticleView struct {
  ID         string
  Title      string
  Content    template.HTML
  CreatedAt  time.Time
  Draft      bool
  Authorized bool
}

type Markdown struct {
  Renderer       *html.Renderer
  Parser         *parser.Parser
  ImageAttribute *ast.Attribute
}

var images = make(map[string]*image.Config)

type IFrame struct {
  ast.Image
  Destination []byte
  Title       []byte
}

func (m *Markdown) RenderHookImage(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
  if entering {
    p, ok := node.(*ast.Paragraph)
    if ok {
      if p.Attribute != nil {
        found := false
        for _, c := range p.Children {
          img, ok := c.(*ast.Image)
          if !ok {
            continue
          }
          if img.Attribute == nil {
            found = true
            *&img.Attribute = *&p.Attribute
          }
        }
        if found {
          p.Attribute = nil
          m.Renderer.Paragraph(w, p, true)
          return ast.GoToNext, true
        }
      }
      return ast.GoToNext, false
    }
  }
  iframe, ok := node.(*IFrame)
  if ok {
    if entering {
      m.Renderer.Outs(w, `<iframe src="`+string(iframe.Destination)+`" height="auto" width="100%"`)
      if len(iframe.Title) > 0 {
        m.Renderer.Outs(w, ` title="`+string(iframe.Title)+`"`)
      }
      m.Renderer.Outs(w, `>`)
    } else {
      m.Renderer.Outs(w, `</iframe>`)
    }
    return ast.GoToNext, true
  }
  img, ok := node.(*ast.Image)
  if !ok {
    return ast.GoToNext, false
  }
  if rune(img.Destination[0]) != '/' {
    return ast.GoToNext, false
  }
  src := "." + string(img.Destination)
  im, found := images[src]
  if !found {
    reader, err := os.Open(src)
    if err == nil {
      defer func(reader *os.File) {
        err := reader.Close()
        if err != nil {
          fmt.Println(err.Error())
        }
      }(reader)
      config, _, err := image.DecodeConfig(reader)
      if err != nil {
        fmt.Println(err.Error())
        return ast.GoToNext, false
      }
      im = &config
      images[src] = im
    }
  }
  if im != nil {
    if !entering {
      iw := strconv.Itoa(im.Width)
      ih := strconv.Itoa(im.Height)
      if img.Attribute == nil {
        *&img.Attribute = *&m.ImageAttribute
      }
      img.Attrs[`height`] = []byte(ih)
      img.Attrs[`width`] = []byte(iw)
      _, exists := img.Attrs[`loading`]
      if !exists {
        img.Attrs[`loading`] = []byte(`lazy`)
      }
      attributes := strings.TrimSuffix(
        strings.Join(html.BlockAttrs(img), " "), `"`)
      m.Renderer.Outs(w, `" `+attributes)
    }
    m.Renderer.Image(w, img, entering)
    return ast.GoToNext, true
  }
  return ast.GoToNext, false
}

var ParserOptions = parser.Options{
  ParserHook: func(data []byte) (ast.Node, []byte, int) {
    i := 0
    if len(data) < 3 {
      return nil, data, 0
    }
    if data[i] != '@' {
      return nil, data, 0
    }
    i += 1
    var title = make([]byte, 0)
    for true {
      if data[i] == '(' {
        break
      }
      i++
      if i > len(data) {
        return nil, data, 0
      }
    }
    if i > 1 {
      title = data[1:i]
    }
    i += 1
    for true {
      if data[i] == ')' {
        break
      }
      i++
      if i > len(data) {
        return nil, data, 0
      }
    }
    linkStart := 2
    if len(title) > 0 {
      linkStart += len(title)
    }
    link := data[linkStart:i]
    return &IFrame{Destination: link, Title: title}, make([]byte, 0), i + 1
  },
}

func InitializeMarkdown() *Markdown {
  var m = new(Markdown)
  var renderHookImage = func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
    return m.RenderHookImage(w, node, entering)
  }
  m.Renderer = html.NewRenderer(html.RendererOptions{
    Flags:          html.CommonFlags | html.HrefTargetBlank | html.FootnoteReturnLinks,
    RenderNodeHook: renderHookImage,
  })
  m.Parser = parser.NewWithExtensions(parser.CommonExtensions | parser.Attributes | parser.Footnotes |
    parser.DefinitionLists | parser.Mmark | parser.Strikethrough | parser.SuperSubscript | parser.Tables)
  m.ImageAttribute = &ast.Attribute{
    Attrs:   make(map[string][]byte),
    Classes: append(make([][]byte, 0), []byte(`img-fluid`), []byte(`rounded`), []byte(`mx-auto`), []byte(`d-block`)),
  }
  m.Parser.Opts = ParserOptions
  return m
}

func createArticleView(article *model.Article, authorized bool) ArticleView {
  m := InitializeMarkdown()
  content := markdown.ToHTML([]byte(article.Content), m.Parser, m.Renderer)
  return ArticleView{
    article.ID.String(),
    article.Title,
    template.HTML(content),
    article.CreatedAt.Time(),
    article.Draft,
    authorized}
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
  id, err := uuid.Parse(mux.Vars(r)["id"])
  if err != nil {
    fmt.Println(err.Error())
    NotFoundHandler(w, r)
    return
  }
  article := model.FindArticle(id)
  if article == nil {
    NotFoundHandler(w, r)
    return
  }
  err = handleAuth(w, r, false)
  renderTemplate("./ui/html/article.page.tmpl", w, createArticleView(article, err == nil), "article")
}

func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
  err := handleAuth(w, r, true)
  if err != nil {
    return
  }
  if r.Method == "GET" {
    renderTemplate("./ui/html/create.page.tmpl", w, nil, "create")
    return
  }
  decoder := json.NewDecoder(r.Body)
  var article model.Article
  decoder.DisallowUnknownFields()
  err = decoder.Decode(&article)
  if err != nil {
    return
  }
  model.CreateArticle(&article)
  http.Redirect(w, r, fmt.Sprintf("/article/%s", article.ID.String()), http.StatusFound)
}

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
  err := handleAuth(w, r, true)
  if err != nil {
    return
  }
  id, err := uuid.Parse(mux.Vars(r)["id"])
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  article := model.FindArticle(id)
  if r.Method == "GET" {
    renderTemplate("./ui/html/create.page.tmpl", w, article, "update")
    return
  }
  decoder := json.NewDecoder(r.Body)
  decoder.DisallowUnknownFields()
  err = decoder.Decode(article)
  if err != nil {
    fmt.Println(err)
    return
  }
  model.UpdateArticle(article)
  http.Redirect(w, r, fmt.Sprintf("/article/%s", article.ID.String()), http.StatusFound)
}
