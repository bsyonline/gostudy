package db

import (
	fileUtils "alert-receiver/pkg/util"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

var db *sql.DB

func init() {
	file := fileUtils.ReadYaml("config/db.yaml")
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var ds DataSource
	err := decoder.Decode(&ds)
	if err != nil {
		log.Fatalf("read db config failed: %v\n", err)
	}
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", ds.MySQL.Username, ds.MySQL.Password, ds.MySQL.Host, ds.MySQL.Port, ds.MySQL.Database)
	fmt.Printf("url: %v\n", url)
	db, _ = sql.Open("mysql", url) //设置连接数据库的参数
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(time.Minute * 60)
	// defer db.Close() //关闭数据库
	err = db.Ping() //连接数据库
	if err != nil {
		log.Fatalf("db connect failed: %v\n", err)
	}
	log.Printf("%v\n", "db connect succeeded")
}

type DataSource struct {
	MySQL struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	}
}

func QueryAlerts() []string {
	var alert_name string
	var alerts []string
	rows, err := db.Query("select alert_name from t_alert")
	if err != nil {
		log.Fatalf("query alerts failed: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&alert_name)
		if err == sql.ErrNoRows {
			return []string{}
		}
		alerts = append(alerts, alert_name)
	}
	return alerts
}

func QueryIdleNodes() []string {
	var node_name string
	var nodes []string
	rows, err := db.Query("select node_name from t_node where node_name not in (select node_name from t_job where status = 1)")
	if err != nil {
		log.Fatalf("query idle nodes failed: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&node_name)
		if err == sql.ErrNoRows {
			return []string{}
		}
		nodes = append(nodes, node_name)
	}
	return nodes
}

func QueryNodeFromDB(ip string) string {
	var nodeName string
	err := db.QueryRow("SELECT node_name nodeName from t_node where host = ?", ip).Scan(&nodeName)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		log.Fatalf("query node failed: %v", err)
	}
	return nodeName
}

func QueryPodsFromDB(nodeName string) []Pod {
	var id int
	var pod_name string
	var node_name string
	podList := []Pod{}
	rows, err := db.Query("select a.id, a.pod_name, a.node_name from t_job a where job_name in (select job_name from t_job where node_name=?)", nodeName)
	if err != nil {
		log.Fatalf("query pods failed: %v", err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &pod_name, &node_name)
		if err != nil {
			if err == sql.ErrNoRows {
				return podList
			}
			log.Fatalf("scan error: %v", err)
		}
		var pod Pod
		pod.id = id
		pod.podName = pod_name
		pod.nodeName = node_name
		podList = append(podList, pod)
	}
	return podList
}

type Pod struct {
	id       int
	podName  string
	nodeName string
}

func ResetDB(t int) {
	db.Exec("delete from t_node")
	if t == 0 {
		db.Exec("insert into t_node (node_name, host, create_time) values ('k8s-node2', '192.168.67.204', now())")
	} else {
		sqlStr := `insert into 
					t_node (
					node_name, 
					host, 
					create_time
					)
					values
					(
					'k8s-master', 
					'192.168.67.201', 
					now()
					),  (
					'k8s-node1', 
					'192.168.67.202', 
					now()
					),  (
					'k8s-node2', 
					'192.168.67.204', 
					now()
					),  (
					'k8s-worker-4090-1.sdns.dev.cloud', 
					'36.103.180.175', 
					now()
					),  (
					'k8s-worker-4090-2.sdns.dev.cloud', 
					'36.103.180.180', 
					now()
					),  (
					'k8s-worker-4090-3.sdns.dev.cloud', 
					'36.103.180.181', 
					now()
					),  (
					'k8s-worker-4090-4.sdns.dev.cloud', 
					'36.103.180.182', 
					now()
					),  (
					'k8s-worker-4090-5.sdns.dev.cloud', 
					'36.103.180.183', 
					now()
				)`
		db.Exec(sqlStr)
	}
}
