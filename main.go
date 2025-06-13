package main

import (
	"fmt"
	"github.com/junglehornet/shellMath/interpreter"
	"os"
	"strings"
)

func main() {
	fmt.Println(interpreter.Evaluate(strings.Join(os.Args[1:], "")))
}
