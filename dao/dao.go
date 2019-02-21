package dao


import (
	"database/sql"
	"fmt"
	"github.com/aWildProgrammer/fconf"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"reflect"
	"time"
)
var (
	db     = &sql.DB{}

)

func Start(){
	c , err := fconf.NewFileConf(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	dbHost := c.String("mysql.dbHost")
	dbUser := c.String("mysql.dbUser")
	dbPass :=  c.String("mysql.dbPass")
	dbName := c.String("mysql.dbName")

	db, _ = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
}

func FindAll(sqlStr string) {
	//方式3 query
	tx, _ := db.Begin()
	defer tx.Commit()
	rows, _ := tx.Query(sqlStr)
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err !=nil{
			log.Fatal(err)
		}
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
}

func random() int {
	rand.Seed(time.Now().UnixNano())
	num := rand.Int()
	return num
}
func Insert(sqlStr string) int64 {
	t, _ := db.Begin()
	res, err := t.Exec(sqlStr, sqlStr)
	if err !=nil{
		log.Fatal(err)
	}
	LastInsertId, err := res.LastInsertId()
	t.Commit()
	return LastInsertId
}
func getInsert(q interface{}) {
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		query := fmt.Sprintf("insert into %s  values(", t)
		v := reflect.ValueOf(q)

		for i := 0; i < v.NumField(); i++ {
			fmt.Println(v.Index(i))
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %d", query, v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, \"%s\"", query, v.Field(i).String())
				}
			default:
				fmt.Println("Unsupported type")
				return
			}
		}
		query = fmt.Sprintf("%s)", query)
		fmt.Println(query)
		return

	}
	fmt.Println("unsupported type")
}
