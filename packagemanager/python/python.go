// Package python provides backends for Python 2 and 3 using Poetry.
package python

import (
	"encoding/json"
	"fmt"
	"github.com/mehmetron/pid1/util"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// this generates a mapping of pypi packages <-> modules
// moduleToPypiPackage pypiPackageToModules are provided
//go:generate go run ./gen_pypi_map -from pypi_packages.json -pkg python -out pypi_map.gen.go

// pypiEntry represents one element of the response we get from
// the PyPI API search results.
type pypiEntry struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Version string `json:"version"`
}

// pypiEntryInfoResponse is a wrapper around pypiEntryInfo
// that matches the format of the REST API
type pypiEntryInfoResponse struct {
	Info pypiEntryInfo `json:"info"`
}

// pypiEntryInfo represents the response we get from the
// PyPI API on doing a single-package lookup.
type pypiEntryInfo struct {
	Author        string   `json:"author"`
	AuthorEmail   string   `json:"author_email"`
	HomePage      string   `json:"home_page"`
	License       string   `json:"license"`
	Name          string   `json:"name"`
	ProjectURL    string   `json:"project_url"`
	PackageURL    string   `json:"package_url"`
	BugTrackerURL string   `json:"bugtrack_url"`
	DocsURL       string   `json:"docs_url"`
	RequiresDist  []string `json:"requires_dist"`
	Summary       string   `json:"summary"`
	Version       string   `json:"version"`
}

// pyprojectTOML represents the relevant parts of a pyproject.toml
// file.
type pyprojectTOML struct {
	Tool struct {
		Poetry struct {
			Name string `json:"name"`
			// interface{} because they can be either
			// strings or maps (why?? good lord).
			Dependencies    map[string]interface{} `json:"dependencies"`
			DevDependencies map[string]interface{} `json:"dev-dependencies"`
		} `json:"poetry"`
	} `json:"tool"`
}

