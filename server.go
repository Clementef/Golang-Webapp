package main
import (
	"fmt"
	"io"
	"strings"
	"time"
	"log"
	"net/http"
	"html/template"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

type Context struct {
    Title  string
    Static string
}

func Home(w http.ResponseWriter, req *http.Request) {
    context := Context{Title: "Home"}
    render(w, "index", context)
}

func About(w http.ResponseWriter, req *http.Request) {
    context := Context{Title: "About"}
    render(w, "about", context)
}

func Carts(w http.ResponseWriter, req *http.Request) {
    context := Context{Title: "Carts"}
    render(w, "carts", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
    context.Static = STATIC_URL
    tmpl_list := []string{"templates/base.html",
        fmt.Sprintf("templates/%s.html", tmpl)}
    t, err := template.ParseFiles(tmpl_list...)
    if err != nil {
        log.Print("template parsing error: ", err)
    }
    err = t.Execute(w, context)
    if err != nil {
        log.Print("template executing error: ", err)
    }
    fmt.Println("Serving Page: ", strings.Title(tmpl))
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
    static_file := req.URL.Path[len(STATIC_URL):]
    if len(static_file) != 0 {
        f, err := http.Dir(STATIC_ROOT).Open(static_file)
        if err == nil {
            content := io.ReadSeeker(f)
            http.ServeContent(w, req, static_file, time.Now(), content)
            return
        }
    }
    http.NotFound(w, req)
}

func main() {
    http.HandleFunc("/", Home)
    http.HandleFunc("/about/", About)
    http.HandleFunc("/carts/", Carts)
    http.HandleFunc(STATIC_URL, StaticHandler)

	fmt.Println("Starting Server")
	err:=http.ListenAndServe(":8000",nil)
	if err != nil {
		log.Fatal("Listen and Serve: ", err)
	}
	fmt.Println("Program Terminating")
}