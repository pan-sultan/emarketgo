package main

import (
	"emarket/cmd/emarket/cmd"
	"emarket/internal/emarket/file"
	"emarket/internal/emarket/html/page"
	"fmt"
	"log"
)

func main() {
	args := cmd.Parse()
	magazService := file.NewMagazineService(args.DataFile)
	magazines, err := magazService.Find()

	if err != nil {
		log.Fatalln(err)
	}

	for _, p := range page.MagazineList(magazines) {
		fmt.Println(string(p))
	}
}
