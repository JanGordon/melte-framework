package dev

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/JanGordon/melte/compile"
	"github.com/fsnotify/fsnotify"
	"github.com/julienschmidt/httprouter"
)

var watcher, err = fsnotify.NewWatcher()
var Server = StartServer()

func Run(port int) {

	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	//server

	done := make(chan bool)
	fmt.Println("Staring dev server... watching for file changes")
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				fmt.Println(event.Name, event.Op)
				_, filename := filepath.Split(event.Name)
				if isSystemFile(filename) {
					// fmt.Println("Ignoring: ", filename)
				} else {
					reBuildFull()
				}
				// dir, filename := filepath.Split(path)
				// if filepath.Ext(path) == ".html" && filename != "out.html" {
				// 	fmt.Println("Rebuilding", path)
				// 	compile.BuildPage(compile.ReplaceComponentWithHTML(path), dir+"out.html", false)
				// }
				//
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	err = watcher.Add("./")
	err = watcher.Add("./components")
	err := filepath.WalkDir("./routes", initWatcher)
	go RunServer(Server)
	if err != nil {
		panic("error reading routes folder")
	}
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
	//compile.BuildPage(compile.ReplaceComponentWithHTML("test.html"), "out.html", false)
}

func reBuildChunk(dir string) {

}

func initWatcher(path string, di fs.DirEntry, err error) error {
	if di.IsDir() && !stringInSlice(path, watcher.WatchList()) {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal("Add failed:", err)
		}
	}
	return nil
}

func reBuildFull() {
	err := filepath.WalkDir("./routes", visitPath)
	if err != nil {
		panic("error reading routes folder")
	}

}

func requestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := "./routes" + r.URL.Path
	fmt.Println("serving", "./routes"+r.URL.Path)
	http.ServeFile(w, r, p)
}

func otherHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := r.URL.Path
	fmt.Println("serving", r.URL.Path)
	http.ServeFile(w, r, p)
}

func routeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := "./routes" + r.URL.Path + "out.html"
	fmt.Println("serving", "./routes"+r.URL.Path+"out.html", " wiht routeHandler")
	http.ServeFile(w, r, p)
}

func fileRouteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := "./routes" + r.URL.Path + ".html"
	fmt.Println("serving", "./routes"+r.URL.Path+".html", " wiht routeHandler")
	http.ServeFile(w, r, p)
}

func visitPath(path string, di fs.DirEntry, err error) error {
	//filename := filepath.Base(path)
	if di.IsDir() && !stringInSlice(path, watcher.WatchList()) {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal("Add failed:", err)
		}
	}
	fmt.Println(Handlers)
	if !di.IsDir() {
		dir, _ := filepath.Split(path)
		if strings.HasPrefix(dir, "./components") {
			fmt.Println("Other Handler listening for ", filepath.Join("/", strings.Replace(dir, "routes", "", 1)))

			Server.GET(filepath.Join("/", dir), otherHandler)

			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "routes", "", 1)))

		} else if filepath.Base(filepath.Join("/", strings.Replace(path, "routes", "", 1))) == "out.html" && !stringInSlice(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), Handlers) {

			fmt.Println("Route Handler listening for ", filepath.Join("/", strings.Replace(dir, "routes", "", 1)))

			Server.GET(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), routeHandler)

			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "routes", "", 1)))

		} else if filepath.Ext(filepath.Base(filepath.Join("/", strings.Replace(path, "routes", "", 1)))) == ".html" && !stringInSlice(filepath.Join("/", strings.Replace(strings.Replace(path, "routes", "", 1), ".html", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(strings.Replace(path, "routes", "", 1), ".html", "", 1)), fileRouteHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(strings.Replace(path, "routes", "", 1), ".html", "", 1)))

			fmt.Println("File Route Handler listening for ", filepath.Join("/", strings.Replace(strings.Replace(path, "routes", "", 1), ".html", "", 1)))

		} else if !stringInSlice(filepath.Join("/", strings.Replace(path, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(path, "routes", "", 1)), requestHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(path, "routes", "", 1)))

			fmt.Println("File Handler listening for ", filepath.Join("/", strings.Replace(path, "routes", "", 1)))
		}
	}
	// make server better and make it work to host the html fil in th e folder if it is just a folder

	dir, filename := filepath.Split(path)
	if filepath.Ext(path) == ".html" && filename != "out.html" {
		fmt.Println("Rebuilding", path)
		compile.BuildPage(compile.ReplaceComponentWithHTML(path), dir+"out.html", dir, false, true)
	}
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isSystemFile(filename string) bool {
	if strings.HasPrefix(filename, "out") || filename == "in.ts" {
		return true
		//is sys file
	} else {
		return false
	}
}