// poetryLock represents the relevant parts of a poetry.lock file, in
// TOML format.
type poetryLock struct {
	Package []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"package"`
}

// moduleMetadata represents the information that could be associated with
// a module using a #upm pragma
type modulePragmas struct {
	Package string `json:"package"`
}

// normalizeSpec returns the version string from a Poetry spec, or the
// empty string. The Poetry spec may be either a string or a
// map[string]interface{} with a "version" key that is a string. If
// neither, then the empty string is returned.
func normalizeSpec(spec interface{}) string {
	switch spec := spec.(type) {
	case string:
		return spec
	case map[string]interface{}:
		switch spec := spec["version"].(type) {
		case string:
			return spec
		}
	}
	return ""
}

// normalizePackageName implements NormalizePackageName for the Python
// backends.
func normalizePackageName(name util.PkgName) util.PkgName {
	nameStr := string(name)
	nameStr = strings.ToLower(nameStr)
	nameStr = strings.Replace(nameStr, "_", "-", -1)
	return util.PkgName(nameStr)
}

// pythonMakeBackend returns a language backend for a given version of
// Python. name is either "python2" or "python3", and python is the
// name of an executable (either a full path or just a name like
// "python3") to use when invoking Python. (This is used to implement
// UPM_PYTHON2 and UPM_PYTHON3.)
func pythonMakeBackend(name string, python string) util.LanguageBackend {

	return util.LanguageBackend{
		Name:             "python-" + name + "-poetry",
		Specfile:         "pyproject.toml",
		Lockfile:         "poetry.lock",
		FilenamePatterns: []string{"*.py"},
		Quirks: util.QuirksAddRemoveAlsoLocks |
			util.QuirksAddRemoveAlsoInstalls,
		NormalizePackageName: normalizePackageName,
		GetPackageDir: func() string {
			// Check if we're already inside an activated
			// virtualenv. If so, just use it.
			if venv := os.Getenv("VIRTUAL_ENV"); venv != "" {
				return venv
			}

			// Ideally Poetry would provide some way of
			// actually checking where the virtualenv will
			// go. But it doesn't. So we have to
			// reimplement the logic ourselves, which is
			// totally fragile and disgusting. (No, we
			// can't use 'poetry run which python' because
			// that will *create* a virtualenv if one
			// doesn't exist, and there's no workaround
			// for that without mutating the global config
			// file.)
			//
			// Note, we don't yet support Poetry's
			// settings.virtualenvs.in-project. That would
			// be a pretty easy fix, though. (Why is this
			// so complicated??)

			outputB := util.GetCmdOutput([]string{
				python, "-m", "poetry",
				"config", "settings.virtualenvs.path",
			})
			var path string
			if err := json.Unmarshal(outputB, &path); err != nil {

				fmt.Printf("parsing output from Poetry: %s\n", err)
			}

			base := ""
			if util.Exists("pyproject.toml") {
				var cfg pyprojectTOML
				if _, err := toml.DecodeFile("pyproject.toml", &cfg); err != nil {
					fmt.Printf("%s\n", err.Error())
				}
				base = cfg.Tool.Poetry.Name
			}

			if base == "" {
				cwd, err := os.Getwd()
				if err != nil {
					fmt.Printf("%s\n", err)
				}
				base = strings.ToLower(filepath.Base(cwd))
			}

			version := strings.TrimSpace(string(util.GetCmdOutput([]string{
				python, "-c",
				`import sys; print(".".join(map(str, sys.version_info[:2])))`,
			})))

			return filepath.Join(path, base+"-py"+version)
		},
		Add: func(pkgs map[util.PkgName]util.PkgSpec, projectName string) {
			// Initalize the specfile if it doesnt exist
			if !util.Exists("pyproject.toml") {
				cmd := []string{python, "-m", "poetry", "init", "--no-interaction"}

				if projectName != "" {
					cmd = append(cmd, "--name", projectName)
				}

				util.RunCmd(cmd)
			}

			cmd := []string{python, "-m", "poetry", "add"}
			for name, spec := range pkgs {
				name := string(name)
				spec := string(spec)

				// NB: this doesn't work if spec has
				// spaces in it, because of a bug in
				// Poetry that can't be worked around.
				// It looks like that bug might be
				// fixed in the 1.0 release though :/
				if spec != "" {
					cmd = append(cmd, name+" "+spec)
				} else {
					cmd = append(cmd, name)
				}
			}
			util.RunCmd(cmd)
		},
		Remove: func(pkgs map[util.PkgName]bool) {
			cmd := []string{python, "-m", "poetry", "remove"}
			for name, _ := range pkgs {
				cmd = append(cmd, string(name))
			}
			util.RunCmd(cmd)
		},
		Lock: func() {
			util.RunCmd([]string{python, "-m", "poetry", "lock"})
		},
		Install: func() {
			// Unfortunately, this doesn't necessarily uninstall
			// packages that have been removed from the lockfile,
			// which happens for example if 'poetry remove' is
			// interrupted. See
			// <https://github.com/sdispater/poetry/issues/648>.
			util.RunCmd([]string{python, "-m", "poetry", "install"})
		},
		ListSpecfile: func() map[util.PkgName]util.PkgSpec {
			pkgs, err := listSpecfile()
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}

			return pkgs
		},
		ListLockfile: func() map[util.PkgName]util.PkgVersion {
			var cfg poetryLock
			if _, err := toml.DecodeFile("poetry.lock", &cfg); err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			pkgs := map[util.PkgName]util.PkgVersion{}
			for _, pkgObj := range cfg.Package {
				name := util.PkgName(pkgObj.Name)
				version := util.PkgVersion(pkgObj.Version)
				pkgs[name] = version
			}
			return pkgs
		},
	}
}

func listSpecfile() (map[util.PkgName]util.PkgSpec, error) {
	var cfg pyprojectTOML
	if _, err := toml.DecodeFile("pyproject.toml", &cfg); err != nil {
		return nil, err
	}
	pkgs := map[util.PkgName]util.PkgSpec{}
	for nameStr, spec := range cfg.Tool.Poetry.Dependencies {
		if nameStr == "python" {
			continue
		}

		specStr := normalizeSpec(spec)
		if specStr == "" {
			continue
		}
		pkgs[util.PkgName(nameStr)] = util.PkgSpec(specStr)
	}
	for nameStr, spec := range cfg.Tool.Poetry.DevDependencies {
		if nameStr == "python" {
			continue
		}

		specStr := normalizeSpec(spec)
		if specStr == "" {
			continue
		}
		pkgs[util.PkgName(nameStr)] = util.PkgSpec(specStr)
	}

	return pkgs, nil
}

// getPython2 returns either "python2" or the value of the UPM_PYTHON2
// environment variable.
func getPython2() string {
	python2 := os.Getenv("UPM_PYTHON2")
	if python2 != "" {
		return python2
	} else {
		return "python2"
	}
}

// getPython3 returns either "python3" or the value of the UPM_PYTHON3
// environment variable.
func getPython3() string {
	python3 := os.Getenv("UPM_PYTHON3")
	if python3 != "" {
		return python3
	} else {
		return "python3"
	}
}

// Python2Backend is a UPM backend for Python 2 that uses Poetry.
var Python2Backend = pythonMakeBackend("python2", getPython2())

// Python3Backend is a UPM backend for Python 3 that uses Poetry.
var Python3Backend = pythonMakeBackend("python3", getPython3())
