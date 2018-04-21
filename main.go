package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/golang-commonmark/markdown"
)

type ImageRender struct {
	BinaryPath *string
}

func fatalErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func markdown2html(text string) string {
	// return string(github_flavored_markdown.Markdown([]byte(markdown)))
	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	return md.RenderToString([]byte(text))
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

func replaceExt(name, newExt string) string {
	ext := path.Ext(name)
	return name[0:len(name)-len(ext)] + "." + newExt
}

func main() {
	binPath := flag.String("bin", "/usr/local/bin/wkhtmltoimage", "wkhtmltoimage bin path")
	markdownPath := flag.String("m", "", "markdown file path")
	outputPath := flag.String("o", "", "output file path (default same as markdown file name)")
	width := flag.Int("w", 960, "output image width")
	quality := flag.Int("q", 80, "output image quality, maxium is 100")
	cssPath := flag.String("css", "", "optional css file path, support any style you like❤️, include fonts!")
	// staticPath := flag.String("static", ".", "static files path")

	debug := flag.Bool("debug", false, "show generated html")
	flag.Parse()

	if *outputPath == "" {
		*outputPath = replaceExt(path.Base(*markdownPath), "png")
	}

	// prepare static files
	// go staticServer(*staticPath)

	imgRender := ImageRender{BinaryPath: binPath}
	md := readFile(*markdownPath)
	html := markdown2html(md)
	if *cssPath != "" {
		html = renderCSS(*cssPath) + html
	}

	if *debug {
		fmt.Println(html)
	}

	imgRender.generateImage(html, "png", *outputPath, *width, *quality)
}
