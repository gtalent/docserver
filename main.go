package main

import (
        "web"
        "io/ioutil"
        "flag"
	"strings"
)

/*
 * Swaps the DML tags in the given input for appropriate HTML tags in the return value.
 * Takes:
 *      doc - the text of the document to convert
 * Returns:
 *      the given text with all special characters escaped and the DML tags swapped for the appropriate HTML tags
 */
func parseDoc(doc string) string {
        return ""
}

func globalServe(val string) string {
	val = strings.Replace(val, "g-dml", "", 1)
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
        flag.Bool("global", false, "global")
	flag.Parse()
        global := false
        for i := 0; i < flag.NArg(); i++ {
                if flag.Arg(i) == "-global" {
                        global = true
                }
        }
        web.Get("/dml/(.*)", contextServe)
        web.Get("/dml", contextServe)
        if global {
		web.Get("/g-dml/(.*)", globalServe)
        	web.Get("/g-dml", contextServe)
	}
        web.Run("0.0.0.0:8080")
}

