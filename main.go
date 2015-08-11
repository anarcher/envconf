package main

import (
	"github.com/namsral/flag"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	CONF_FILES = "ENV_CONF_FILES"
)

var (
	env             map[string]string = make(map[string]string)
	confPaths       []string
	templateFuncMap template.FuncMap
)

func init() {
	flag.Parse()

	if len(flag.Args()) >= 0 {
		for _, path := range flag.Args() {
			confPaths = append(confPaths, path)
		}
	}

	for _, pair := range os.Environ() {
		item := strings.SplitN(pair, "=", 2)
		if len(item) < 2 {
			log.Fatalf("invalid item from os.Environ %s", pair)
		}
		env[item[0]] = item[1]
	}

	if val, ok := env[CONF_FILES]; ok {
		paths := strings.Split(val, ":")
		for _, path := range paths {
			confPaths = append(confPaths, path)
		}
	}

	templateFuncMap = template.FuncMap{
		"default": func(args ...string) string {
			return args[0]
		},
	}
}

func main() {
	for _, path := range confPaths {
		err := transformFile(path)
		if err != nil {
			log.Printf("Err: %s,%s", path, err)
		} else {
			log.Printf("Transformed %s", path)
		}
	}
}

func transformFile(path string) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Err: Parsing file %s %s", path, err)
		return err
	}
	t = t.Funcs(templateFuncMap)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
