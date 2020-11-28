package main

import (
	"emarket/impl"
	"emarket/model"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s --web-root <path> --listen <ip:port> --data <path>\n", os.Args[0])
		os.Exit(1)
	}

	webRootOpt := flag.String("web-root", "", "<path>")
	listenOpt := flag.String("listen", "", "<ip:port>")
	dataOpt := flag.String("data", "", "<path>")
	flag.Parse()

	if webRootOpt == nil || *webRootOpt == "" {
		fmt.Println("web root not specified")
		os.Exit(1)
	}

	if listenOpt == nil || *listenOpt == "" {
		fmt.Println("listen ip:port not specified")
		os.Exit(1)
	}

	if dataOpt == nil || *dataOpt == "" {
		fmt.Println("listen ip:port not specified")
		os.Exit(1)
	}

	dataFile, err := filepath.Abs(*dataOpt)
	webRoot, err := filepath.Abs(*webRootOpt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	emarket, err := impl.NewEMarket(webRoot, model.NewDB(dataFile))

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      emarket,
		Addr:         *listenOpt,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("started %v\n", dir)
	log.Fatal(srv.ListenAndServe())
}