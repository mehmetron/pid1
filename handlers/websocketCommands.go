package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mehmetron/pid1/filemanager"
	"io"
	"log"
	"net/http"
)

var upgrader2 = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func messageClient(client *websocket.Conn, msg filemanager.WsCommand) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
	}
}

func WebsocketCommands(w http.ResponseWriter, r *http.Request) {
	fmt.Println("9 WebsocketCommands")

	ws, err := upgrader2.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()

	for {
		var msg filemanager.WsCommand
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("36 ", err)
			break
		}
		// send new message to the channel
		fmt.Println("40  ", msg)

		switch msg.Command {
		case "execute":
			fmt.Println("first command")
		case "update":
			fmt.Println("second command")
		case "kill":
			fmt.Println("third command")
		case "port":
			fmt.Println("fourth command")
		default:
			fmt.Println("default command")
		}
	}

}
