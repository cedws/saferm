package main

import (
	"C"
	"os"
	"path"
)

var trashdir string

func main() {}

// See https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html.
func init() {
	datahome := os.Getenv("XDG_DATA_HOME")
	if datahome == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		datahome = path.Join(homedir, ".local/share")
	}
	trashdir = path.Join(datahome, "Trash/files")
}

func saferm(pathname string) C.int {
	newpath := path.Join(trashdir, path.Base(pathname))
	os.Rename(pathname, newpath)
	return 0
}

//export remove
func remove(pathname *C.char) C.int {
	return saferm(C.GoString(pathname))
}

//export unlinkat
func unlinkat(dirfd C.int, pathname *C.char, flags int) C.int {
	return saferm(C.GoString(pathname))
}

//export rmdir
func rmdir(pathname *C.char) C.int {
	return saferm(C.GoString(pathname))
}
