package model

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Article struct {
	ID        uuid.UUID          `bson:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
	Tags      []string           `bson:"tags" json:"tags"`
}

func (a *Article) HasTag(tag string) bool {
	for _, v := range a.Tags {
		if v == tag {
			return true
		}
	}
	return false
}

func FindArticle(uuid uuid.UUID) *Article {
	var article Article
	err := database.Collection("article").FindOne(ctx, bson.M{"id": uuid}).Decode(&article)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &article
}

func CreateArticle(article *Article) *mongo.InsertOneResult {
	article.ID = uuid.New()
	article.CreatedAt = primitive.DateTime(time.Now().UnixMilli())
	article.UpdatedAt = primitive.DateTime(time.Now().UnixMilli())
	result, err := database.Collection("article").InsertOne(ctx, article)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return result
}

func UpdateArticle(article *Article) *mongo.UpdateResult {
	article.UpdatedAt = primitive.DateTime(time.Now().UnixMilli())
	result, err := database.Collection("article").UpdateOne(
		ctx, bson.M{"id": article.ID}, bson.D{{"$set", article}},
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return result
}

func ArticlesByQuery(query interface{}) *[]Article {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})
	result, err := database.Collection("article").Find(ctx, query, findOptions)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	var articles []Article
	err = result.All(ctx, &articles)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &articles
}

func ArticlesByTag(tag string) *[]Article {
	return ArticlesByQuery(bson.M{"tags": tag})
}

func AllArticles() *[]Article {
	return ArticlesByQuery(bson.M{})
}
