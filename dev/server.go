package dev

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var Handlers []string

func StartServer() *httprouter.Router {
	router := httprouter.New()
	// router.GET("/", Index)
	// router.GET("/hello/:name", Hello)
	return router
}
func RunServer(router *httprouter.Router, port int) {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

// func Serve(path string) {
// 	http.HandleFunc(path, serveHTML)
// 	http.HandleFunc("/out.js", serveFiles)
// 	fmt.Println("Serving new :", filepath.Join("/", path))
// }

// func serveHTML(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.URL.Path)
// 	p := "./routes" + r.URL.Path + "/out.html"
// 	fmt.Println("serving", "./routes"+r.URL.Path+"out.html")
// 	http.ServeFile(w, r, p)
// }

// func serveFiles(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("r.URL.Path")
// 	p := "./routes" + r.URL.Path
// 	fmt.Println("serving", "./routes"+r.URL.Path)
// 	http.ServeFile(w, r, p)
// }
