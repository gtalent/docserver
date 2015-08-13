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
	"log"
	"os"
	"strings"
)

func dirList(dir string) string {
	out := `<html>
	<head>
		<title>` + dir + `</title>
	</head>
	<body>
	<h3>` + dir + `</h3><ul>`
	list, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, v := range list {
			name := v.Name()
			if v.IsDir() || name[len(name)-3:] == ".md" {
				out += "<li><a href=\"" + dir + "/" + name + "\">" + name + "</a></li>"
			}
		}
	} else {
		log.Println("error:", err)
		return "404: Not Found"
	}
	out += `</ul>
	</body>
</html>`
	return out
}

func mkServer(contextDir string, format bool) func(*web.Context, string) string {
	return func(params *web.Context, val string) string {
		if len(val) == 0 || (len(val) == 1 && val == "/") {
			val = "."
		}
		fullPath := contextDir + val
		if fi, err := os.Stat(fullPath); err == nil && fi.IsDir() {
			return dirList(fullPath)
		} else if err != nil {
			log.Println("error:", err)
			return "404: File not found: " + val
		} else {
			file, err := ioutil.ReadFile(fullPath)
			if err != nil {
				log.Println("error:", err)
				return "404: File not found: " + val
			} else if strings.HasSuffix(val, ".md") {
				var text string
				if format {
					text = string(blackfriday.MarkdownBasic(file))
				} else {
					text = string(file)
				}
				return `<html>
	<head>
		<title>` + val + `</title>
	</head>
	<body>
` + text + `
	</body>
</html>`
			}
			return string(file)
		}
	}
}

func main() {
	var contextDir = ""
	global := flag.Bool("global", false, "Allow the server to access any files that the user running it has access to.")
	remote := flag.Bool("remote", false, "Allow the remote clients to access the server.")
	port := "15448"
	flag.Parse()
	//read the context from the input and override whats in the settings file if something was there
	if flag.NArg() != 0 {
		contextDir = flag.Arg(0)
	}
	//make sure the context is a directory
	if len(contextDir) != 0 && !strings.HasSuffix(contextDir, "/") {
		contextDir += "/"
	}

	contextServe := mkServer(contextDir, true)
	rawServe := mkServer(contextDir, false)
	globalServe := mkServer("", true)
	web.Get("/doc/(.*)", contextServe)
	web.Get("/doc", contextServe)
	web.Get("/raw/(.*)", rawServe)
	web.Get("/raw", rawServe)
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
