// A simple port of the core functionality of node-static as a learning exercise
package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"os"
	"path"
)

const CLI_NAME = "go-static"
const CLI_VERSION = "1.0.0"

var port int
var hostAddress string
var cache int
var directory string

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
	w.ResponseWriter.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", cache))
	w.ResponseWriter.Header().Set("Server", fmt.Sprintf("%s/%s", CLI_NAME, CLI_VERSION))
	log.Printf("[%d]: %s", status, w.Path)
	w.ResponseWriter.WriteHeader(status)
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		w := StaticResponseWriter{
			ResponseWriter: rw,
			Path:           r.URL.Path,
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	app := cli.NewApp()
	app.Name = CLI_NAME
	app.Version = CLI_VERSION
	app.Usage = "a simple port of the core functionality of github.com/cloudhead/node-static"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "TCP port at which the files will be served",
			Value:       8080,
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "host-address, a",
			Usage:       "the local network interface at which to listen",
			Value:       "127.0.0.1",
			Destination: &hostAddress,
		},
		cli.IntFlag{
			Name:        "cache, c",
			Usage:       "\"Cache-Control\" max-age header setting",
			Value:       3600,
			Destination: &cache,
		},
	}

	app.Action = func(c *cli.Context) {
		wd, _ := os.Getwd()
		directory = path.Join(wd, c.Args().First())
		addr := fmt.Sprintf("%s:%d", hostAddress, port)
		log.Printf("serving \"%s\" at http://%s\n", directory, addr)
		log.Fatal(http.ListenAndServe(addr, headerMiddleware(http.FileServer(http.Dir(directory)))))
	}

	app.Run(os.Args)
}
