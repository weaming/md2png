package main

import (
	"fmt"
)

func renderCssPath(path string) string {
	block := "\n\n<style>\n%v\n</style>\n\n"
	return fmt.Sprintf(block, ReadFile(path))
}

func renderCssUrl(url string) string {
	block := `<link rel="stylesheet" media="all" type="text/css; charset=UTF-8" href="%v"/>`
	return fmt.Sprintf(block, url)
}
