package models

type Post struct {
	ID              string
	Title           string
	ContentHtml     string // надо чтобы отображать
	ContentMarkdown string // надо чтобы редактировать
}

func NewPost(id, title, contentHtml, ContentMarkdown string) *Post {
	return &Post{id, title, contentHtml, ContentMarkdown}
}
