package main

import (
	"C"
	"github.com/pkg/errors"
	"log"
	"os"
	"path"
)

var trashdir string
var logger = log.New(os.Stdout, "saferm: ", 0)

func main() {}

func init() {
	trashdir = TrashDir()

	// Check if the user's trash directory exists.
	// Make these checks early so less disruption is caused where possible.
	_, err := os.Stat(trashdir)
	if err != nil {
		directoryError(err)
	}
}

func directoryError(err error) {
	if os.IsNotExist(err) {
		errinfo := errors.Wrap(err, "User's trash directory does not exist")
		logger.Fatalf("%v", errinfo)
	}
	if os.IsPermission(err) {
		errinfo := errors.Wrap(err, "User's trash directory is not accessible")
		logger.Fatalf("%v", errinfo)
	}
	logger.Fatalf("%v", err)
}

func TrashDir() string {
	// Compute user's trash directory.
	// See https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html.
	datadir := os.Getenv("XDG_DATA_HOME")
	if datadir == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			logger.Fatalf("Couldn't find user's home directory")
		}
		datadir = path.Join(homedir, ".local/share")
	}
	trashdir := path.Join(datadir, "Trash/files")
	return trashdir
}

func SafeRemove(pathname string) C.int {
	newpath := path.Join(trashdir, path.Base(pathname))
	// TODO: Check for file collisions?
	err := os.Rename(pathname, newpath)
	if err != nil {
		directoryError(err)
	}
	return 0
}

//export remove
func remove(pathname *C.char) C.int {
	return SafeRemove(C.GoString(pathname))
}

//export unlinkat
func unlinkat(_ C.int, pathname *C.char, _ int) C.int {
	return SafeRemove(C.GoString(pathname))
}

//export rmdir
func rmdir(pathname *C.char) C.int {
	return SafeRemove(C.GoString(pathname))
}
