package handlers

import (
	"errors"
	"fmt"
	"github.com/mehmetron/pid1/filemanager"
	"github.com/mehmetron/pid1/helpers"
	"log"
	"net/http"
	"os"
)

func UpdateFilesVite(w http.ResponseWriter, r *http.Request) {
	fmt.Println("13 execute handler")

	var command filemanager.UpdateFiles
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

	if len(command.Deleted) > 0 {
		for i, s := range command.Deleted {
			fmt.Println("204 Deleting file ", i, s)
			var err = os.Remove(s)
			if err != nil {
				fmt.Println("215 ", err)
			}
		}
	}

	if len(command.Updated) > 0 {
		for i, s := range command.Updated {
			fmt.Println("213 Creating File", i, s)
			// check if file exists
			_, err := os.Stat(s.Name)

			// create file if not exists
			if os.IsNotExist(err) {
				var file, err = os.Create(s.Name)
				if err != nil {
					fmt.Println("223 ", err)
				}
				defer file.Close()
			}

			fmt.Println("214 Updating file", i, s)

			file, err := os.OpenFile(s.Name, os.O_RDWR, 0644)
			if err != nil {
				fmt.Println("218 ", err)
			}
			defer file.Close()

			// Write some text line-by-line to file.
			_, err = file.WriteString(s.Content)
			if err != nil {
				fmt.Println("225 ", err)
			}

			// Save file changes.
			err = file.Sync()
			if err != nil {
				fmt.Println("235 ", err)
			}

		}
	}

	err = helpers.WriteJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		fmt.Println("77 ", err)
	}
}
