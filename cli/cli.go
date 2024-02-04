package cli

import (
	"flag"
	"fmt"
	"github.com/Hyojip/housecoin/explorer"
	"github.com/Hyojip/housecoin/rest"
	"os"
	"strconv"
)

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	mode := flag.String("m", "rest", "Choose Between 'rest' and 'explorer'")
	iPort := flag.Int("p", 4000, "Set the port of the server")
	flag.Parse()
	port := ":" + strconv.Itoa(*iPort)

	switch *mode {
	case "explorer":
		explorer.Start(port)
	case "rest":
		rest.Start(port)
	default:
		usage()
	}

}

func usage() {
	fmt.Printf("Welcome House Coin\n\n")
	fmt.Printf("Choose Application Mode\n\n")
	fmt.Printf("-m:    Choose Between 'rest' and 'explorer'\n")
	fmt.Printf("-p:    Set the port of the server\n")
	os.Exit(0)
}
