package main

import (
        "web"
        "io/ioutil"
        "flag"
	"strings"
)

func globalServe(val string) string {
	val = strings.Replace(val, "dml/g", "", 1)
	if len(val) == 0 || (len(val) == 1 && val == "/") {
                val += "index.htm"
        }
        file, err := ioutil.ReadFile(val)
        if err != nil {
                err = nil
                file, err = ioutil.ReadFile(val + "l")
                if err != nil {
                        return "404: File not found."
                }
        }
        return string(file)
}

func contextServe(val string) string {
	val = strings.Replace(val, "dml", "", 1)
	if len(val) == 0 || (len(val) == 1 && val == "/") {
                val = "index.htm"
        }
        file, err := ioutil.ReadFile(val)
        if err != nil {
                err = nil
                file, err = ioutil.ReadFile("index.htm" + "l")
                if err != nil {
                        return "404: File not found."
                }
        }
        return string(file)
}

func main() {
        global := flag.Bool("global", false, "Allow the server to access any files that the user running it has access to.")
	flag.Parse()
        web.Get("/dml/(.*)", contextServe)
        web.Get("/dml", contextServe)
        if *global {
		web.Get("/dml/g/(.*)", globalServe)
        	web.Get("/dml/g", globalServe)
	}
        web.Run("0.0.0.0:8080")
}

