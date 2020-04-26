package main

import (
	"fmt"

	"github.com/ravenscroftj/papersleuth/sleuth"
)

func main() {
	fmt.Println("Hello")

	//client, err := sleuth.GetDefaultClient()
	client, err := sleuth.GetDefaultCrossrefClient()

	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.GetWorkByDOI("10.1192/bjp.bp.107.044677")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Paper title %s,abstract: %s\n", resp.Title[0], resp.Abstract)
	}

}
