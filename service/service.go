package service

import (
	"io/ioutil"
	"log"
	"os"
)

func WriteFile(file []byte) {
	ioutil.WriteFile("courses.json", file, 0644)
}

func ReadFile() []byte {
	fileContent, err := os.ReadFile("courses.json")

	if err != nil {
		log.Fatal(err)
	}

	return fileContent
}
