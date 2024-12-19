package main

import (
	"fmt"
	"log"

	"github.com/feelbeatapp/feelbeatserver/internal/thirdparty/ytsearch"
)

func main() {
	result, err := ytsearch.Search("System of a down Aerials")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range result {
		fmt.Printf("%v\n%v\n%v\n\n", r.VideoId, r.Title, r.Duration)
	}
}
