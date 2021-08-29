package handlers

import (
	"errors"
	"fmt"
	"github.com/mehmetron/pid1/filemanager"
	"github.com/mehmetron/pid1/helpers"
	"github.com/mehmetron/pid1/packagemanager"
	"log"
	"net/http"
)

var server2 = "3000"

func SetupProject(w http.ResponseWriter, r *http.Request) {
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
	if command.Env == "js" || command.Env == "vitejs" {
		packagemanager.RunAdd("nodejs", command.Packages, true, "playground")
	} else if command.Env == "go" {
		packagemanager.RunAdd("go", command.Packages, true, "playground")
	} else if command.Env == "python" {
		packagemanager.RunAdd("python", command.Packages, true, "playground")
	}

	// Set user port based on .hackr config file
	server2 = command.Port

	err = helpers.WriteJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		fmt.Println("77 ", err)
	}
}
