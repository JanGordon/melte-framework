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
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var watcher, err = fsnotify.NewWatcher()
var Server = StartServer()
var upgrader = websocket.Upgrader{}
var conn *websocket.Conn

func hotReloadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	fmt.Println("Upgrading to ws")
	defer conn.Close()

	// Continuosly read and write message
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		runS(conn, message, mt)
		fmt.Println(message)
		message = []byte("reload")
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write failed:", err)
			break
		}
	}
}
func Run(port int) {
	err := watcher.Add("./")
	err = watcher.Add("./components")
	err = watcher.Add("./hotReload")
	err = filepath.WalkDir("./hotReload", initWatcher)
	err = filepath.WalkDir("./routes", initWatcher)
	Server.GET("/hotReloadWS", hotReloadHandler)
	RunServer(Server)

	if err != nil {
		panic("error reading routes folder")
	}
	if err != nil {
		log.Fatal("Add failed:", err)
	}
}
func runS(conn *websocket.Conn, message []byte, mt int) {

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

				// fmt.Println(event.Name, event.Op)
				_, filename := filepath.Split(event.Name)
				if isSystemFile(filename) {
					// fmt.Println("Ignoring: ", filename)
				} else {
					reBuildFull()
					reload(conn, mt)

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

	<-done
	//compile.BuildPage(compile.ReplaceComponentWithHTML("test.html"), "out.html", false)
}

func reload(conn *websocket.Conn, mt int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	if conn != nil {
		conn.WriteMessage(mt, []byte("reload"))
		if err := recover(); err != nil {
			log.Println("write failed :", err)
		}
		//wait for reloaded message and if doest come within p.5s resend reload request
	}
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
	if !di.IsDir() {
		dir, filename := filepath.Split(path)
		if filename == "out.html" && !stringInSlice(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), routeHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "routes", "", 1)))
			fmt.Println("Handling ", filepath.Join("/", strings.Replace(dir, "routes", "", 1)), " ", path)

		} else if strings.HasPrefix(path, "routes") && !stringInSlice(filepath.Join("/", strings.Replace(path, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(path, "routes", "", 1)), fileInRouteHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(path, "routes", "", 1)))
			fmt.Println("Handling ", filepath.Join("/", strings.Replace(path, "routes", "", 1)), " ", path)

		} else if !stringInSlice(filepath.Join("/", path), Handlers) {
			Server.GET(filepath.Join("/", path), otherHandler)
			Handlers = append(Handlers, filepath.Join("/", path))
			fmt.Println("Handling ", filepath.Join("/", path), " ", path)

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

func fileInRouteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := filepath.Join("./routes", r.URL.Path)
	fmt.Println("serving", "./routes"+r.URL.Path)
	http.ServeFile(w, r, p)
}

func otherHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := filepath.Join("./", r.URL.Path)
	fmt.Println("serving", filepath.Join("./", r.URL.Path))
	http.ServeFile(w, r, p)
}

func routeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println(r.URL.Path)
	p := filepath.Join("./routes", r.URL.Path, "out.html")
	fmt.Println("serving", filepath.Join("./routes", r.URL.Path, "out.html"), " wiht routeHandler")
	http.ServeFile(w, r, p)
}

func goHandler() {
	//
}

func visitPath(path string, di fs.DirEntry, err error) error {
	//filename := filepath.Base(path)
	if di.IsDir() && !stringInSlice(path, watcher.WatchList()) {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal("Add failed:", err)
		}
	}
	if !di.IsDir() {
		dir, filename := filepath.Split(path)
		if filename == "out.html" && !stringInSlice(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), routeHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "routes", "", 1)))
			fmt.Println("Handling ", filepath.Join("/", strings.Replace(dir, "routes", "", 1)), " ", path)

		} else if strings.HasPrefix(path, "routes") && !stringInSlice(filepath.Join("/", strings.Replace(path, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(path, "routes", "", 1)), fileInRouteHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(path, "routes", "", 1)))
			fmt.Println("Handling ", filepath.Join("/", strings.Replace(path, "routes", "", 1)), " ", path)

		} else if !stringInSlice(filepath.Join("/", path), Handlers) {
			Server.GET(filepath.Join("/", path), otherHandler)
			Handlers = append(Handlers, filepath.Join("/", path))
			fmt.Println("Handling ", filepath.Join("/", path), " ", path)

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
