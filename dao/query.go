package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func SaveLog(worklist chan []string, unseenLinks chan bool) {
	Start()
	t, _ := db.Begin()
	for list := range worklist {
		for _, link := range list {
			sqlStr := fmt.Sprintf("insert into crawl_log(url,status) values(?,?)", link)
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
	unseenLinks <- true
}
func SaveArticle(articlelist chan map[string]string, unseenLinks chan bool) {
	Start()
	t, _ := db.Begin()
	for article := range articlelist {
		sqlStr := fmt.Sprintf("insert into article(title,content,author,tags,dynasty,url) ")
		sqlStr += fmt.Sprintf(" values(?,?,?,?,?,?)")
		res, err := t.Exec(sqlStr, article["title"], article["content"], article["author"], strings.Trim(article["tags"], ","),
			article["dynasty"], article["url"])
		if err != nil {
			log.Fatal(err)
		}
		LastInsertId, err := res.LastInsertId()
		fmt.Printf("LastInsertId: %d \n", LastInsertId)
	}
	t.Commit()
	unseenLinks <- true
}