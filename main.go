package main

import (
	"errors"
	"fmt"
	"github.com/mehmetron/pid1/executor"
	"github.com/mehmetron/pid1/userPort"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	res := executor.Handler()
	if res == "" {
		fmt.Fprintf(w, "Command was killed")
		return
	}
	fmt.Fprintf(w, res)
}

func getOpenPort(w http.ResponseWriter, r *http.Request) {
	res := userPort.GetOpenPort()
	fmt.Fprintf(w, fmt.Sprintf("Output: %s \n", res))
}

func getOpenPortLib(w http.ResponseWriter, r *http.Request) {
	demoPort := userPort.GetOpenPortLib()
	fmt.Fprintf(w, fmt.Sprintf("done i think %d", demoPort))
}

type Result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	OpenedPort uint16 `json:"opened_port"`
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in executeHandler")
	var command executor.Command
	err := decodeJSONBody(w, r, &command)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	//executor.Command{run, "main.go", content}
	fmt.Println("command 51: ", command)
	stdout, stderr, demoPort, err := command.Running()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("stdout: ", stdout)
	fmt.Println("stderr: ", stderr)

	//createdLesson := env.lessons.CreateLesson(lesson)
	err = writeJSON(w, http.StatusOK, Result{stdout, stderr, demoPort}, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("hello world")
	//err := InstallPackage("github.com/ermos/gomon")
	//if err != nil {
	//    fmt.Println(err)
	//}

	//go PipingStd()
	//WatchFiles()

	conf := LoadConfiguration("langs.json")
	fmt.Println("this is it ", conf)

	mux := http.NewServeMux()
	mux.HandleFunc("/executeHandler", executeHandler)
	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "first, %s!", r.URL.Path[1:])
	//})

	mux.HandleFunc("/bobhandler", bobHandler)
	mux.HandleFunc("/bob", handler)
	mux.HandleFunc("/ws", cmdToResponse)
	mux.HandleFunc("/page", serveHtml)
	mux.HandleFunc("/ports", getOpenPort)
	mux.HandleFunc("/portslib", getOpenPortLib)
	err := http.ListenAndServe(":8080", CorsHeaders(mux))
	fmt.Println("running main server")
	if err != nil {
		fmt.Printf("97 error %s\n", err)
	}

}
