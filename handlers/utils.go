package handlers

import (
	"bufio"
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

func loadConf(configure_path string) map[interface{}]interface{} {
	file, err := os.Open(configure_path)
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		log.Panic(err.Error())
	}

	bytes := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bytes)
	if err != nil {
		log.Panic(err.Error())
	}
	results := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(bytes), &results)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return results
}
