/*
   Copyright 2010-2014 gtalent2@gmail.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"flag"
	"github.com/hoisie/web"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"strings"
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
		return "404: File not found: doc-g/" + val
	}
	if strings.HasSuffix(val, ".md") {
		return string(blackfriday.MarkdownBasic(file))
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
		val = "index.md"
	}
	file, err := ioutil.ReadFile(contextDir + val)
	if err != nil {
		return "404: File not found: doc/" + val
	}
	if strings.HasSuffix(val, ".md") {
		return string(blackfriday.MarkdownBasic(file))
	}
	return string(file)
}

func main() {
	global := flag.Bool("global", false, "Allow the server to access any files that the user running it has access to.")
	remote := flag.Bool("remote", false, "Allow the remote clients to access the server.")
	port := "15448"
	settingsFile, err := ioutil.ReadFile("docserver.conf")
	if err == nil {
		settings := strings.Split(string(settingsFile), "\n")
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
			} else if strings.HasPrefix(settings[i], "AllowRemote:") {
				g := strings.Trim(strings.Replace(settings[i], "AllowRemote:", "", 1), "\t ")
				if strings.ToLower(g) == "true" {
					*remote = true
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
	web.Get("/doc/(.*)", contextServe)
	web.Get("/doc", contextServe)
	if *global {
		web.Get("/doc-g/(.*)", globalServe)
		web.Get("/doc-g", globalServe)
	}
	if *remote {
		web.Run("0.0.0.0:" + port)
	} else {
		web.Run("127.0.0.1:" + port)
	}
}
