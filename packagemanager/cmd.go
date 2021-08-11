package packagemanager

import (
	"fmt"
	"github.com/mehmetron/pid1/util"
	"os"
	"strings"
)

// Quiet is true if --quiet was passed on the command line.
var Quiet bool

// subroutineSilencer is used to easily enable and restore
// Quiet for part of a function.
type subroutineSilencer struct {
	origQuiet bool
}

// silenceSubroutines turns on Quiet and returns a struct that
// can be used to restore its value. This only happens if
// UPM_SILENCE_SUBROUTINES is non-empty.
func silenceSubroutines() subroutineSilencer {
	s := subroutineSilencer{origQuiet: Quiet}
	if os.Getenv("UPM_SILENCE_SUBROUTINES") != "" {
		Quiet = true
	}
	return s
}

// restore restores the previous value of Quiet.
func (s *subroutineSilencer) restore() {
	Quiet = s.origQuiet
}

// pkgNameAndSpec is a tuple of a PkgName and a PkgSpec. It's used to
// put both of them as a value in the same map entry.
type pkgNameAndSpec struct {
	name util.PkgName
	spec util.PkgSpec
}

// maybeLock either runs lock or not, depending on the backend, store,
// and command-line options. It returns true if it actually ran lock.
func maybeLock(b util.LanguageBackend) bool {
	if b.QuirksIsNotReproducible() {
		return false
	}

	if !util.Exists(b.Specfile) {
		return false
	}

	if !util.Exists(b.Lockfile) {
		b.Lock()
		return true
	}

	return false
}

// maybeInstall either runs install or not, depending on the backend,
// store, and command-line options.
func maybeInstall(b util.LanguageBackend) {
	if b.QuirksIsReproducible() {
		if !util.Exists(b.Lockfile) {
			return
		}
		b.Install()
	} else {
		if !util.Exists(b.Specfile) {
			return
		}
		b.Install()
	}
}

// deleteLockfile deletes the project's lockfile, if one exists.
func deleteLockfile(b util.LanguageBackend) {
	if util.Exists(b.Lockfile) {
		fmt.Println("-->delete " + b.Lockfile)
		os.Remove(b.Lockfile)
	}
}

// runInstall implements 'upm install'.
func runInstall(language string) {
	b := GetBackend(language)

	maybeInstall(b)
}

// runLock implements 'upm lock'.
func runLock(language string, upgrade bool) {
	b := GetBackend(language)

	if upgrade {
		deleteLockfile(b)
	}

	didLock := maybeLock(b)

	if !(didLock && b.QuirksDoesLockAlsoInstall()) {
		maybeInstall(b)
	}

}

// RunAdd implements 'upm add'.
func RunAdd(language string, args []string, upgrade bool, name string) {

	b := GetBackend(language)

	fmt.Println("119 ", b)

	// Map from normalized package names to the corresponding
	// original package names and specs.
	normPkgs := map[util.PkgName]pkgNameAndSpec{}
	for _, arg := range args {
		fmt.Println("125 ", arg)
		fields := strings.SplitN(arg, " ", 2)
		fmt.Println("127 ", fields)
		name := util.PkgName(fields[0])
		var spec util.PkgSpec
		if len(fields) >= 2 {
			spec = util.PkgSpec(fields[1])
		}
		fmt.Println("133 ", spec)

		fmt.Println("135 ", b.NormalizePackageName(name))
		fmt.Println("136 ", normPkgs)
		normPkgs[name] = pkgNameAndSpec{
			name: name,
			spec: spec,
		}

		fmt.Println("140 ")
		fmt.Println("141 ", normPkgs)
	}

	fmt.Println("141 ")
	fmt.Println("142", normPkgs)

	fmt.Println("149 ", util.Exists(b.Specfile))

	// Delete from normPkgs map  packages that are already installed
	if util.Exists(b.Specfile) {
		s := silenceSubroutines()
		for name, _ := range b.ListSpecfile() {
			fmt.Println("153  ", name)
			delete(normPkgs, b.NormalizePackageName(name))
			//delete(normPkgs, name)
		}
		s.restore()
	}

	// Install everything from scratch to install latest dependency versions
	if upgrade {
		deleteLockfile(b)
	}

	// Install packages in normPkgs
	if len(normPkgs) >= 1 {
		pkgs := map[util.PkgName]util.PkgSpec{}
		for _, nameAndSpec := range normPkgs {
			pkgs[nameAndSpec.name] = nameAndSpec.spec
		}
		b.Add(pkgs, name)
	}

	// Generate lock file and/or install
	if len(normPkgs) == 0 || b.QuirksDoesAddRemoveNotAlsoLock() {
		didLock := maybeLock(b)

		if !(didLock && b.QuirksDoesLockAlsoInstall()) {
			maybeInstall(b)
		}
	} else if len(normPkgs) == 0 || b.QuirksDoesAddRemoveNotAlsoInstall() {
		maybeInstall(b)
	}

}

// runRemove implements 'upm remove'.
func runRemove(language string, args []string, upgrade bool,
	forceLock bool, forceInstall bool) {

	b := GetBackend(language)

	if !util.Exists(b.Specfile) {
		return
	}

	s := silenceSubroutines()
	specfilePkgs := b.ListSpecfile()
	s.restore()

	// Map whose keys are normalized package names.
	normSpecfilePkgs := map[util.PkgName]bool{}
	for name := range specfilePkgs {
		normSpecfilePkgs[b.NormalizePackageName(name)] = true
		//normSpecfilePkgs[name] = true
	}

	// Map from normalized package names to original package
	// names.
	normPkgs := map[util.PkgName]util.PkgName{}
	for _, arg := range args {
		name := util.PkgName(arg)
		norm := b.NormalizePackageName(name)
		//norm := name
		if _, ok := normSpecfilePkgs[norm]; ok {
			normPkgs[norm] = name
		}
	}

	if upgrade {
		deleteLockfile(b)
	}

	if len(normPkgs) >= 1 {
		pkgs := map[util.PkgName]bool{}
		for _, name := range normPkgs {
			pkgs[name] = true
		}
		b.Remove(pkgs)
	}

	if len(normPkgs) == 0 || b.QuirksDoesAddRemoveNotAlsoLock() {
		didLock := maybeLock(b)

		if !(didLock && b.QuirksDoesLockAlsoInstall()) {
			maybeInstall(b)
		}
	} else if len(normPkgs) == 0 || b.QuirksDoesAddRemoveNotAlsoInstall() {
		maybeInstall(b)
	}

}
