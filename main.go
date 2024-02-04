package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//go explorer.Start(":8080")
	//rest.Start(":4000")
	if len(os.Args) < 2 {
		usage()
	}

	rest := flag.NewFlagSet("rest", flag.ExitOnError)
	portFlag := rest.Int("p", 4000, "Set the port or the server")

	switch os.Args[1] {
	case "explorer":
		fmt.Printf("Start explorer")
	case "rest":
		rest.Parse(os.Args[2:])
	default:
		usage()
	}

	if rest.Parsed() {
		fmt.Println(*portFlag)
		fmt.Println("Start REST API")
	}
}

func usage() {
	fmt.Printf("Welcome House Coin\n\n")
	fmt.Printf("Choose Application Type\n\n")
	fmt.Printf("explorer:    Start the HTML Explorer\n")
	fmt.Printf("rest:        Start the Rest API(Recommended)\n")
	os.Exit(0)
}
