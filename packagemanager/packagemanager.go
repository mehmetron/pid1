// Package packagemanager contains the language-specific UPM backends, and
// logic for selecting amongst them.
package packagemanager

import (
	"fmt"
	"github.com/mehmetron/pid1/packagemanager/golang"
	"github.com/mehmetron/pid1/packagemanager/nodejs"
	"github.com/mehmetron/pid1/packagemanager/python"
	"github.com/mehmetron/pid1/util"
	"strings"
)

// languageBackends is a slice of language backends which may be used
// from the command line.
//
// If more than one backend might match the same project, then one
// that comes first in this list will be used.
var languageBackends = []util.LanguageBackend{
	nodejs.NodejsNPMBackend,
	golang.GoBackend,
	python.Python3Backend,
	python.Python2Backend,
}

// matchesLanguage checks if a language backend matches a value for
// the --lang argument. For example, the python-python3-poetry backend
// will match --lang=python-poetry and --lang=python3 but not
// --lang=python2. This is used to filter the available language
// backends before trying to guess which one should be used.
func matchesLanguage(b util.LanguageBackend, language string) bool {
	bParts := map[string]bool{}
	for _, bPart := range strings.Split(b.Name, "-") {
		bParts[bPart] = true
	}
	for _, lPart := range strings.Split(language, "-") {
		if !bParts[lPart] {
			return false
		}
	}
	return true
}

// GetBackend returns the language backend for a given --lang argument
// value. If none is applicable, it exits the process.
func GetBackend(language string) util.LanguageBackend {
	backends := languageBackends
	if language != "" {
		filteredBackends := []util.LanguageBackend{}
		for _, b := range backends {
			if matchesLanguage(b, language) {
				filteredBackends = append(filteredBackends, b)
			}
		}
		switch len(filteredBackends) {
		case 0:
			fmt.Printf("no such language: %s\n", language)
		case 1:
			return filteredBackends[0]
		default:
			backends = filteredBackends
		}

	}
	for _, b := range backends {
		if util.Exists(b.Specfile) &&
			util.Exists(b.Lockfile) {
			return b
		}
	}
	for _, b := range backends {
		if util.Exists(b.Specfile) ||
			util.Exists(b.Lockfile) {
			return b
		}
	}
	for _, b := range backends {
		for _, p := range b.FilenamePatterns {
			if util.PatternExists(p) {
				return b
			}
		}
	}
	if language == "" {
		fmt.Println("could not autodetect a language for your project")
	}
	return backends[0]
}

// GetBackendNames returns a slice of the canonical names (e.g.
// python-python3-poetry, not just python3) for all the backends
// listed in languageBackends.
func GetBackendNames() []string {
	backendNames := []string{}
	for _, b := range languageBackends {
		backendNames = append(backendNames, b.Name)
	}
	return backendNames
}
