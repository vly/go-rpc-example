package main

import(
	"fmt"
	"unicode/utf8"
)

// testing out rune encoding, hoping for a stable 2 byte rep of utf8
func main() {
	r := 'a'
	buf := make([]byte, 2)

	n := utf8.EncodeRune(buf, r)

	fmt.Println(buf)
	fmt.Println(n)
}