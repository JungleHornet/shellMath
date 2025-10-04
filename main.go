package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/junglehornet/shellMath/interpreter"
	"github.com/lunixbochs/vtclean"
	"golang.org/x/term"
)

func quit(scanner *bufio.Scanner) {
	fmt.Print("Are you sure you want to quit? (y/N) > ")
	scanner.Scan()
	switch strings.TrimSpace(strings.ToLower(scanner.Text())) {
	case "y":
		os.Exit(0)
	default:
		return
	}
}

func main() {
	errText := color.New(color.FgRed).Add(color.Bold).SprintFunc()

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		fmt.Println(errors.New(errText("error getting terminal size: ")))
		fmt.Println(errText(err))
		fmt.Println(errText("Defaulting to terminal width of 40 characters.\n"))
		width = 40
	}

	title := "Welcome to shellMath! Enter \"q\" to quit."
	paddingLen := float64(width - len(title)%width)
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
	var ans float64
	noAns := true
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		if line == "q" {
			quit(scanner)
		} else {
			firstType, _, _ := interpreter.DecodeTokenTypeInString(line)
			if firstType != interpreter.Value && firstType != interpreter.OpenParen {
				line = "ans" + line
			}

			var output string

			ansErr := false
			if strings.Contains(line, "ans") {
				if !noAns {
					line = strings.Replace(line, "ans", strconv.FormatFloat(ans, 'f', -1, 64), -1)
				} else {
					output = errText("Error: cannot use ans as there is no valid answer to use")
					ansErr = true
				}
			}

			res, err := interpreter.Evaluate(line)
			if err != nil || ansErr {
				if !ansErr {
					output = errText("Error: ", err)
				}
				noAns = true
			} else {
				var resStr string

				noAns = false
				ans = res
				if x := int(math.Round(res)); float64(x) == res {
					resStr = strconv.Itoa(x)
				} else {
					resStr = strconv.FormatFloat(res, 'f', -1, 64)
				}
				output = "= " + resStr
			}
			bufferWidth := width - utf8.RuneCountInString(vtclean.Clean(output, false)) - 1
			if bufferWidth < 0 {
				bufferWidth = width - utf8.RuneCountInString(vtclean.Clean(output, false))%width - 1
			}
			fmt.Println(strings.Repeat(" ", bufferWidth), output)
		}
		fmt.Println(strings.Repeat("-", width))
	}
}
