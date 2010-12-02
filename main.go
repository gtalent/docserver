package main

import (
	"io/ioutil"
	"flag"
	"strings"
	"web"
	"dml"
)

var contextDir = ""

func globalServe(params *web.Context, val string) string {
	cssPath := params.Params["style"]
	if len(cssPath) == 0 {
		cssPath = "default.css"
	} else {
		cssPath += ".css"
	}
	file, err := ioutil.ReadFile("/" + val)
	if err != nil {
		return "404: File not found: dml-g/" + val
	}
	if strings.HasSuffix(val, ".dml") {
		css, err := ioutil.ReadFile(contextDir + cssPath)
		if err == nil {
			return dml.ToHTMLCSS(val, string(file), string(css))
		}
		return dml.ToHTML(val, string(file))
	}
	return string(file)
}

func contextServe(params *web.Context, val string) string {
	cssPath := params.Params["style"]
	if len(cssPath) == 0 {
		cssPath = "default.css"
	} else {
		cssPath += ".css"
	}
	if len(val) == 0 || (len(val) == 1 && val == "/") {
		val = "index.dml"
	}
	file, err := ioutil.ReadFile(contextDir + val)
	if err != nil {
		return "404: File not found: dml/" + val
	}
	if strings.HasSuffix(val, ".dml") {
		css, err := ioutil.ReadFile(contextDir + cssPath)
		if err == nil {
			return dml.ToHTMLCSS(val, string(file), string(css))
		}
		return dml.ToHTML(val, string(file))
	}
	return string(file)
}

func main() {
	global := flag.Bool("global", false, "Allow the server to access any files that the user running it has access to.")
	port := "15448"
	settingsFile, err := ioutil.ReadFile("dmld.conf")
	if err == nil {
		settings := strings.Split(string(settingsFile), "\n", -1)
		for i := 0; i < len(settings); i++ {
			if strings.HasPrefix(settings[i], "Port:") {
				port = strings.Trim(strings.Replace(settings[i], "Port:", "", 1), "\t ")
			} else if strings.HasPrefix(settings[i], "Context:") {
				contextDir = strings.Trim(strings.Replace(settings[i], "Context:", "", 1), "\t ")
			} else if strings.HasPrefix(settings[i], "Global:") {
				g := strings.Trim(strings.Replace(settings[i], "Global:", "", 1), "\t ")
				if strings.ToLower(g) == "true" {
					*global = true
				}
			}
		}
	}
	flag.Parse()
	//read the context from the input and override whats in the settings file if something was there
	if flag.NArg() != 0 {
		contextDir = flag.Arg(0)
	}
	//make sure the context is a directory
	if len(contextDir) != 0 && !strings.HasSuffix(contextDir, "/") {
		contextDir += "/"
	}
	dml.SetDefaultCSSPath(contextDir + "default.css")
	web.Get("/dml/(.*)", contextServe)
	web.Get("/dml", contextServe)
	if *global {
		web.Get("/dml-g/(.*)", globalServe)
		web.Get("/dml-g", globalServe)
	}
	web.Run("0.0.0.0:" + port)
}
