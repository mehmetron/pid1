package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os/exec"
)

var html = []byte(
	`<html>
	<body>
		<h1>blah.exe</h1>
		<code></code>
		<script>
			var ws = new WebSocket("ws://localhost:8080/ws")
			ws.onmessage = function(e) {
				document.querySelector("code").innerHTML += e.data + "<br>"
			}
		</script>
	</body>
</html>
`)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//func main() {
//	http.HandleFunc("/ws", cmdToResponse)
//	http.HandleFunc("/", serveHtml)
//	http.HandleFunc("/bob", handler)
//
//	log.Println("Listening on :8000")
//	err := http.ListenAndServe(":8000", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

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

func cmdToResponse(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("", err)))
		return
	}
	defer ws.Close()

	// discard received messages
	go func(c *websocket.Conn) {
		for {
			if _, _, err := c.NextReader(); err != nil {
				c.Close()
				break
			}
		}
	}(ws)

	ws.WriteMessage(1, []byte("Starting...\n"))

	// execute and get a pipe
	//cmd := exec.Command("blah.exe")
	cmd := exec.Command("/bin/sh", "-c", "ls && sleep 5 && ls")
	//cmd.Dir = "/Users/mehmetcureoglu/go/src"
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		ws.WriteMessage(1, s.Bytes())
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return
	}

	ws.WriteMessage(1, []byte("Finished\n"))
}

func serveHtml(w http.ResponseWriter, r *http.Request) {
	w.Write(html)
}
