package main

import "fmt"

func main() {
	s := "有些人喜欢在中间加#标签 然后"
	fmt.Printf("len: %d\n", len(s))
	for i, v := range s {
		fmt.Printf("%d, %d, %s, #: %v, Nul: %v\n", i, v, s[i:i+1], v == '#', v == ' ')
	}

	for i := 0; i < len(s); i++ {
		fmt.Printf("%d, %d, %s, #: %v, Nul: %v\n", i, s[i], s[i:i+1], s[i] == '#', s[i] == ' ')
	}
}
