package explorer

import (
	"fmt"
	"github.com/Hyojip/housecoin/blockchain"
	"html/template"
	"log"
	"net/http"
)

const templateDir = "explorer/templates/"

var port string
var templates *template.Template

type homeData struct {
	Title  string
	Blocks []*blockchain.Block
}

func Start(aPort string) {
	port = aPort

	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/add", add)

	fmt.Printf("web page listen localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func add(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		templates.ExecuteTemplate(writer, "add", nil)
	case "POST":
		request.ParseForm()
		blockchain.GetBlockchain().AddBlock()
		http.Redirect(writer, request, "/", http.StatusPermanentRedirect)
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	parameters := homeData{"Home", blockchain.FindBlocks()}
	templates.ExecuteTemplate(writer, "home", parameters)
}
