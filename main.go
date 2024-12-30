package main

import (
	"fmt"

	"github.com/ghf-go/glib/gcache"
)

func main() {
	c := gcache.NewFileCache("/tmp")

	fmt.Println(c.Incr("a1"))

}
