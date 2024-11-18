package main

import (
	"fmt"

	"github.com/seanhuebl/blog_aggregator/internal/config"
)

func main() {
	conf := config.Read()
	conf.SetUser("Sean")
	conf = config.Read()
	fmt.Println(conf)
}
