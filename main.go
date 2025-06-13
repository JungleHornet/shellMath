package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/junglehornet/shellMath/interpreter"
	"golang.org/x/term"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func quit(scanner *bufio.Scanner) {
	fmt.Print("Are you sure you want to quit? (y/N) > ")
	scanner.Scan()
	switch strings.TrimSpace(strings.ToLower(scanner.Text())) {
	case "y":
		os.Exit(0)
	case "n":
		fmt.Print(strings.Repeat("\b", 50))
	default:
		fmt.Print(strings.Repeat("\b", 50))
	}
}

func main() {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil || width <= 0 {
		fmt.Println(errors.New("error getting terminal size: "))
		fmt.Println(err)
		fmt.Println("Defaulting to terminal width of 40 characters\n")
		width = 40
	}

	title := "Welcome to shellMath! Enter \"q\" to quit."
	paddingLen := float64(width - len(title))
	var paddingLeft, paddingRight string
	if int(paddingLen)%2 != 0 {
		paddingLeft = strings.Repeat(" ", int(paddingLen/2+0.5))
		paddingRight = strings.Repeat(" ", int(paddingLen/2-0.5))
	} else {
		paddingLeft = strings.Repeat(" ", int(paddingLen/2))
		paddingRight = paddingLeft
	}
	fmt.Println(paddingLeft + title + paddingRight)
	fmt.Println(strings.Repeat("=", width))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "q" {
			quit(scanner)
		} else {
			res, err := interpreter.Evaluate(line)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			} else {
				var resStr string
				if x := int(math.Round(res)); float64(x) == res {
					resStr = strconv.Itoa(x)
				} else {
					resStr = strconv.FormatFloat(res, 'f', -1, 64)
				}
				output := "= " + resStr
				fmt.Println(strings.Repeat(" ", width-utf8.RuneCountInString(output)-1), output)
			}
		}
		fmt.Println(strings.Repeat("-", width))
	}
}
