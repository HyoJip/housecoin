package main

import (
	"github.com/Hyojip/housecoin/cli"
	"github.com/Hyojip/housecoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
