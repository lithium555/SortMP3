package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/lithium555/SortMP3/web_app/go-blog-example/models"
	"github.com/martini-contrib/render" // это middleware которые упрощают нам жизнь
)

var (
	posts   map[string]*models.Post
	counter int
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error()) // вернем ошибку прямо в браузер
	}

	fmt.Println(posts)
	fmt.Printf("counter = '%v'\n", counter)

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

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	post, found := posts[id]
	if !found {
		rnd.Redirect("/")
		return
	}

	rnd.HTML(200, "write", post)
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
		return
	}

	delete(posts, id)

	rnd.Redirect("/")
}

func main() {
	fmt.Println("Listening on port :3000")

	posts = make(map[string]*models.Post, 0)
	counter = 0

	m := martini.Classic() // объект в котрый включено логгирование обработка, статических файлов и так далее/ Совместим с интерфесом http

	//m.Use(func(r *http.Request){   мы можем на кадждом запросе что-то отфильтровать по запросу и что-либо сделать
	//	if r.URL.Path == "/write"{
	//		counter++
	//	}
	//})

	// 	m.Get("/test", func() string {
	//		return "test" // мартини просто видит что возвращается строк аи он выведет ее на респонс
	//	})

	m.Use(render.Renderer(render.Options{
		Directory:       "templates",                    // Specify what path to load the templates from.
		Layout:          "layout",                       // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions:      []string{".tmpl", ".html"},     // Specify extensions to load for templates.
		Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		Delims:          render.Delims{"{[{", "}]}"},    // Sets delimiters to the specified strings.
		Charset:         "UTF-8",                        // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON:      true,                           // Output human readable JSON
		IndentXML:       true,                           // Output human readable XML
		HTMLContentType: "application/xhtml+xml",        // Output XHTML content type instead of default "text/html"
	}))

	//
	//	http.Handle("/assets", http.FileServer(http.Dir("./assets")))   =====>>>>>> example.com/assets/css/app.css   - ищем файл app.css по такому пути: 'assets/css/app.css'
	// Нам так не надо поэтому делаем StripPrefix()
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	staticOptions := martini.StaticOptions{Prefix: "asserts"}
	m.Use(martini.Static("assets", staticOptions)) // what about statuc files. In martini we have jandler
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler) // header.html есть ссылка на <li><a href="/write">New Post</a></li>

	m.Get("/edit", editHandler)
	m.Get("/delete", deleteHandler)
	// В write.html пост запрос на /SavePost: <form role="form" method="POST" action="/SavePost">
	m.Post("/SavePost", savePostHandler)

	//http.ListenAndServe(":3000", nil)
	m.Run()
}
