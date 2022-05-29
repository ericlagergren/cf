package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	cf := cf{
		fn: "exp(z)",
		bN: []string{
			"1", "2-z", "6", "10", "14",
		},
		aN: []string{
			"2z", "z^2",
		},
	}
	fmt.Println(cf)
}

type cf struct {
	fn string
	bN []string
	aN []string
}

func (c cf) String() string {
	fn := utf8.RuneCountInString(c.fn)
	if fn != 0 {
		fn += len(" = ")
	}

	// Calculate the maximum line length.
	max := fn
	for _, v := range c.bN {
		max += len(v) + 1
	}
	max += utf8.RuneCountInString("─────")

	var b strings.Builder
	b.WriteString("// ")

	bpad := 0
	dash := max
	for i, v := range c.bN {
		d := len(v)
		dash -= d + 1

		var a string
		if i < len(c.aN) {
			a = c.aN[i]
		} else {
			a = c.aN[len(c.aN)-1]
		}
		apad := (dash - 1) / 2
		apad -= len(a) / 2
		apad += bpad + d + len(" + ")
		if i == 0 && fn != 0 {
			apad += fn
		}
		for i := 0; i < apad; i++ {
			b.WriteByte(' ')
		}
		b.WriteString(a)
		b.WriteByte('\n')
		b.WriteString("// ")

		if i == 0 && fn != 0 {
			b.WriteString(c.fn)
			b.WriteString(" = ")
		}
		for i := 0; i < bpad; i++ {
			b.WriteByte(' ')
		}
		b.WriteString(v)
		b.WriteString(" + ")
		bpad += d + 1
		if i == 0 && fn != 0 {
			bpad += fn
		}

		for i := 0; i < dash; i++ {
			b.WriteRune('─')
		}
		b.WriteByte('\n')
		b.WriteString("// ")
	}
	for i := 0; i < max+fn-1; i++ {
		b.WriteByte(' ')
	}
	b.WriteString("...")
	return b.String()
}
