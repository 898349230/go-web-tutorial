package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	// 引入时会 init() 函数会运行并自我注册
	_ "github.com/go-sql-driver/mysql" // 导入包但不使用 init()
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql",
		"root:123456@tcp(localhost:3306)/shop_cloud_alibaba?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected")

	// 获取单条数据
	one, _ := getOne(106353117364359169)
	fmt.Println(one)

	price := one.totalPrice + 100

	// 更新数据
	one.updatePrice(int(price))
	one, _ = getOne(106353117364359169)
	fmt.Println(one)

	// 获取多条数据
	many, err := getMany(2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(many)

	a := order{
		id:         2,
		userId:     2,
		userName:   "lisi",
		phone:      "13789899898",
		address:    "天津",
		totalPrice: 18.9,
	}

	// 插入值
	err = a.insert()
	if err != nil {
		log.Fatal(err.Error())
	}

}

func getOne(id int64) (a order, err error) {
	a = order{}
	// QueryRow 查询返回一条数据
	err = db.QueryRow(
		"select id, t_user_id, t_user_name, t_phone, t_address, t_total_price from t_order where id= ? ", id).Scan(
		&a.id, &a.userId, &a.userName, &a.phone, &a.address, &a.totalPrice)
	return
}

func getMany(id int) (orders []order, err error) {
	rows, err := db.Query("select id, t_user_id, t_user_name, t_phone, t_address, t_total_price from t_order where id > ? ", limit)
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		a := order{}
		err = rows.Scan(&a.id, &a.userId, &a.userName, &a.phone, &a.address, &a.totalPrice)
		if err != nil {
			log.Fatal(err.Error())
		}
		orders = append(orders, a)
	}
	return
}

// 更新
func (a *order) updatePrice(price int) (err error) {
	_, error := db.Exec("update t_order set t_total_price = ? where id = ? ", price, a.id)
	if error != nil {
		log.Fatal(error.Error())
	}
	return
}

// 插入
func (a *order) insert() (err error) {
	sql := `insert into t_order 
	(id, t_user_id, t_user_name, t_phone, t_address, t_total_price)
	values (? , ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	// 需要执行
	result, err := stmt.Exec(a.id, a.userId, a.userName, a.phone, a.address, a.totalPrice)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(result)
	return
}

type order struct {
	id         int
	userId     int
	userName   string
	phone      string
	address    string
	totalPrice float32
}
