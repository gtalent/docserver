package main

import (
        "io/ioutil"
        "flag"
	"strings"
        "web"
	"dml"
)

var contextDir = "";

func globalServe(val string) string {
        file, err := ioutil.ReadFile("/" + val)
        if err != nil {
		return "404: File not found: dml-g/" + val
        }
	if strings.HasSuffix(val, ".dml") {
		return dml.ToHTML(val, string(file))
	}
	return string(file)
}

func contextServe(val string) string {
	if len(val) == 0 || (len(val) == 1 && val == "/") {
                val = "index.dml"
        }
        file, err := ioutil.ReadFile(contextDir + val)
        if err != nil {
		return "404: File not found: dml/" + val
        }
	if strings.HasSuffix(val, ".dml") {
		return dml.ToHTML(val, string(file))
	}
	return string(file)
}

func main() {
        global := flag.Bool("global", false, "Allow the server to access any files that the user running it has access to.")
	flag.Parse()
	if flag.NArg() != 0 {
		contextDir = flag.Arg(0)
		if !strings.HasSuffix(contextDir, "/") {
			contextDir += "/"
		}
	}
	web.Get("/dml/(.*)", contextServe)
        web.Get("/dml", contextServe)
        if *global {
		web.Get("/dml-g/(.*)", globalServe)
		web.Get("/dml-g", globalServe)
	}
        web.Run("0.0.0.0:8080")
}

