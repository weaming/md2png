package main

import (
	"fmt"
)

func renderCSS(path string) string {
	block := "<style>\n%v\n</style>\n\n"
	return fmt.Sprintf(block, readFile(path))
}
