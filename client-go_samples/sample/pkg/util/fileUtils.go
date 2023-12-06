package util

import (
	"log"
	"os"
)

func ReadYaml(dbConfigFile string) *os.File {
	// 打开文件
	file, err := os.Open(dbConfigFile)
	if err != nil {
		log.Fatalf("open yaml file failed: %v\n", err)
	}
	return file
}
