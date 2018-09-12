package main

import (
	"flag"
	"fmt"
	"path"

	"github.com/golang-commonmark/markdown"
)

var cssUrlList ArrayFlags
var cssFileList ArrayFlags

func main() {
	binPath := flag.String("bin", "/usr/local/bin/wkhtmltoimage", "wkhtmltoimage bin path")
	markdownPath := flag.String("m", "", "markdown file path")
	outputPath := flag.String("o", "", "output file path (default same as markdown file name)")
	htmlFile := flag.String("html", "", "file path of HTML content. CSS relative options will ignore if include-css flag is not true")

	width := flag.Int("w", 960, "output image width")
	quality := flag.Int("q", 80, "output image quality, maxium is 100")

	flag.Var(&cssUrlList, "cssurl", "CSS URLs [repeatable, optional]")
	flag.Var(&cssFileList, "cssfile", "CSS file path, support any style you like❤️ , include fonts! [repeatable, optional]")
	cssName := flag.String("cssname", "", "use builtin CSS from github.com/mixu/markdown-styles:"+cssListHelpText)
	//staticPath := flag.String("static", ".", "static files path")

	print := flag.Bool("print", false, "print generated html")
	includeCSS := flag.Bool("include-css", false, "include css file as html header when read content from *.html file")
	flag.Parse()

	if *cssName != "" {
		cssUrlList = append(cssUrlList, getCssUrl(*cssName))
	}

	if *outputPath == "" {
		*outputPath = ReplaceExt(path.Base(*markdownPath), "png")
	}

	//prepare static files
	//go staticServer(*staticPath)

	header := ""
	if *htmlFile == "" || *includeCSS {
		header = `<meta http-equiv="content-Type" content="text/html; charset=UTF-8" />`
		for _, f := range cssFileList {
			header += renderCssPath(f)
		}
		for _, u := range cssUrlList {
			header += renderCssUrl(u)
		}
	}

	html := ""
	if *htmlFile == "" {
		md := ReadFile(*markdownPath)
		html = fmt.Sprintf("%v\n\n<div class='main container content article'>\n%v\n</div>", header, markdown2html(md))
	} else {
		if *includeCSS {
			html = header + ReadFile(*htmlFile)
		} else {
			html = ReadFile(*htmlFile)
		}
	}

	// renderHTML
	if *print {
		fmt.Println(html)
	}

	imgRender := ImageRender{BinaryPath: binPath}
	imgRender.generateImage(html, "png", *outputPath, *width, *quality)
}

type ImageRender struct {
	BinaryPath *string
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

// builtin css

var cssListHelpText = `
	jasonm23-dark
	jasonm23-foghorn
	jasonm23-markdown
	jasonm23-swiss
	markedapp-byword
	thomasf-solarizedcssdark
	thomasf-solarizedcsslight
`

func getCssUrl(name string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/mixu/markdown-styles/master/output/%v/assets/style.css", name)
}
