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

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	id := r.FormValue("id") // считываем айдии ищем пост(сообщение) в нашей мапе with key "id"
	post, found := posts[id]
	if !found {
		http.NotFound(w, r)
	}

	t.ExecuteTemplate(w, "write", post)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = GenerateID()
		post := models.NewPost(id, title, content)
		posts[post.ID] = post
	}

	http.Redirect(w, r, "/", 302) // TODO: надо почитать об этом.
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.NotFound(w, r)
	}

	delete(posts, id)

	http.Redirect(w, r, "/", 302) // TODO: надо почитать об этом.
}

func main() {
	fmt.Println("Listening on port :3000")

	posts = make(map[string]*models.Post)
	//
	//	http.Handle("/assets", http.FileServer(http.Dir("./assets")))   =====>>>>>> example.com/assets/css/app.css   - ищем файл app.css по такому пути: 'assets/css/app.css'
	// Нам так не надо поэтому делаем StripPrefix()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler) // header.html есть ссылка на <li><a href="/write">New Post</a></li>
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)

	// В write.html пост запрос на /SavePost: <form role="form" method="POST" action="/SavePost">
	http.HandleFunc("/SavePost", savePostHandler)

	http.ListenAndServe(":3000", nil)
}
