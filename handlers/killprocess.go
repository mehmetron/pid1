package handlers

import (
	"fmt"
	"github.com/mehmetron/pid1/executor"
	"github.com/mehmetron/pid1/helpers"
	"net/http"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Killing running process")

	s := executor.KillProcess()

	err := helpers.WriteJSON(w, http.StatusOK, s, nil)
	if err != nil {
		fmt.Println("77 ", err)
	}
}
