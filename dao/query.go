package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

var tokens = make(chan struct{}, 10)

func SaveCategory(categorylist chan map[string][]string) {
	Start()
	t, _ := db.Begin()
	for cate := range categorylist {
		for i, cateUrl := range cate["url"] {

			sqlStr := fmt.Sprintf("insert into category(url,name) values(?,?)")
			fmt.Println(sqlStr, cateUrl, cate["name"][i])
			res, err := t.Exec(sqlStr,cateUrl, cate["name"][i])
			if err != nil {
				log.Fatal(err)
			}
			LastInsertId, err := res.LastInsertId()
			fmt.Printf("\n Category LastInsertId: %d \n", LastInsertId)
		}
	}
	t.Commit()
	close(categorylist)
}
func SaveLog(worklist chan []string, endChan chan bool) {
	Start()
	t, _ := db.Begin()
	for list := range worklist {
		for _, link := range list {
			sqlStr := fmt.Sprintf("insert into crawl_log(url,status) values(?,?)")
			fmt.Println(sqlStr, link, 1)
			res, err := t.Exec(sqlStr)
			if err != nil {
				log.Fatal(err)
			}
			LastInsertId, err := res.LastInsertId()
			fmt.Printf("LastInsertId: %d \n", LastInsertId)
		}
	}
	t.Commit()
	endChan <- true
}
func SaveArticle(articlelist chan map[string]string, endChan chan bool) {
	Start()
	t, _ := db.Begin()
	num := 0
	urlmap := make(map[string]bool)
	for article := range articlelist {
		num++
		tokens <- struct{}{}
		if urlmap[article["url"]] {
			continue
		}
		urlNum := findArticleNum(article["url"])
		if urlNum > 0 {
			fmt.Println(article["url"],urlNum)
			continue
		}
		sqlStr := fmt.Sprintf("insert into article(title,content,author,tags,dynasty,url) ")
		sqlStr += fmt.Sprintf(" values(?,?,?,?,?,?)")
		res, err := t.Exec(sqlStr, article["title"], article["content"], article["author"], strings.Trim(article["tags"], ","),
			article["dynasty"], article["url"])
		if err != nil {
			fmt.Println(article)
			log.Fatal(err)
		}
		LastInsertId, err := res.LastInsertId()
		fmt.Printf(" \n Article LastInsertId: %d \n", LastInsertId)
		urlmap[article["url"]] = true
		<- tokens
		if (num % 100) == 0 {
			t.Commit()
			t, _ = db.Begin()
		}
	}
	<-endChan
}

func findArticleNum(url string) int64 {
	var count int64
	err := db.QueryRow("select count(id) from article where url=$1",url).Scan(&count)
	if err != nil{
		return 0
	}
	return count
}
