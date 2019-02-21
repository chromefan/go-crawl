package main

import (
	"./links"
	"fmt"
	"go-crawl/dao"
	"log"
	"os"
)

var tokens = make(chan struct{}, 10)

func crawl(url string) []string {
	fmt.Println(url)
	//tokens <- struct{}{}

	list, err := links.ExtractList(url)
	//<- tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

func crawArticle(worklist chan []string,articlelist chan map[string]string)  {
	for list := range worklist {
		for _, link := range list {
			tokens <- struct{}{}
			article,err:=links.ExtractArticle(link)
			if err != nil{
				log.Fatal(err)
				break
			}
			go func() {
				fmt.Printf("%v\n", article)
				articlelist <- article
			}()
			<- tokens
		}
	}
	close(articlelist)
}
func main() {

	var worklist = make(chan []string, 1)
	var articlelist = make(chan map[string]string, 1)
	var unseenLinks = make(chan bool)

	baseUrl := os.Args[1]
	foundLinks := crawl(baseUrl)
	go func() {
		worklist <- foundLinks
		close(worklist)
	}()
	//go dao.SaveLog(worklist, unseenLinks)
	go crawArticle(worklist,articlelist)
	go dao.SaveArticle(articlelist, unseenLinks)
	<-unseenLinks
}
