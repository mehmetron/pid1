package handlers

import (
	"bufio"
	"fmt"
	"github.com/mehmetron/pid1/executor"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var html = []byte(
	`<html>
	<body>
		<h1>blah.exe</h1>
		<code></code>
		<script>
			var ws = new WebSocket("ws://pid1.localhost:8080/ws?lang=vitejs")
			ws.onmessage = function(e) {
				document.querySelector("code").innerHTML += e.data + "<br>"
			}
		</script>
	</body>
</html>
`)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	defer func() {
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		fmt.Println("scan: ", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

//func main() {
//	http.HandleFunc("/ws", CmdToResponse)
//	http.HandleFunc("/", ServeHtml)
//	// http.HandleFunc("/bob", handler)
//
//	log.Println("Listening on :8000")
//	err := http.ListenAndServe(":8000", nil)
//	if err != nil {
//		fmt.Println("82 ", err)
//		// log.Fatal(err)
//	}
//}

//
//func handler(w http.ResponseWriter, r *http.Request) {
//
//	cmd := exec.Command("/bin/sh", "-c", "ls && sleep 5 && cd test1 && ls")
//	cmd.Dir = "/Users/mehmetcureoglu/go/src"
//
//	pr, pw := io.Pipe()
//	defer pw.Close()
//
//	cmd.Stdout = pw
//	cmd.Stderr = pw
//	go io.Copy(w, pr)
//
//	cmd.Run()
//}

//func reverseProxy() {
//	fmt.Println("in reverseProxy")
//	router2 := http.NewServeMux()
//
//	router2.HandleFunc("/", helpers.ReverseProxy)
//	srv := &http.Server{
//		Handler:      router2,
//		Addr:         ":8090",
//	}
//	log.Fatal(srv.ListenAndServe())
//	//log.Fatal(http.ListenAndServe(":8081", router2))
//}

func CmdToResponse(w http.ResponseWriter, r *http.Request) {

	langs, ok := r.URL.Query()["lang"]
	if !ok || len(langs[0]) < 1 {
		log.Println("Url Param 'lang' is missing")
		return
	}
	// Query()["lang"] will return an array of items,
	// we only want the single item.
	lang := langs[0]
	log.Println("Url Param 'lang' is: " + string(lang))

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("130 ", err)
		w.Write([]byte(fmt.Sprintf("", err)))
		return
	}
	defer ws.Close()

	fmt.Println("135 here is it")

	// discard received messages
	go func(c *websocket.Conn) {
		for {
			if _, _, err := c.NextReader(); err != nil {
				fmt.Println("140 ", err)
				c.Close()
				break
			}
		}
	}(ws)

	fmt.Println("149 here is it")

	ws.WriteMessage(1, []byte("Starting...\n"))

	cmd := executor.ExecuteRuntime(lang)

	fmt.Println("156 here is it")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("127 ", err)
		return
	}

	fmt.Println("164 here is it")

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("132 ", err)
		return
	}

	fmt.Println("172 here is it")

	if err := cmd.Start(); err != nil {
		fmt.Println("137 ", err)
		return
	}

	fmt.Println("179 here is it")

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, io.MultiReader(stdout, stderr), stdoutDone)

	fmt.Println("184 here is it")

	//s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	//for s.Scan() {
	//	ws.SetWriteDeadline(time.Now().Add(writeWait))
	//	err := ws.WriteMessage(websocket.TextMessage, s.Bytes())
	//	if err != nil {
	//		ws.Close()
	//		return
	//	}
	//}

	if err := cmd.Wait(); err != nil {
		fmt.Println("155  ", err)
		return
	}

	fmt.Println("201 here is it")

	ws.WriteMessage(1, []byte("Finished\n"))
}

func ServeHtml(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write(html)
}
