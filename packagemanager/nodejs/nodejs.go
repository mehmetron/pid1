package nodejs

import (
	"encoding/json"
	"fmt"
	"github.com/mehmetron/pid1/util"
	"io/ioutil"
)

// packageLockJSON represents the relevant data in a package-lock.json
// file.
type packageLockJSON struct {
	Dependencies map[string]struct {
		Version string `json:"version"`
	} `json:"dependencies"`
}

// packageJSON represents the relevant data in a package.json file.
type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// normalizePackageName implements NormalizePackageName for the Python
// backends.
func normalizePackageName(name util.PkgName) util.PkgName {
	nameStr := string(name)
	//nameStr = strings.ToLower(nameStr)
	//nameStr = strings.Replace(nameStr, "_", "-", -1)
	return util.PkgName(nameStr)
}

// nodejsPatterns is the FilenamePatterns value for NodejsBackend.
var nodejsPatterns = []string{"*.js", "*.ts", "*.jsx", "*.tsx"}

// nodejsListSpecfile implements ListSpecfile for nodejs-yarn and
// nodejs-npm.
func nodejsListSpecfile() map[util.PkgName]util.PkgSpec {
	contentsB, err := ioutil.ReadFile("package.json")
	if err != nil {
		fmt.Printf("package.json: %s\n", err)
	}
	var cfg packageJSON
	if err := json.Unmarshal(contentsB, &cfg); err != nil {
		fmt.Printf("package.json: %s\n", err)
	}
	pkgs := map[util.PkgName]util.PkgSpec{}
	for nameStr, specStr := range cfg.Dependencies {
		pkgs[util.PkgName(nameStr)] = util.PkgSpec(specStr)
	}
	for nameStr, specStr := range cfg.DevDependencies {
		pkgs[util.PkgName(nameStr)] = util.PkgSpec(specStr)
	}
	return pkgs
}

// NodejsNPMBackend is a UPM backend for Node.js that uses NPM.
var NodejsNPMBackend = util.LanguageBackend{
	Name:             "nodejs-npm",
	Specfile:         "package.json",
	Lockfile:         "package-lock.json",
	FilenamePatterns: nodejsPatterns,
	Quirks: util.QuirksAddRemoveAlsoLocks |
		util.QuirksAddRemoveAlsoInstalls |
		util.QuirksLockAlsoInstalls,
	NormalizePackageName: normalizePackageName,
	GetPackageDir: func() string {
		return "node_modules"
	},
	Add: func(pkgs map[util.PkgName]util.PkgSpec, projectName string) {
		if !util.Exists("package.json") {
			util.RunCmd([]string{"npm", "init", "-y"})

			packageFile := `
{
  "name": "playground",
  "version": "1.0.0",
  "description": "",
  "main": "vite.config.js",
  "scripts": {
    "dev":"vite"
  },
  "keywords": [],
  "author": "",
  "license": "ISC"
}
`
			//d1 := []byte("{\n  \"name\": \"playground\",\n  \"version\": \"1.0.0\",\n  \"description\": \"\",\n  \"main\": \"vite.config.js\",\n  \"scripts\": {\n    \"dev\":\"vite\"\n  },\n  \"keywords\": [],\n  \"author\": \"\",\n  \"license\": \"ISC\"\n}\n")
			err := ioutil.WriteFile("/app/playground/package.json", []byte(packageFile), 0644)
			if err != nil {
				fmt.Println("70 new package.json failed ", err)
			}
		}
		cmd := []string{"npm", "install"}
		for name, spec := range pkgs {
			arg := string(name)
			if spec != "" {
				arg += "@" + string(spec)
			}
			cmd = append(cmd, arg)
		}
		util.RunCmd(cmd)
	},
	Remove: func(pkgs map[util.PkgName]bool) {
		cmd := []string{"npm", "uninstall"}
		for name, _ := range pkgs {
			cmd = append(cmd, string(name))
		}
		util.RunCmd(cmd)
	},
	Lock: func() {
		util.RunCmd([]string{"npm", "install"})
	},
	Install: func() {
		util.RunCmd([]string{"npm", "ci"})
	},
	ListSpecfile: nodejsListSpecfile,
	ListLockfile: func() map[util.PkgName]util.PkgVersion {
		contentsB, err := ioutil.ReadFile("package-lock.json")
		if err != nil {
			fmt.Printf("package-lock.json: %s", err)
		}
		var cfg packageLockJSON
		if err := json.Unmarshal(contentsB, &cfg); err != nil {
			fmt.Printf("package-lock.json: %s\n", err)
		}
		pkgs := map[util.PkgName]util.PkgVersion{}
		for nameStr, data := range cfg.Dependencies {
			pkgs[util.PkgName(nameStr)] = util.PkgVersion(data.Version)
		}
		return pkgs
	},
}
