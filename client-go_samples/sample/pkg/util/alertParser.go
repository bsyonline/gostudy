package util

import (
	"encoding/json"

	"log"
)

func HasAlert(alerts1 []string, alerts2 []string) bool {
	// 判断alerts1中是否存在alerts2的元素
	for _, alert1 := range alerts1 {
		for _, alert2 := range alerts2 {
			if alert1 == alert2 {
				return true
			}
		}
	}
	return false
}

type Alert struct {
	Alerts []struct {
		annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
		}
		Labels struct {
			Alertname string `json:"alertname"`
			Severity  string `json:"severity"`
			NodeName  string `json:"node_name"`
		} `json:"labels"`
		StartsAt string `json:"startsAt"`
	} `json:"alerts"`
}

func Parse(jsonstr string) map[string][]string {
	alertMap := make(map[string][]string)
	var alertObj Alert
	err := json.Unmarshal([]byte(jsonstr), &alertObj)
	if err != nil {
		log.Fatalf("unmarshal failed: %v", err)
	}
	// 解析
	for _, alert := range alertObj.Alerts {
		alertName := alert.Labels.Alertname
		severity := alert.Labels.Severity
		nodeName := alert.Labels.NodeName
		if nodeName != "" && severity == "critical" {
			value, has := alertMap[nodeName]
			if has {
				alertMap[nodeName] = append(value, alertName)
			} else {
				var arr []string
				arr = append(arr, alertName)
				alertMap[nodeName] = arr
			}
		}
	}
	return alertMap
}
