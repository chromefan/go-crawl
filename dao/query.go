package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"regexp"
	"strings"
)

type category struct {
	id    int
	title string
	ctime string
}
type crawl_log struct {
	id     int
	url    string
	ctime  string
	status int
}
type article struct {
	title   string
	content string
	author  string
	dynasty string
	tags    string
	url     string
}

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

// 正则过滤sql注入的方法
// 参数 : 要匹配的语句
func FilteredSQLInject(to_match_str string) bool {
	//过滤 ‘
	//ORACLE 注解 --  /**/
	//关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		panic(err.Error())
		return false
	}
	return re.MatchString(to_match_str)
}
