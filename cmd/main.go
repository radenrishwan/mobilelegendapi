package main

import (
	"context"
	"log"

	"github.com/radenrishwan/mobilelegendapi"
)

const (
	ENDPOINT = "https://m.mobilelegends.com/id/news/EVENTS/list"
)

func main() {
	core := mobilelegendapi.NewCore(context.Background())

	news := core.GetNews(context.Background())
	// print news
	for _, n := range news {
		log.Println(n)
	}
}
