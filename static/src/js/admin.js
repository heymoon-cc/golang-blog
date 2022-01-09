const showdown  = require('showdown');
const converter = new showdown.Converter();
/** @var article HTMLTextAreaElement */
/** @var preview HTMLDivElement */
/** @var save HTMLFormElement */
/** @var title HTMLInputElement */
const previewMd = () => {
  preview.innerHTML = converter.makeHtml(article.value);
}
article.addEventListener('input', () => {
  previewMd();
});
previewMd();

const tags = new Set;

save.addEventListener('submit', (e) => {
  e.preventDefault();
  fetch('', {
    method: 'post',
    body: JSON.stringify({
      title: title.value,
      content: article.value,
      tags: Array.from(tags)
    })
  }).then((response) => {
    location.href = response.url;
  });
});

const selectTag = (e) => {
  e.preventDefault();
  const target = e.currentTarget;
  if (target.classList.contains('active')) {
    target.classList.remove('active');
    tags.delete(target.dataset.tag);
  } else {
    target.classList.add('active');
    tags.add(target.dataset.tag);
  }
}

document.querySelectorAll('[data-tag]').forEach((nav) => {
  if (nav.classList.contains('active')) {
    tags.add(nav.dataset.tag);
  }
  nav.addEventListener('click', selectTag);
});
