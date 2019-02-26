package main

import (
	"./links"
	"fmt"
	"go-crawl/dao"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var tokens = make(chan struct{}, 10)

func crawl(url string) []string {
	fmt.Println(url)
	//tokens <- struct{}{}
	list, err := links.ExtractList(url)
	//<-tokens
	if err != nil {
		log.Print(err)
	}
	return list
}
func crawCategory(baseUrl string, worklist chan []string) {
	//var categorylist = make(chan map[string][]string,1)
	category, _ := links.ExtractCategory(baseUrl)
	for _, cateUrl := range category["url"] {
		foundLinks := crawl(cateUrl)
		worklist <- foundLinks
	}

	/*go func() {
		categorylist <- category
	}()
	go dao.SaveCategory(categorylist)
	//go func() { <-categorylist}()
	close(worklist)*/
}
func crawArticle(worklist chan []string, articlelist chan map[string]string) {
	num := 0
	for list := range worklist {
		for _, link := range list {
			num++
			tokens <- struct{}{}
			article, err := links.ExtractArticle(link)
			if err != nil {
				log.Fatal(err)
				break
			}
			go func() {
				fmt.Printf("%v\n", article)
				articlelist <- article
			}()
			<-tokens
			sec := rand.Intn(9)+2
			if num % 30 == 0{
				time.Sleep(time.Duration(sec)*time.Second)
			}
		}
	}
	close(articlelist)
}
func main() {

	var worklist = make(chan []string, 20)
	var articlelist = make(chan map[string]string,10)

	var endChan = make(chan bool)

	baseUrl := os.Args[1]

	go crawCategory(baseUrl, worklist)
	//go dao.SaveLog(worklist, unseenLinks)
	go crawArticle(worklist, articlelist)
	go dao.SaveArticle(articlelist, endChan)

	select {
	case <-endChan:
		fmt.Println("\n—-----—done—-----—")
		return
	}
}
func tracefile(str_content string) {
	fd, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join([]string{fd_time, " : ", str_content, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}
