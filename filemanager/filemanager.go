package filemanager

import (
	"fmt"
	"io/ioutil"
	"os"
)

type WsCommand struct {
	Command  string   `json:"command"`
	Env      string   `json:"env"`
	Files    Files    `json:"files"`
	Packages []string `json:"packages"`
	Port     string   `json:"port"`
}

type Command struct {
	Env      string   `json:"env"`
	Files    Files    `json:"files"`
	Packages []string `json:"packages"`
	Port     string   `json:"port"`
}

type UpdateFiles struct {
	Deleted []string `json:"deleted"`
	Updated Files    `json:"updated"`
}

type Files []struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Take files - Updates playground to match object of files
// get all files in playground and put in map name>>contents
// for each file in object:
//  if exists in map, check if contents equal, if not, update contents
//  if not exists in map, create file
//  Delete from map
// for leftovers in map:
//  delete those files

func Create(files Files) error {

	existingFiles := make(map[string]string)
	filesToIgnore := map[string]bool{
		// Javascript files to ignore
		"package-lock.json": true,
		"package.json":      true,
		// Python files to ignore
		"pyproject.toml": true,
		"poetry.lock":    true,
	}

	// Get all files in playground folder
	files2, err := ioutil.ReadDir("/app/playground")
	if err != nil {
		return fmt.Errorf("33 problem reading dir: %v", err)
	}

	// Go over all files in playground folder and store in map
	for _, file := range files2 {
		fmt.Println("36 ", file.Name())

		_, found := filesToIgnore[file.Name()]
		if file.IsDir() || found {
			fmt.Println("48 Ignoring file: ", file.Name())
		} else {
			content, err := ReadFile(file.Name())
			if err != nil {
				return fmt.Errorf("41 problem reading file: %v", err)
			}
			existingFiles[file.Name()] = content
		}
	}

	// Create all files in passed in struct
	for _, file := range files {

		// If file in map, then delete it from map
		if x, found := existingFiles[file.Name]; found {
			fmt.Println("50 ", x)
			delete(existingFiles, file.Name)
		}

		// Create file
		err := CreateFile(file.Name, file.Content)
		if err != nil {
			return fmt.Errorf("43 problem creating file: %v", err)
		}

	}

	// Delete any files still left in map
	for key, value := range existingFiles {
		fmt.Println(key, value)
		err := DeleteFile(key)
		if err != nil {
			return fmt.Errorf("70 problem after calling Deletefile function: %v", err)
		}
		delete(existingFiles, key)
	}

	return nil
}

func CreateFile(fileName, code string) error {
	fmt.Println("Creating file: ", fileName)

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(code)
	if err != nil {
		return fmt.Errorf("86 problem writing to file: %v", err)
	}

	return nil
}

func DeleteFile(fileName string) error {
	fmt.Println("Deleting file: ", fileName)

	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("97 problem deleting file: %v", err)
	}

	return nil
}

func ReadFile(fileName string) (string, error) {
	fmt.Println("Reading file: ", fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("106 problem reading file: %v", err)
	}
	text := string(content)

	return text, nil
}

func UpdateFile(fileName, content string) error {
	fmt.Println("Updating file: ", fileName)

	err := ioutil.WriteFile(fileName, []byte(content), 0)
	if err != nil {
		return fmt.Errorf("118 problem writing to file: %v", err)
	}

	return nil
}
