package filemanager

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/mehmetron/pid1/executor"
//	"io/ioutil"
//	"net/http"
//	"os"
//)
//
//type command executor.Command
//
//func (c *command) CreateFile() {
//	f, err := os.Create(c.Entrypoint)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = f.WriteString(c.Content)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//
//type File struct {
//	ID       int    `json:"id"`
//	Name     string `json:"name"`
//	Contents string `json:"contents"`
//	Folder   Folder `json:"folder"`
//	Edited   bool
//}
//
//type Folder struct {
//	ID   int    `json:"id"`
//	Name string `json:"name"`
//}
//
//type FileTree struct {
//	Main    File     `json:"Main"`
//	Files   []File   `json:"Files"`
//	Folders []Folder `json:"Folders"`
//}
//
//
//func Create(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer r.Body.Close()
//
//	fmt.Println(string(body))
//
//	var b FileTree
//	err = json.Unmarshal(body, &b)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	// fmt.Printf("what %s - %s", b.Name, b.Contents)
//
//	for _, folder := range b.Folders {
//		err := os.Mkdir(folder.Name, 0755)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	for _, file := range b.Files {
//		f, err := os.Create(file.Name)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		_, err = f.WriteString(file.Contents)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//
//	// f, err := os.Create(b.Name)
//	// check(err)
//	// _, err = f.WriteString(b.Contents)
//	// check(err)
//
//	// cmd := exec.Command("go", "run", "bob.go")
//	// err = cmd.Run()
//	// check(err)
//
//	// out, err := exec.Command("go", "run", b.Name).Output()
//	// if err != nil {
//	// 	fmt.Printf("%s", err)
//	// }
//	// output := string(out[:])
//	// fmt.Printf(output)
//
//	// err = os.Remove(b.Name)
//	// check(err)
//
//	fmt.Fprintf(w, "It worked")
//
//}
//
//func createDirs(m FileTree) {
//
//	for _, folder := range m.Folders {
//		os.Mkdir(folder.Name, 0755)
//	}
//
//	for _, file := range m.Files {
//
//		os.Chdir(file.Folder.Name)
//		f, err := os.Create(file.Name)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		_, err = f.WriteString(file.Contents)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//	}
//
//}
//
//
//
//func CreateFolder(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer r.Body.Close()
//
//	fmt.Println(string(body))
//
//	var b Folder
//	err = json.Unmarshal(body, &b)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	err = os.Mkdir(b.Name, 0755)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//}
//
//func CreateFiles(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer r.Body.Close()
//
//	fmt.Println(string(body))
//
//	var b File
//	err = json.Unmarshal(body, &b)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	f, err := os.Create(b.Name)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = f.WriteString(b.Contents)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//
//}
//
//
//// Update file/folder tree recursively
//func updateDirs(m FileTree) {
//	for _, file := range m.Files {
//		if file.Edited {
//			// update file content
//		}
//	}
//}
//
//
//func transplant() {
//	b := []byte(`{"Folders": [{"Id":593893508,"Name":"Ahmet", "IsDir": false, "Edited":false, "Contents":"123;\n234;\n345;"}, {"Id":8833984,"Name":"Pickle", "IsDir": true, "Edited":false, "Children":[{"Id":93957932,"Name":"James", "IsDir": true, "Edited":false, "Children":[{"Id":93957932,"Name":"Kerem", "IsDir": false, "Edited":true, "Contents":"something in here"}, {"Id":93957932,"Name":"Gopher", "IsDir": true, "Edited":false, "Children":[]}]}, {"Id":1133422455,"Name":"Bob", "IsDir": false, "Edited":false, "Contents":"while butt is big slap"}]}]}`)
//	var m FileTree
//	err := json.Unmarshal(b, &m)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//
//	parent := ""
//	for i := 0; i < len(m.Folders); i++ {
//		//create(m.Folders[i], parent)
//		fmt.Println(parent)
//	}
//}
//
