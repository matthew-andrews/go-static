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

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", cache))
		w.Header().Set("Server", fmt.Sprintf("%s/%s", CLI_NAME, CLI_VERSION))
		log.Printf("[xxx]: %s", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func serve() {
	addr := fmt.Sprintf("%s:%d", hostAddress, port)
	log.Printf("serving \"%s\" at http://%s\n", directory, addr)
	log.Fatal(http.ListenAndServe(addr, headerMiddleware(http.FileServer(http.Dir(directory)))))
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
		serve()
	}

	app.Run(os.Args)
}
