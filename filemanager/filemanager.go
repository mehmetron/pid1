package filemanager

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Files []struct {
	Name    string
	Content string
}

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

func DeleteFile(fileName string) {
	fmt.Println("Deleting file")

	err := os.Remove(fileName)
	if err != nil {
		fmt.Println(err)
	}
}

func ReadFile(fileName string) (string, error) {
	fmt.Println("Reading file")

	content, err := ioutil.ReadFile(fileName)
	text := string(content)

	return text, err

}

func UpdateFile(fileName, content string) {
	fmt.Println("Updating file")

	err := ioutil.WriteFile(fileName, []byte(content), 0)
	if err != nil {
		fmt.Println(err)
	}
}
