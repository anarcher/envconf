package main

import (
	"fmt"
	"github.com/namsral/flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	CONF_FILES = "ENV_CONF_FILES"
)

var (
	templateFuncMap template.FuncMap
)

func init() {

	flag.Usage = func() {
		fmt.Printf("Usage: envconv [options] file [files...]\n\n")
		flag.PrintDefaults()
	}

	version := flag.Bool("v", false, "Prints current version")

	flag.Parse()

	if *version {
		fmt.Println("envconf version " + VERSION)
		os.Exit(0)
	}

	templateFuncMap = template.FuncMap{
		"default": func(args ...string) string {
			defer recovery()
			return args[0]
		},
		"in": func(arg string, slice []string) bool {
			defer recovery()
			for _, i := range slice {
				if arg == i {
					return true
				}
			}
			return false
		},
	}
}

func main() {
	paths := confFilePaths()
	env := environ()
	for _, path := range paths {
		err := transformFile(path, env)
		if err != nil {
			log.Printf("Err: %s,%s", path, err)
		} else {
			log.Printf("Transformed %s", path)
		}
	}
}

func environ() map[string]interface{} {
	env := make(map[string]interface{})
	for _, pair := range os.Environ() {
		item := strings.SplitN(pair, "=", 2)
		if len(item) < 2 {
			log.Fatalf("invalid item from os.Environ %s", pair)
		}
		if strings.Contains(item[1], ",") {
			val := strings.Split(item[1], ",")
			env[item[0]] = val
		} else {
			env[item[0]] = item[1]
		}
	}
	return env
}

func confFilePaths() []string {
	var paths []string

	if len(flag.Args()) >= 0 {
		for _, path := range flag.Args() {
			paths = append(paths, path)
		}
	}

	if val := os.Getenv(CONF_FILES); val != "" {
		for _, path := range strings.Split(val, ":") {
			paths = append(paths, path)
		}
	}
	return paths
}

func transformFile(path string, env map[string]interface{}) error {
	tname := filepath.Base(path)
	t := template.New(tname).Funcs(templateFuncMap)
	t, err := t.ParseFiles(path)
	if err != nil {
		log.Printf("Err: Parsing file %s %s", path, err)
		return err
	}

	var fileMode os.FileMode = 0644 //default
	if fi, err := os.Stat(path); err == nil {
		fileMode = fi.Mode()
	} else {
		log.Printf("Err: %s", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		log.Printf("Err: Can't open a file %s %s", path, err)
		return err
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	if err = t.Execute(file, env); err != nil {
		return err
	}

	return nil
}

func recovery() {
	if r := recover(); r != nil {
		log.Println("Panic in", r)
	}
}
