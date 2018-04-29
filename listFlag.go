package main

import (
	"strings"
	"fmt"
)

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return fmt.Sprintf("[%v]", strings.Join(*i, ", "))
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
