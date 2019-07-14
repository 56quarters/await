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

const DEFAULT_INTERVAL = 2

type filetest interface {
	satisfied() (bool, error)
}

type filefreshness struct {
	path string
	age  time.Duration
}

func newFileFreshness(path string, age time.Duration) filefreshness {
	return filefreshness { path: path, age: age }
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
	return fileexists { path: path }
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
	return filenotexists { path: path }
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
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [PID]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Repeatedly try to stop a process with SIGTERM and eventually SIGKILL\n\n")
		flag.PrintDefaults()
	}

	interval := flag.Duration("interval", DEFAULT_INTERVAL, "How long between checks, in seconds")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("File path is required")
	}

	path := flag.Arg(0)
	check := newFileFreshness(path, 120 * time.Second)

	for {
		if done, err := check.satisfied(); done {
			break
		} else if err != nil {
			log.Fatalf("Problem checking file: %s", err)
		}

		time.Sleep(*interval * time.Second)
	}
}
