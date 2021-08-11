package golang

import (
	"encoding/json"
	"fmt"
	"github.com/mehmetron/pid1/util"
)

type GoPackage struct {
	Path     string `json:"Path"`
	Version  string `json:"Version"`
	Indirect bool   `json:"Indirect"`
}

type AllPackages struct {
	Go      string      `json:"Go"`
	Require []GoPackage `json:"Require"`
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

// normalizePackageName implements NormalizePackageName for the Python
// backends.
func normalizePackageName(name util.PkgName) util.PkgName {
	nameStr := string(name)
	//nameStr = strings.ToLower(nameStr)
	//nameStr = strings.Replace(nameStr, "_", "-", -1)
	return util.PkgName(nameStr)
}

// GoBackend is a UPM backend for Node.js that uses NPM.
var GoBackend = util.LanguageBackend{
	Name:             "go-modules",
	Specfile:         "go.mod",
	Lockfile:         "go.sum",
	FilenamePatterns: []string{"*.go"},
	Quirks: util.QuirksAddRemoveAlsoLocks |
		util.QuirksAddRemoveAlsoInstalls,
	NormalizePackageName: normalizePackageName,
	GetPackageDir: func() string {
		// I have no clue
		return "i_dont_know"
		//return "node_modules"
	},
	Add: func(pkgs map[util.PkgName]util.PkgSpec, projectName string) {
		if !util.Exists("go.mod") {
			util.RunCmd([]string{"go", "mod", "init", projectName})
		}
		for name, spec := range pkgs {
			cmd := []string{"go", "get", "-v"}
			arg := string(name)
			if spec != "" {
				arg += "@" + string(spec)
			}
			cmd = append(cmd, arg)

			util.RunCmd(cmd)
		}
	},
	Remove: func(pkgs map[util.PkgName]bool) {
		//cmd := []string{"npm", "uninstall"}
		//for name, _ := range pkgs {
		//	cmd = append(cmd, string(name))
		//}
		//util.RunCmd(cmd)
		util.RunCmd([]string{"go", "mod", "tidy"})

	},
	Lock: func() {
		util.RunCmd([]string{"go", "mod", "init", "playground"})
	},
	Install: func() {
		util.RunCmd([]string{"go", "mod", "download"})
	},
	ListSpecfile: func() map[util.PkgName]util.PkgSpec {
		depList := util.GetCmdOutput([]string{"go", "mod", "edit", "-json"})
		fmt.Println("32 ", depList)
		fmt.Println("33 ", string(depList))

		dynamic := AllPackages{}
		json.Unmarshal(depList, &dynamic)

		pkgs := map[util.PkgName]util.PkgSpec{}
		for nameStr, data := range dynamic.Require {
			fmt.Println("83 ", nameStr)
			fmt.Println("84 ", data)
			pkgs[util.PkgName(data.Path)] = util.PkgSpec(trimLeftChar(data.Version))
		}
		fmt.Println("96 ", pkgs)
		return pkgs

	},
	ListLockfile: func() map[util.PkgName]util.PkgVersion {
		// Can't list deps in go.sum file

		depList := util.GetCmdOutput([]string{"go", "mod", "edit", "-json"})
		fmt.Println("32 ", depList)
		fmt.Println("33 ", string(depList))

		dynamic := AllPackages{}
		json.Unmarshal(depList, &dynamic)

		pkgs := map[util.PkgName]util.PkgVersion{}
		for nameStr, data := range dynamic.Require {
			fmt.Println("83 ", nameStr)
			fmt.Println("84 ", data)
			pkgs[util.PkgName(data.Path)] = util.PkgVersion(trimLeftChar(data.Version))
		}
		fmt.Println("96 ", pkgs)
		return pkgs

	},
}
