package filemanager

import (
	"fmt"
	"log"
	"os"
)

func CreateFile(fileName, code string) {
	fmt.Println("Creating file")

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	_, err = f.WriteString(code)

	if err != nil {
		log.Fatal(err)
	}
}

func DeleteFile() {
	fmt.Println("Deleting file")
}

func ReadFile() {
	fmt.Println("Reading file")
}

func UpdateFile() {
	fmt.Println("Updating file")
}
