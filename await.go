// Await - Wait for a file to be created/deleted/modified while blocking
//
// Copyright 2019 TSH Labs
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [FILE]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Wait for a file to be created/deleted/modified while blocking.\n\n")
		flag.PrintDefaults()
	}

	interval := flag.Duration("interval", DEFAULT_INTERVAL, "How long to wait between checking that a file has been created/deleted/modified. T is a \"duration\" argument and so accepts values like '1s', '5m', '1h', etc. The default is one second ('1s').")
	notexists := flag.Bool("notexists", false, "Check that the provided FILE does not exist and exit as soon as it does not")
	exists := flag.Bool("exists", true, "Check that the provided FILE exists and exit as soon as it does")
	fresh := flag.Duration("fresh", DEFAULT_DURATION, "Check if that the provided FILE has been modified in the last T duration and exit as soon as it has not. T is a \"duration\" argument and so accepts values like '1s', '5m', '1h', etc.")

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
