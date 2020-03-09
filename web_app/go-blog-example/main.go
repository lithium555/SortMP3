package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/lithium555/SortMP3/web_app/go-blog-example/models"
	"github.com/martini-contrib/render" // middleware
	"github.com/russross/blackfriday"
)

var (
	posts   map[string]*models.Post
	counter int
)

func indexHandler(rnd render.Render) {
	fmt.Printf("counter: '%v'\n", counter)

	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML(200, "write", nil) // in write we dont need any object, that is why we send nil
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"] // считываем айди, ищем пост(сообщение) в нашей мапе with key "id"
	post, found := posts[id]
	if !found {
		rnd.Redirect("/")
		return
	}

	rnd.HTML(200, "index", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := string(blackfriday.Run([]byte(contentMarkdown)))

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.ContentHtml = contentHtml
		post.ContentMarkdown = contentMarkdown
	} else {
		id = GenerateID()
		post := models.NewPost(id, title, contentHtml, contentMarkdown)
		posts[post.ID] = post
	}

	rnd.Redirect("/")
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

// main on http:
//func main() {
//	fmt.Println("Listening on port :3000")
//
//	posts = make(map[string]*models.Post)
//	//
//	//	http.Handle("/assets", http.FileServer(http.Dir("./assets")))   =====>>>>>> example.com/assets/css/app.css   - ищем файл app.css по такому пути: 'assets/css/app.css'
//	// Нам так не надо поэтому делаем StripPrefix()
//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
//	http.HandleFunc("/", indexHandler)
//	http.HandleFunc("/write", writeHandler) // header.html есть ссылка на <li><a href="/write">New Post</a></li>
//	http.HandleFunc("/edit", editHandler)
//	http.HandleFunc("/delete", deleteHandler)
//
//	// В write.html пост запрос на /SavePost: <form role="form" method="POST" action="/SavePost">
//	http.HandleFunc("/SavePost", savePostHandler)
//
//	http.ListenAndServe(":3000", nil)
//}

func getHtmlHandler(rnd render.Render, r http.Request) {
	md := r.FormValue("md")
	htmlBytes := blackfriday.Run([]byte(md))

	rnd.JSON(200, map[string]interface{}{"html": string(htmlBytes)})
}

func unescape(x string) interface{} { // если не будем юзать то HTML будет показываться просто тегами
	return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :3000")

	posts = make(map[string]*models.Post)
	counter = 0

	m := martini.Classic() // include logging, validation of statistics files

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{ // в тех хендлерах, в которых мы пропишем rnd render.Render, он будет автоматически инджектится
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // specify helper function maps fortemplatesto access
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	// For statistic files in martini we have handler martini.Static():
	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)   // header.html есть ссылка на <li><a href="/write">New Post</a></li>
	m.Get("/edit/:id", editHandler) // роутинг это передаа айди в урле. Вместо id моно передавать любой параметр и вместо него любой текст помто можно прочитать.
	m.Get("/delete/:id", deleteHandler)
	// В write.html пост запрос на /SavePost: <form role="form" method="POST" action="/SavePost">
	m.Post("/SavePost", savePostHandler)
	m.Post("/gethtml", getHtmlHandler)

	//http.ListenAndServe(":3000", nil)
	m.Run() // he will listen on port 3000
}
