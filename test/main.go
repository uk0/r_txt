package main

import (
	"bytes"
	"fmt"
	"strings"
)


// Wraps text at the specified column lineWidth on word breaks

func word_wrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	fmt.Println(len(words))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	fmt.Println(strings.Split(wrapped,"\n"))

	return wrapped

}

func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i + 1) % n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

func main() {

	longTextStr :=
		`体柔弱的宫装美妇人，从外面推门走进来，看着坐在床榻上的张若尘，带着关切的眼神，`

	fmt.Printf("Original: [%v] \n", longTextStr)
	fmt.Println("--------------------------------")

	wrapped := word_wrap(longTextStr, 4)

	fmt.Println(SplitSubN(longTextStr, 4)[0])
	// Some minimal html fixups
	// Note: this can introduce newlines inside class attributes, but that's perfectly
	// valid html (nb: http://stackoverflow.com/a/14928606)

	fmt.Println(wrapped)
}
