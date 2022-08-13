package main

import "fmt"

type tree []node

type node struct {
	label    string
	children []int // indexes into tree
}

func vis(t tree) {
	if len(t) == 0 {
		fmt.Println("<empty>")
		return
	}

	var f func(int, string)

	f = func(n int, pre string) {
		ch := t[n].children
		if len(ch) == 0 {
			fmt.Println("╴", t[n].label)
			return
		}
		fmt.Println("┐", t[n].label)
		last := len(ch) - 1
		for _, ch := range ch[:last] {
			fmt.Print(pre, "├─")
			f(ch, pre+"│ ")
		}
		fmt.Print(pre, "└─")
		f(ch[last], pre+"  ")
	}

	f(0, "")
}
