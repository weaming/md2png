package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/shurcooL/github_flavored_markdown"
)

type ImageRender struct {
	BinaryPath *string
}

func fatalErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func markdown2html(markdown string) string {
	return string(github_flavored_markdown.Markdown([]byte(markdown)))
}

func (r *ImageRender) generateImage(html, format, output string, width, quality int) []byte {
	c := ImageOptions{
		BinaryPath: *r.BinaryPath,
		Input:      "-",
		Html:       html,
		Format:     format,
		Width:      width,
		Quality:    quality,
		Output:     output,
	}
	out, err := GenerateImage(&c)
	fatalErr(err)
	return out
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func readFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	fatalErr(err)
	return string(bytes)
}

func main() {
	binPath := flag.String("bin", "/usr/local/bin/wkhtmltoimage", "wkhtmltoimage bin path")
	markdownPath := flag.String("m", "", "markdown file path")
	outputPath := flag.String("o", "output.png", "output file path")
	width := flag.Int("w", 500, "output image width")
	flag.Parse()

	imgRender := ImageRender{BinaryPath: binPath}
	md := readFile(*markdownPath)
	html := markdown2html(md)
	fmt.Println(html)

	imgRender.generateImage(html, "png", *outputPath, *width, 100)
}
