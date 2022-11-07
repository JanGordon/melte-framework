package dev

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/JanGordon/melte-framework/compile"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var cwd, _ = os.Getwd()
var watcher, err = fsnotify.NewWatcher()
var Server = StartServer()
var upgrader = websocket.Upgrader{}
var conn *websocket.Conn
var wConn = 0

func hotReloadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	fmt.Println("Upgrading to ws")
	wConn = 1
	if reloadToDo == 1 {
		listenForSuccess = 1
		reload(conn, 1)
	}

	defer conn.Close()
	defer func() { wConn = 0 }()
	// Continuosly read and write message
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		runS(conn, message, mt)
		message = []byte("reload")
		if listenForSuccess == 1 {
			if string(message) == "reloaded" {
				reloadToDo = 0
				fmt.Println("reload succeeded")
			} else {
				reloadToDo = 1
				fmt.Println("reload failed")
			}
		}
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write failed:", err)
			break
		}
	}
}
func Run(port int) {
	err := watcher.Add(cwd)
	// err = watcher.Add(cwd + "/components")
	// err = watcher.Add("./hotReload")
	err = filepath.WalkDir(cwd, initWatcher)
	// err = filepath.WalkDir(cwd+"/hotReload", initWatcher)
	// err = filepath.WalkDir(cwd+"/routes", initWatcher)
	Server.GET("/clientSideRouting/out.js", devHandler)
	Server.GET("/hotReload/WebSocket.js", devHandler)
	Server.GET("/hotReloadWS", hotReloadHandler)
	reBuildFull()
	fmt.Println("Starting dev server for ", cwd, " ...")
	fmt.Printf("connect to http://127.0.0.1:%v/ to begin live rebuild\n", port)

	RunServer(Server, port)

	if err != nil {
		panic("error reading routes folder")
	}
	if err != nil {
		log.Fatal("Add failed:", err)
	}
}

var reloadToDo = 0
var listenForSuccess = 0

func runS(conn *websocket.Conn, message []byte, mt int) {

	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	//server

	done := make(chan bool)
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
					listenForSuccess = 1
					err := reload(conn, mt)
					fmt.Println(err)
					// if err != "" {
					// 	reloadToDo = 1
					// 	fmt.Println("an error occured sending reload request")
					// }
					// mt, message, err := conn.ReadMessage()
					// if err != nil {
					// 	log.Fatal(err)
					// 	break
					// }
					// handleMessage(mt, string(message))

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

func handleMessage(mt int, message string) {
}

func reload(conn *websocket.Conn, mt int) string {
	var er = "none"
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			er = "e"

		}
	}()
	if conn != nil {
		conn.WriteMessage(mt, []byte("reload"))
		if err := recover(); err != nil {
			log.Println("write failed :", err)
		} else {
			return "none"
		}

		//wait for reloaded message and if doest come within p.5s resend reload request
	}
	return er
	// go func(conn *websocket.Conn) {
	// 	time.Sleep(3 * time.Second)
	// 	mt, message, err := conn.ReadMessage()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if string(message) != "reloaded" {
	// 		reload(conn, mt)
	// 	}
	// }(conn)
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
		dir = strings.Replace(dir, cwd, "", 1)
		//fmt.Println("dir : ", dir+filename, " has prefix : ", strings.HasPrefix(dir, "/routes"), " is in handlers : ", !stringInSlice(filepath.Join("/", strings.Replace(dir, "/routes", "", 1)), Handlers), " route : ", filepath.Join("/", strings.Replace(dir, "/routes", "", 1)))
		if filename == "out.html" && !stringInSlice(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), routeHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "routes", "", 1)))
			// fmt.Println("route Handling ", filepath.Join("/", strings.Replace(dir, "routes", "", 1)), " ", path)

		} else if strings.HasPrefix(dir, "/routes") && !stringInSlice(filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename), fileInRouteHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename))
			// fmt.Println("fileinRoute Handling ", filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename), " ", path)

		} else if !stringInSlice(filepath.Join("/r", dir, filename), Handlers) {
			Server.GET(filepath.Join("/r", dir, filename), otherHandler)
			Handlers = append(Handlers, filepath.Join("/r", dir, filename))
			// fmt.Println("other Handling ", filepath.Join("/r", dir, filename), " ", path)

		}
	}
	return nil
}

