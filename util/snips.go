package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Ternary[T any](condition bool, If, Else T) T {
	if condition {
		return If
	}
	return Else
}

func TernaryPointer[T any](condition bool, If, Else func() T) T {
	if condition {
		return If()
	}
	return Else()
}

func PrintType[T any](x T) string {
	return fmt.Sprintf("%T", x)
}

func NewSlice(start, end, step int) []int {
	if step <= 0 || end < start {
		return []int{}
	}
	s := make([]int, 0, 1+(end-start)/step)
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

func AskForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Is this correct fuzz pipeline setup (Y/n): ")
		text, _ := reader.ReadString('\n')

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if (strings.Compare("y", text) == 0) || (strings.Compare("Y", text) == 0) || (strings.Compare("", text) == 0) {
			fmt.Println("✅ Fuzz pipeline is confirmed, starting soon ...")
			return true
		} else {
			fmt.Println("🚫 Fuzz pipeline declined. Exiting ...")
			return false
		}
	}
}
