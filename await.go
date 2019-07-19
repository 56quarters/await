//
//
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const DEFAULT_INTERVAL = 1 * time.Second
const DEFAULT_DURATION = 0

type filetest interface {
	satisfied() (bool, error)
}

type filefreshness struct {
	path string
	age  time.Duration
}

func newFileFreshness(path string, age time.Duration) filefreshness {
	return filefreshness{path: path, age: age}
}

func (f filefreshness) satisfied() (bool, error) {
	res, err := os.Stat(f.path)
	if err != nil {
		return false, err
	}

	mod := res.ModTime()
	return time.Since(mod) > f.age, nil
}

type fileexists struct {
	path string
}

func newFileExists(path string) fileexists {
	return fileexists{path: path}
}

func (f fileexists) satisfied() (bool, error) {
	if _, err := os.Stat(f.path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

type filenotexists struct {
	path string
}

func newFileNotExists(path string) filenotexists {
	return filenotexists{path: path}
}

func (f filenotexists) satisfied() (bool, error) {
	if _, err := os.Stat(f.path); err == nil {
		return false, nil
	} else if os.IsNotExist(err) {
		return true, nil
	} else {
		return false, err
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [PATH]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Check if a file exists (or meets other criteria) in a "+
			"loop, blocking until it does\n\n")
		flag.PrintDefaults()
	}

	interval := flag.Duration("interval", DEFAULT_INTERVAL,
		"How long between file checks, as a duration (example '1s', '30s', '1m', etc.)")
	notexists := flag.Bool("notexists", false, "Check if the file does *not* exist")
	exists := flag.Bool("exists", true, "Check if the file exists. This is the default behavior")
	fresh := flag.Duration("fresh", DEFAULT_DURATION, "Check if the file has *not* been modified "+
		"within the specified amount of time, as a duration (example '1s', '30s', '1m', etc.)")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("File path is required")
	}

	path := flag.Arg(0)

	var check filetest = nil
	if *fresh != DEFAULT_DURATION {
		check = newFileFreshness(path, *fresh)
	} else if *notexists {
		check = newFileNotExists(path)
	} else if *exists {
		check = newFileExists(path)
	} else {
		log.Fatal("No file check specified")
	}

	for {
		if done, err := check.satisfied(); done {
			break
		} else if err != nil {
			log.Fatalf("Problem checking file: %s", err)
		}

		time.Sleep(*interval)
	}
}
