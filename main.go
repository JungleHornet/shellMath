package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/junglehornet/shellMath/interpreter"
	"github.com/lunixbochs/vtclean"
	"github.com/mattn/go-tty"
	"golang.org/x/term"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func quit(t *tty.TTY) {
	prompt := "Are you sure you want to quit? (y/N) > "
	fmt.Print(prompt)
	inpt, _ := t.ReadRune()
	switch strings.TrimSpace(strings.ToLower(string(inpt))) {
	case "y":
		os.Exit(0)
	case "n":
		returningText := "n | Returning to program..."
		fmt.Print(returningText)
		time.Sleep(1 * time.Second)
		rep := len(prompt) + len(returningText)
		backspaces := strings.Repeat("\b", rep)
		fmt.Print(backspaces + strings.Repeat(string(rune(127)), rep) + backspaces)
		return
	default:
		invalidText := " | Invalid answer, continuing program."

		fmt.Print(string(inpt) + invalidText)
		time.Sleep(1 * time.Second)
		rep := len(prompt) + 1 + len(invalidText)
		backspaces := strings.Repeat("\b", rep)
		fmt.Print(backspaces + strings.Repeat(string(rune(127)), rep) + backspaces)
	}
}

func runLine(line string, noAns bool, errText func(a ...interface{}) string, ans float64, width int) (bool, float64) {
	first, _ := utf8.DecodeRuneInString(line)
	if interpreter.IsOperatorRune(first) {
		line = "ans" + line
	}

	var output string

	ansErr := false
	if strings.Contains(line, "ans") {
		if !noAns {
			line = strings.Replace(line, "ans", "("+strconv.FormatFloat(ans, 'f', -1, 64)+")", -1)
		} else {
			output = errText("Error: cannot use ans as there is no valid answer to use")
			ansErr = true
		}
	}

	line = strings.Replace(line, "//", "÷", -1)
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
	fmt.Println(strings.Repeat("-", width))

	return noAns, ans
}

func main() {
	errText := color.New(color.FgRed).Add(color.Bold).SprintFunc()

	width, _, err := term.GetSize(int(os.Stdin.Fd()))
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

	t, err := tty.Open()
	if err != nil {
		fmt.Println(errText(err))
	}
	defer t.Close()

	var ans float64
	noAns := true
	prevLines := []string{
		"",
	}
	var lineCursor int
	var line string
	prevLen := 0
	cursor := 0
	inSequence := false
	for {
		r, err := t.ReadRune()
		width, _, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil || width <= 0 {
			fmt.Println(errors.New(errText("error getting terminal size: ")))
			fmt.Println(errText(err))
			fmt.Println(errText("Defaulting to terminal width of 40 characters.\n"))
			width = 40
		}

		switch r {
		case 'q':
			quit(t)
		case 13:
			fmt.Println()
			cursor = 0
			prevLines = slices.Insert(prevLines, 0, "")
			noAns, ans = runLine(line, noAns, errText, ans, width)
			line = ""
			prevLen = 0
		case 27:
			inSequence = true
		case 91:
			inSequence = true
		case 127:
			cursor--
			if cursor < 0 {
				cursor = 0
			}
			line = line[:cursor] + line[cursor+1:]
			fmt.Print(strings.Repeat("\b", prevLen), line+" \b"+strings.Repeat("\b", len(line)-cursor))
			prevLen--
		default:
			if inSequence {
				inSequence = false
				chLine := false
				switch r {
				case 65:
					// up
					chLine = true
				case 66:
					// down
					chLine = true
				case 67:
					// right
					//fmt.Print(string(rune(27)) + "[1C")
					if cursor >= len(line) {
						fmt.Print("\b")
						cursor--
					}
				case 68:
					// left
					if cursor < 1 {
						cursor--
					} else if cursor == 1 {
						fmt.Print("\b")
						cursor -= 2
					} else if cursor > 1 {
						fmt.Print("\b\b")
						cursor -= 2
					}
				default:
					//fmt.Println(int(r))
				}
				if chLine {
					if lineCursor == 0 {
						if r == 66 {
							continue
						}
					}
					if r == 66 {
						lineCursor--
					} else if r == 65 && lineCursor < (len(prevLines)-1) {
						lineCursor++
					}
					line = prevLines[lineCursor]
				}
			} else {
				if interpreter.IsOperatorRune(r) && cursor == 0 && !noAns {
					line = "ans"
					fmt.Print("   ")
					cursor += 3
				}
				line = line[:cursor] + string(r) + line[cursor:]
			}
			lineCursor = 0
			prevLines[0] = line
			cursor++
			prefixLen := cursor - 1
			if prefixLen < 0 {
				prefixLen = 0
			}
			fmt.Print(strings.Repeat("\b", prefixLen) + line + strings.Repeat("\b", len(line)-cursor))
			prevLen = len(line)
		}
	}
}
