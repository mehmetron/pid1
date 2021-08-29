package main

import (
	"fmt"
	"github.com/mehmetron/pid1/handlers"
	"github.com/mehmetron/pid1/helpers"
	"github.com/mehmetron/pid1/userPort"
	"log"
	"net/http"
	"time"
)

func pid1() {
	fmt.Println("hello world second")

	mux := http.NewServeMux()
	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/", static)
	//mux.HandleFunc("/", handlers.ServeHtml)
	mux.HandleFunc("/ws", handlers.CmdToResponse)

	mux.HandleFunc("/kill", handlers.Kill)
	mux.HandleFunc("/port", func(w http.ResponseWriter, r *http.Request) {
		demoPort := userPort.GetOpenPortLib()

		err := helpers.WriteJSON(w, http.StatusOK, demoPort, nil)
		if err != nil {
			fmt.Println("77 ", err)
		}
	})

	mux.HandleFunc("/execute", handlers.SetupProject)
	mux.HandleFunc("/update", handlers.UpdateFilesVite)

	mux.HandleFunc("/websocketcommands", handlers.WebsocketCommands)

	port := ":8100"
	err := http.ListenAndServe(port, helpers.CorsHeaders(mux))
	fmt.Printf("Running main server on port %s\n", port)
	if err != nil {
		fmt.Printf("97 error %s\n", err)
	}
}

func main() {

	go pid1()

	router := http.NewServeMux()
	router.HandleFunc("/", handlers.ReverseProxy)
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 35 * time.Second,
		ReadTimeout:  35 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
