package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var debug = false

func main() {
	path := flag.String("p", "", "file path")
	line := flag.Int("l", -1, "line number")

	flag.Parse()
	if *path == "" || *line == -1 {
		flag.Usage()
		os.Exit(1)
	}

	// generate absolute file path if needed
	if !filepath.IsAbs(*path) {
		d, e := os.Getwd()
		if e != nil {
			log.Fatal(e)
		}
		*path = filepath.Join(d, *path)
	}
	if debug {
		log.Printf("path[%s], line[%d]\n", *path, *line)
	}

	// get git repository directory
	pd, err := getProjectDir(*path)
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.Printf("project dir[%s]\n", pd)
	}

	// get into repository and blame
	err = blame(pd, *path, *line)
	if err != nil {
		log.Fatal(err)
	}
}

func getProjectDir(absPath string) (string, error) {
	d := filepath.Dir(absPath)
	for d != string(filepath.Separator) {
		df, err := os.Open(d)
		if err != nil {
			return "", err
		}

		dirs, err := df.Readdirnames(-1)
		if err != nil {
			df.Close()
			return "", err
		}
		df.Close()

		for _, dir := range dirs {
			if dir == ".git" {
				return d, nil
			}
		}

		d = filepath.Dir(d)
	}

	return "", errors.New("can't find git repository")
}

func blame(pp, fp string, line int) error {
	rp, err := filepath.Rel(pp, fp)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "--no-pager", "blame",
		"-L", fmt.Sprintf("%d,%d", line, line+1),
		rp)
	cmd.Dir = pp

	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	return err
}
