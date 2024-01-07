package explorer

import (
	"fmt"
	"github.com/Hyojip/housecoin/blockchain"
	"html/template"
	"log"
	"net/http"
)

const (
	port        = ":8080"
	templateDir = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	Title  string
	Blocks []*blockchain.Block
}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("web page listen localhost:%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func add(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		templates.ExecuteTemplate(writer, "add", nil)
	case "POST":
		request.ParseForm()
		blockchain.GetBlockchain().AddBlock(request.Form.Get("blockData"))
		http.Redirect(writer, request, "/", http.StatusPermanentRedirect)
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	parameters := homeData{"Home", blockchain.AllBlock()}
	templates.ExecuteTemplate(writer, "home", parameters)
}
