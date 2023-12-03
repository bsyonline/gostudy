package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, _ := sql.Open("mysql", "root:123456@(192.168.67.204:30306)/test") // 设置连接数据库的参数
	defer db.Close()                                                      //关闭数据库
	err := db.Ping()                                                      //连接数据库
	if err != nil {
		fmt.Printf("%v\n", "数据库连接失败")
		return
	}
	fmt.Printf("%v\n", "数据库连接成功")

	var id int
	var podName string
	var node string

	rows, _ := db.Query("select id,pod_name podName,node from t_pod where node='k8s-node2'")
	for rows.Next() {
		rows.Scan(&id, &podName, &node)
		fmt.Println(id, "--", podName)
	}

	fmt.Println("--")
	// rows, _ := db.Query("select * from t1")
	// for rows.Next() {
	// 	rows.Scan(&id, &podName, &node)
	// 	fmt.Println(id, "--", podName)
	// }

}
