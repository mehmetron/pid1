package main

import (
	"errors"
	"fmt"
	"github.com/mehmetron/pid1/executor"
	"github.com/mehmetron/pid1/filemanager"
	"github.com/mehmetron/pid1/helpers"
	"github.com/mehmetron/pid1/packagemanager"
	"log"
	"net/http"
	"time"
)

type Result struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func execute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("13 execute handler")

	var command filemanager.Command
	err := helpers.DecodeJSONBody(w, r, &command)
	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Setup files
	err = filemanager.Create(command.Files)
	if err != nil {
		fmt.Println("48 ", err)
	}

	// Installs packages based on .hackr config file
	//args := []string{"faker 5.5.3"}
	//for i, s := range command.Packages {
	//fmt.Println("44 ", s, i)
	if command.Env == "js" || command.Env == "vitejs" {
		packagemanager.RunAdd("nodejs", command.Packages, true, "playground")
	} else if command.Env == "go" {
		packagemanager.RunAdd("go", command.Packages, true, "playground")
	} else if command.Env == "python" {
		packagemanager.RunAdd("python", command.Packages, true, "playground")
	}
	//}

	stdout, stderr := "port", ""
	// Execute code
	//go func() {
	//	//executor.CmdToResponse(w, r)
	//	stdout, stderr, err = executor.Running(command.Env)
	//	if err != nil {
	//		fmt.Println("38 ", err)
	//	}
	//}()

	time.Sleep(8 * time.Second)

	err = helpers.WriteJSON(w, http.StatusOK, Result{stdout, stderr}, nil)
	if err != nil {
		fmt.Println("77 ", err)
	}
}

func main() {
	fmt.Println("hello world")

	mux := http.NewServeMux()
	mux.HandleFunc("/", executor.ServeHtml)
	mux.HandleFunc("/ws", executor.CmdToResponse)

	mux.HandleFunc("/kill", executor.Kill)

	mux.HandleFunc("/execute", execute)

	port := ":8080"
	err := http.ListenAndServe(port, helpers.CorsHeaders(mux))
	fmt.Printf("Running main server on port %s\n", port)
	if err != nil {
		fmt.Printf("97 error %s\n", err)
	}

}
