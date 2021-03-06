# go-static

A simple go port of the core functionality of the extremely useful [node-static](https://github.com/cloudhead/node-static) command line tool.  Just a little exercise to learn a bit more Go.

## Installation

Just `go get` this.

```sh
go get github.com/matthew-andrews/go-static
```

## Example usage

```sh
# serve up the current directory
$ go-static
2016/03/06 11:28:50 serving "{{current directory}}" at http://127.0.0.1:8080

# serve up a different directory
$ go-static public
serving "{{current directory}}/public" at http://127.0.0.1:8080

# specify additional headers (this one is useful for development)
$ go-static -H '{"Cache-Control": "no-cache, must-revalidate"}'
serving "{{current directory}}" at http://127.0.0.1:8080

# set cache control max age
$ go-static -c 7200
serving "{{current directory}}" at http://127.0.0.1:8080

# expose the server to your local network
$ go-static -a 0.0.0.0
serving "{{current directory}}" at http://0.0.0.0:8080

# show help message, including all options
$ go-static -h
```
