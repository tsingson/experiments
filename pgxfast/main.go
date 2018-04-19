package main

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
)

var pool *pgx.ConnPool

type Row struct {
	Id         int    `json:"id"`
	Createtime string `json:"createtime"`
	Goods_name string `json:"goods_name"`
	Nickname   string `json:"nickname"`
	Mobileno   string `json:"mobileno"`
}

type Data struct {
	Return_code string `json:"return_code"`
	Rows        []Row  `json:"data"`
}

func httpHandle(w *fasthttp.RequestCtx) {
	offset := string(w.QueryArgs().Peek("offset"))
	if offset == "" {
		offset = "0"
	}

	sql := `
    SELECT 
          orders.id,orders.createtime::text,
          orders.goods_name,
          users.nickname,users.mobileno 
    FROM
          orders 
          INNER JOIN users ON orders.users_id=users.id
    ORDER BY 
          orders.id DESC
    OFFSET
          ` + offset + `
    LIMIT 10
    `

	rows, err := pool.Query(sql)
	checkErr(err)
	defer rows.Close()
	w.SetContentType("text/html")

	var data Data = Data{}
	data.Rows = make([]Row, 0)
	data.Return_code = "FAIL"

	for rows.Next() {
		var row Row
		err = rows.Scan(&row.Id, &row.Createtime, &row.Goods_name, &row.Nickname, &row.Mobileno)
		checkErr(err)
		data.Rows = append(data.Rows, row)
	}

	if len(data.Rows) > 0 {
		data.Return_code = "SUCCESS"
	}

	ret, _ := json.Marshal(data)
	fmt.Fprintf(w, "%s", string(ret))
}

func main() {
	var err error
	poolnum := 8
	checkErr(err)
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     "postgres",
			Password: "postgres",
			Database: "car_goods_matching",
			Port:     9410,
		},
		MaxConnections: poolnum,
	}

	pool, err = pgx.NewConnPool(connPoolConfig)
	checkErr(err)

	if err := fasthttp.ListenAndServe("0.0.0.0:8091", httpHandle); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
