package main

import (
	"io/ioutil"
	"os"
	"path"
	"log"
)

func ExistFile(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func ReadFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	fatalErr(err)
	return string(bytes)
}

func ReplaceExt(name, newExt string) string {
	ext := path.Ext(name)
	return name[0:len(name)-len(ext)] + "." + newExt
}

func fatalErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
