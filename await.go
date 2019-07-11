//
//
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	//"syscall"
	//"time"
)

const DEFAULT_INTERVAL = 2

type filetest interface {
	satisfied() (bool, error)
}

type filefreshness struct {
	path string
	age int64
}

func (f filefreshness) statisfiled() (bool, error) {
	return true, nil
}

type fileexists struct {
	path string
}

func (f fileexists) statisfied() (bool, error) {
	return true, nil
}

type filenotexists struct {
	path string
}

func (f filenotexists) satisfied() (bool, error) {
	return true, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [PID]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Repeatedly try to stop a process with SIGTERM and eventually SIGKILL\n\n")
		flag.PrintDefaults()
	}

	//_ interval := flag.Int("interval", DEFAULT_INTERVAL, "How long between checks, in seconds")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("PID is required")
	}
}
