package main

import (
	"C"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"os"
	"path"
	"syscall"
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
		logger.Fatal(errinfo)
	}
	if os.IsPermission(err) {
		errinfo := errors.Wrap(err, "User's trash directory is not accessible")
		logger.Fatal(errinfo)
	}
	logger.Fatal(err)
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

func SameDevice(a os.FileInfo, b os.FileInfo) bool {
	a_dev := a.Sys().(*syscall.Stat_t).Dev
	b_dev := b.Sys().(*syscall.Stat_t).Dev

	return a_dev == b_dev
}

func SafeRemove(dirfd int, pathname string, flags int) error {
	o_flags := unix.O_RDONLY | unix.O_NONBLOCK | unix.O_CLOEXEC
	pathfd, err := unix.Openat(dirfd, pathname, o_flags, 0)
	if err != nil {
		return err
	}
	// Create a Go File from the file descriptor obtained earlier.
	srcfile := os.NewFile(uintptr(pathfd), pathname)

	srcstat, err := srcfile.Stat()
	if err != nil {
		return err
	}

	deststat, err := os.Stat(trashdir)
	if err != nil {
		return err
	}

	dest := path.Join(trashdir, path.Base(pathname))

	// Compare source and destination devices to see if they are the same.
	// If they are, it's a cheap and simple operation. If not, we have to copy everything.
	if SameDevice(srcstat, deststat) {
		err := unix.Renameat(dirfd, pathname, unix.AT_FDCWD, dest)
		if err != nil {
			return err
		}
	} else {
		var err error

		destfile, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer destfile.Close()

		_, err = io.Copy(destfile, srcfile)
		if err != nil {
			return err
		}

		// The SameDevice == true branch 'removes' the path
		// as a side effect. We should do the same in this branch.
		err = unix.Unlinkat(dirfd, pathname, flags)
		if err != nil {
			return err
		}
	}

	return nil
}

//export unlinkat
func unlinkat(dirfd C.int, pathname *C.char, flags int) C.int {
	if flags&unix.AT_REMOVEDIR != unix.AT_REMOVEDIR {
		err := SafeRemove(int(dirfd), C.GoString(pathname), flags)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		err := unix.Unlinkat(int(dirfd), C.GoString(pathname), flags)
		if err != nil {
			logger.Fatal(err)
		}
	}

	return 0
}
