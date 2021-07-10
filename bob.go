package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func createBob(code string) {
	f, err := os.Create("bob.go")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal(err)
	}
}
func runBob(code string) {

	createBob(code)

	cmd := exec.Command("go", "run", "bob.go")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

type Execute struct {
	Code string `json:"code,omitempty"`
}

func bobHandler(w http.ResponseWriter, r *http.Request) {

	var Execute Execute
	err := decodeJSONBody(w, r, &Execute)
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

	runBob(Execute.Code)

	//createdLesson := env.lessons.CreateLesson(lesson)
	//err = writeJSON(w, http.StatusOK, createdLesson, nil)
	//if err != nil {
	//	fmt.Println(err)
	//}
}