func reBuildFull() {
	err := filepath.WalkDir(cwd+"/routes", visitPath)
	if err != nil {
		panic("error reading routes folder")
	}

}

func fileInRouteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Println(r.URL.Path)
	p := filepath.Join(cwd, "/routes", r.URL.Path)
	fmt.Println("serving", "./routes"+r.URL.Path, " with file in route handler")
	http.ServeFile(w, r, p)
}

func otherHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Println(r.URL.Path)
	p := filepath.Join(cwd, strings.Replace(r.URL.Path, "/r", "", 1))
	fmt.Println("serving", filepath.Join("./", r.URL.Path), " with other handler")
	http.ServeFile(w, r, p)
}

func devHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Println(r.URL.Path)
	if r.URL.Path == "/clientSideRouting/out.js" {
		fmt.Fprint(w, router())
		// http.ServeFile(w, r, "./clientSideRouting/src.js")
	} else if r.URL.Path == "/hotReload/WebSocket.js" {
		fmt.Fprint(w, hotReload())
	} else {
		p := filepath.Join("./", r.URL.Path)
		fmt.Println("serving", filepath.Join("./", r.URL.Path), " with other handler")
		http.ServeFile(w, r, p)
	}

}

func routeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//fmt.Println(r.URL.Path)
	p := filepath.Join("./routes", r.URL.Path, "out.html")
	fmt.Println("serving", filepath.Join("./routes", r.URL.Path, "out.html"), " wiht routeHandler")
	//os .write
	f, err := os.OpenFile(filepath.Join(cwd, filepath.Join("./routes", r.URL.Path, "out.html")), os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(f, checkServerFuncs(filepath.Join("./routes", r.URL.Path, "out.html")))
	http.ServeFile(w, r, p)
}
func goHandler() {
	//
}

func visitPath(path string, di fs.DirEntry, err error) error {
	if di.IsDir() && !stringInSlice(path, watcher.WatchList()) {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal("Add failed:", err)
		}
	}
	if !di.IsDir() {
		dir, filename := filepath.Split(path)
		dir = strings.Replace(dir, cwd, "", 1)
		//fmt.Println("dir : ", dir+filename, " has prefix : ", strings.HasPrefix(dir, "/routes"), " is in handlers : ", !stringInSlice(filepath.Join("/", strings.Replace(dir, "/routes", "", 1)), Handlers), " route : ", filepath.Join("/", strings.Replace(dir, "/routes", "", 1)))
		if filename == "out.html" && !stringInSlice(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(dir, "routes", "", 1)), routeHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "routes", "", 1)))
			//fmt.Println("route Handling ", filepath.Join("/", strings.Replace(dir, "routes", "", 1)), " ", path)

		} else if strings.HasPrefix(dir, "/routes") && !stringInSlice(filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename), Handlers) {
			Server.GET(filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename), fileInRouteHandler)
			Handlers = append(Handlers, filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename))
			//fmt.Println("fileinRoute Handling ", filepath.Join("/", strings.Replace(dir, "/routes", "", 1), filename), " ", path)

		} else if !stringInSlice(filepath.Join("/r", dir, filename), Handlers) {
			Server.GET(filepath.Join("/r", dir, filename), otherHandler)
			Handlers = append(Handlers, filepath.Join("/r", dir, filename))
			//fmt.Println("other Handling ", filepath.Join("/r", dir, filename), " ", path)

		}
	}
	// make server better and make it work to host the html fil in th e folder if it is just a folder

	dir, filename := filepath.Split(path)
	if filepath.Ext(path) == ".html" && filename != "out.html" && !strings.HasPrefix(filename, "layout") {
		compile.BuildPage(compile.ReplaceComponentWithHTML(compile.ParseHTMLFragmentFromPath(path), false, dir+"out.html"), dir+"out.html", dir, false, true, true)
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

func checkServerFuncs(route string) string {
	f, err := os.ReadFile(route)
	file := string(f)
	if err != nil {
		panic("failed to read route")
	}
	for _, f := range compile.ServerFunctions {
		fmt.Println(f.Route, filepath.Join(cwd, route))
		if f.Route == filepath.Join(cwd, route) {
			file = f.F(route, file)
			fmt.Println("Hello")
		}
	}
	return file
}
