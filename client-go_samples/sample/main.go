package main

import (
	"alert-receiver/pkg/db"
	kube "alert-receiver/pkg/k8s"
	alertParser "alert-receiver/pkg/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/webhook", func(c *gin.Context) {
		param := make(map[string]interface{})
		c.BindJSON(&param)
		paramBytes, _ := json.Marshal(&param)
		log.Printf("=========param:%v", string(paramBytes))
		paramJson := string(paramBytes)
		// 根据告警类型进行处理
		// {"alerts":[{"annotations":{"description":"test gpu cores","summary":"test gpu cores"},"labels":{"alertname":"TestGpuCore","node_name":"k8s-worker-4090-1.sdns.dev.cloud","severity":"critical"},"startsAt":"2023-12-05T08:21:49.049Z","status":"firing"},{"annotations":{"description":"test gpu cores","summary":"test gpu cores"},"labels":{"alertname":"TestGpuCore","node_name":"k8s-worker-4090-1.sdns.dev.cloud","severity":"critical"},"startsAt":"2023-12-05T08:21:49.049Z","status":"firing"}]}
		paramJson = `{"alerts":[{"annotations":{"description":"test gpu cores","summary":"test gpu cores"},"labels":{"alertname":"TestGpuCore","node_name":"k8s-node2","severity":"critical"},"startsAt":"2023-12-05T08:21:49.049Z","status":"firing"},{"annotations":{"description":"test gpu cores","summary":"test gpu cores"},"labels":{"alertname":"TestGpuCore","node_name":"k8s-worker-4090-1.sdns.dev.cloud","severity":"critical"},"startsAt":"2023-12-05T08:21:49.049Z","status":"firing"}]}`
		alertMap := alertParser.Parse(paramJson)
		alertsNeedDoRestart := db.QueryAlerts()
		log.Printf("need dealwith alerts: %v\n", alertsNeedDoRestart)
		for nodeName, alertList := range alertMap {
			fmt.Println(nodeName, alertList)
			if alertParser.HasAlert(alertList, alertsNeedDoRestart) {
				idleNodes := db.QueryIdleNodes()
				hasIdleNode := false
				if len(idleNodes) > 0 {
					hasIdleNode = true
				}
				log.Printf("idleNodes: %v\n", len(idleNodes))
				dealwith(nodeName, hasIdleNode)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": nil,
		})
	})
	r.POST("/reset", func(c *gin.Context) {
		param := make(map[string]interface{})
		c.BindJSON(&param)
		t := param["type"].(float64)
		reset()
		db.ResetDB(int(t))
		log.Println("reset")
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": nil,
		})
	})
	r.Run(":5000")
}

func dealwith(nodeName string, hasIdleNode bool) {
	configPath := "./config/config"
	clientset := kube.CreateClient(configPath)
	dd := kube.CreateDynamicClient(configPath)
	// 查询node信息
	node := kube.GetNode(clientset, nodeName)
	if node == nil {
		return
	}
	// 给node打taint
	key := "err"
	value := "gpu-disabled"
	kube.AddTaint(clientset, *node, key, value)
	kube.RestartJob(clientset, dd, hasIdleNode)
}

func reset() {
	configPath := "./config/config"
	clientset := kube.CreateClient(configPath)
	dd := kube.CreateDynamicClient(configPath)
	node := kube.GetNode(clientset, "k8s-node2")
	// 给node打taint
	key := "err"
	value := "gpu-disabled"
	kube.DelTaint(clientset, *node, key, value)
	kube.RestartJob(clientset, dd, true)
}
