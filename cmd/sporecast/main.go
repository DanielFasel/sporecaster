package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/you/sporecast/internal/loader"
	"github.com/you/sporecast/internal/verify"
)

func main() {
	doVerify := flag.Bool("verify", false, "verify the spore against the codebase")
	flag.Parse()

	if *doVerify {
		spore, err := loader.Load("spore.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("verifying spore: %s\n\n", spore.App)
		if !verify.Run(spore, ".") {
			os.Exit(1)
		}
		return
	}

	flag.Usage()
}
