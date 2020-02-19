package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/lithium555/SortMP3/web_app/go-blog-example/models"
)

var (
	posts map[string]*models.Post
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error()) // вернем ошибку прямо в браузер
	}

	fmt.Println(posts)

	t.ExecuteTemplate(w, "index", posts) // in index.html:   {{ range $key, $value := . }}   точка этотекущий контекст а текущий контекст это posts map[string]*models.Post
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "write", nil)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := models.NewPost(id, title, content)
	posts[post.ID] = post

	http.Redirect(w, r, "/", 302)
}

func main() {
	fmt.Println("Listening on port :3000")

	posts = make(map[string]*models.Post, 0)
	//
	//	http.Handle("/assets", http.FileServer(http.Dir("./assets")))   =====>>>>>> example.com/assets/css/app.css   - ищем файл app.css по такому пути: 'assets/css/app.css'
	// Нам так не надо поэтому делаем StripPrefix()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler) // header.html есть ссылка на <li><a href="/write">New Post</a></li>

	// В write.html пост запрос на /SavePost: <form role="form" method="POST" action="/SavePost">
	http.HandleFunc("/SavePost", savePostHandler)

	http.ListenAndServe(":3000", nil)
}
