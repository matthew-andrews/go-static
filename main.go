// A simple port of the core functionality of node-static as a learning exercise
package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"os"
	"path"
)

const CLI_NAME = "go-static"
const CLI_VERSION = "1.0.0"

var headers map[string]string

// Define a ‘StaticResponseWriter’ in order do some manipulation and logging
// of the response immediately before the header is flushed in ‘WriteHeader’
type StaticResponseWriter struct {
	ResponseWriter http.ResponseWriter
	Path           string
}

func (w StaticResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w StaticResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func (w StaticResponseWriter) WriteHeader(status int) {
	for header, value := range headers {
		w.ResponseWriter.Header().Set(header, value)
	}
	log.Printf("[%d]: %s", status, w.Path)
	w.ResponseWriter.WriteHeader(status)
}

func main() {
	app := cli.NewApp()
	app.Name = CLI_NAME
	app.Version = CLI_VERSION
	app.Usage = "a simple port of the core functionality of github.com/cloudhead/node-static"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Usage: "TCP port at which the files will be served",
			Value: 8080,
		},
		cli.StringFlag{
			Name:  "host-address, a",
			Usage: "the local network interface at which to listen",
			Value: "127.0.0.1",
		},
		cli.IntFlag{
			Name:  "cache, c",
			Usage: "\"Cache-Control\" max-age header setting",
			Value: 3600,
		},
		cli.StringFlag{
			Name:  "headers, H",
			Usage: "additional headers (in JSON format)",
			Value: "{}",
		},
	}

	app.Action = func(c *cli.Context) {
		headers = make(map[string]string)
		headers["Server"] = fmt.Sprintf("%s/%s", CLI_NAME, CLI_VERSION)
		headers["Cache-Control"] = fmt.Sprintf("max-age=%d", c.GlobalInt("cache"))
		err := json.Unmarshal([]byte(c.GlobalString("headers")), &headers)
		if err != nil {
			log.Fatal("Headers option is invalid JSON")
		}

		wd, _ := os.Getwd()
		directory := path.Join(wd, c.Args().First())
		addr := fmt.Sprintf("%s:%d", c.GlobalString("host-address"), c.GlobalInt("port"))
		log.Printf("serving \"%s\" at http://%s\n", directory, addr)
		fileServer := http.FileServer(http.Dir(directory))
		log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			w := StaticResponseWriter{
				ResponseWriter: rw,
				Path:           r.URL.Path,
			}
			fileServer.ServeHTTP(w, r)
		})))
	}

	app.Run(os.Args)
}
