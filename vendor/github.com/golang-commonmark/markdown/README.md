markdown [![GoDoc](http://godoc.org/github.com/golang-commonmark/markdown?status.svg)](http://godoc.org/github.com/golang-commonmark/markdown) [![License](https://img.shields.io/badge/licence-BSD--2--Clause-blue.svg)](https://opensource.org/licenses/BSD-2-Clause) [![Build Status](https://drone.io/github.com/golang-commonmark/markdown/status.png)](https://drone.io/github.com/golang-commonmark/markdown/latest)
========

golang-commonmark/markdown package provides CommonMark-compliant markdown parser and renderer, written in Go.

## Installation

    go get github.com/golang-commonmark/markdown

You can also go get [mdtool](https://github.com/golang-commonmark/mdtool), an example command-line tool:

    go get github.com/golang-commonmark/mdtool

## Standards support

Currently supported CommonMark spec: [v0.20](http://spec.commonmark.org/0.20/).

## Extensions

Besides the features required by CommonMark, golang-commonmark/markdown supports:

  * Tables (GFM)
  * Strikethrough (GFM)
  * Autoconverting plain-text URLs to links
  * Typographic replacements (smart quotes and other)

## Usage

``` go
md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
fmt.Println(md.RenderToString([]byte("Header\n===\nText")))
```

Check out [the source of mdtool](https://github.com/golang-commonmark/mdtool/blob/master/main.go) for a more complete example.

The following options are currently supported:

  Name            |  Type  |                        Description                          | Default
  --------------- | ------ | ----------------------------------------------------------- | ---------
  HTML            | bool   | whether to enable raw HTML                                  | false
  Tables          | bool   | whether to enable GFM tables                                | true
  Linkify         | bool   | whether to autoconvert plain-text URLs to links             | true
  Typographer     | bool   | whether to enable typographic replacements                  | true
  Quotes          | string | double + single quote replacement pairs for the typographer | “”‘’
  MaxNesting      | int    | maximum nesting level                                       | 20
  LangPrefix      | string | CSS language prefix for fenced blocks                       | language-
  Breaks          | bool   | whether to convert newlines inside paragraphs into `<br>`   | false
  Nofollow        | bool   | whether to add `rel="nofollow"` to links                    | false
  XHTMLOutput     | bool   | whether to output XHTML instead of HTML                     | false

## Benchmarks

Rendering spec/spec-0.20.txt on a Intel(R) Core(TM) i5-2400 CPU @ 3.10GHz

    BenchmarkRenderSpecNoHTML         300     5751062 ns/op    2696668 B/op     7861 allocs/op
    BenchmarkRenderSpec               100    15423729 ns/op    4659899 B/op    38725 allocs/op
    BenchmarkRenderSpecBlackFriday    200     7531597 ns/op    2750346 B/op    37247 allocs/op

## TODO

  * Improve performance with the raw HTML option enabled
  * Write an URL parser/encoder that would support decoding punycode and counting matching `[(` and `)]` in URLs

## References / Thanks

[golang-commonmark/markdown](https://github.com/golang-commonmark/markdown) is a partial port of an awesome JavaScript markdown parser [markdown-it](https://github.com/markdown-it/markdown-it), written by Alex Kocharin and Vitaly Puzrin.

## See also

https://github.com/jgm/CommonMark — the reference CommonMark implementations in C and JavaScript,
  also contains the latest spec and an online demo.

http://talk.commonmark.org — the CommonMark forum, a good place to join together the efforts of the developers.
