package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ericlagergren/decimal"
)

// https://math.stackexchange.com/a/716976/153292
//
// pi = 3.14159
// a0 = floor(pi) = 3
// b0 = pi-a0 = 0.14159
//
// pi = a0 + b0
//    = a0 + 1/x1 where x1 = 1/b0
//
//             1
//    = 3 + --------,
//             1
//
//               z    z^2   z^2   z^2
//    tan(z) = ----- ----- ----- -----
//               1  -  3  -  5  -  7  - ···
//
//                         2z              a1
// exp(z) = 1 + ────────────────────────── b0
//                           z^2           a2
//            2-z + ────────────────────── b1
//                            z^2          a3
//                6 + ──────────────────── b2
//                              z^2        a4
//                  10 + ───────────────── b3
//                               z^2       a5
//                     14 + ────────────── b3
//                                     ...

func main() {
	main3()

	println()

	printGen(&e{}, 11)
	printGen(&rat{3, 8}, 11)
	printGen(&rat{13, 11}, 11)
}

type gen interface {
	next() bool
	get() uint64
}

func printGen(g gen, n int) {
	fmt.Printf("[%d", g.get())
	fmt.Printf("; %d", g.get())
	for i := 0; i < 11 && g.next(); i++ {
		fmt.Printf(", %d", g.get())
	}
	fmt.Println("]")
}

type rat struct {
	n, d uint64
}

func (r *rat) next() bool {
	return r.d != 0
}

func (r *rat) get() uint64 {
	n, d := r.n, r.d
	p := n / d
	r.n = d
	r.d = n - p*d
	return p
}

type e struct {
	n uint64
}

func (e) next() bool {
	return true
}

func (e *e) get() uint64 {
	n := e.n
	e.n++
	switch {
	case n == 0:
		return 2
	case n%3 == 2:
		return 2*(n/3) + 2
	default:
		return 1
	}
}

func main3() {
	ctx := decimal.Context{
		Precision: 50,
	}
	x := ctx.Pi(new(decimal.Big))
	fmt.Println(cfstr(cf(ctx, x, 36)))
	r, _ := new(big.Rat).SetString(x.String())
	fmt.Println(cfstr(cf2(r, 36)))
	fmt.Println("[3; 7, 15, 1, 292, 1, 1, 1, 2, 1, 3, 1, 14, 2, 1, 1, 2, 2, 2, 2, 1, 84, 2, 1, 1, 15, 3, 13, 1, 4, 2, 6, 6, 99, 1, 2]")

	fmt.Println()

	x.SetString("1.694889244410333714141783611437197494892623622551650491315726964531624162040e28")
	fmt.Println(cfstr(cf(ctx, x, 31)))
	r, _ = new(big.Rat).SetString(x.String())
	fmt.Println(cfstr(cf2(r, 31)))
	fmt.Println("[16948892444103337141417836114; 2, 1, 2, 4, 1, 3, 1, 3, 3, 1, 1, 5, 1, 2, 1, 115, 1, 1, 1, 1, 1, 5, 1, 1]")

	r, _ = new(big.Rat).SetString("1.694889244410333714141783611437197494892623622551650491315726964531624162040e28")
	fmt.Println(cfstr(cfinf(r)))
}

func cfstr(a []big.Int) string {
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < len(a); i++ {
		switch v := &a[i]; i {
		case 0:
			b.WriteString(v.String())
		case 1:
			b.WriteString("; ")
			b.WriteString(v.String())
		default:
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteByte(']')
	return b.String()
}

func cfinf(x *big.Rat) []big.Int {
	var a []big.Int
	var z big.Rat
	z.Set(x)
	n, d := z.Num(), z.Denom()
	for i := 0; i < 100; i++ {
		var t big.Int
		// a_i = \floor{x}
		// x = 1 / (x - \floor{x})
		t.QuoRem(n, d, n)
		a = append(a, t)
		if n.Sign() == 0 {
			break
		}
		n, d = d, n
	}
	return a

}

func cf2(x *big.Rat, k int) []big.Int {
	a := make([]big.Int, k)
	var z big.Rat
	z.Set(x)
	n, d := z.Num(), z.Denom()
	for i := 0; i < k; i++ {
		// a_i = \floor{x}
		// x = 1 / (x - \floor{x})
		a[i].QuoRem(n, d, n)
		n, d = d, n
	}
	return a

}

var one = decimal.New(1, 0)

func cf(ctx decimal.Context, x *decimal.Big, n int) []big.Int {
	//a := make([]big.Int, n)
	//z := new(decimal.Big).Copy(x)
	var z big.Rat
	x.Rat(&z)
	return cf2(&z, n)
	// var t decimal.Big
	// for i := 0; i < n; i++ {
	// 	ctx.Floor(&t, z)
	// 	t.Int(&a[i])
	// 	ctx.Sub(z, z, &t)
	// 	ctx.Quo(z, one, z)
	// }
	// return a
}

func main2() {
	a := []int{1, 0, 0, 2}
	b := lagrange(a, 10)
	fmt.Println(b)

	var x big.Int
	// x[0] + x[1]*_B^1 + ... x[n-1]*_B^(n-1)
	x.SetBits([]big.Word{1, 0, 0, 2})
	fmt.Println(x.String())
	// x[n-1]*_B^(n-1) + ... + x[1]*_B^1 + x[0]
	x.SetBits([]big.Word{2, 0, 0, 1})
	fmt.Println(x.String())
}

//     a : list - coefficients of the polynomial,
//         i.e. f(x) = a[0] + a[1]*x + ... + a[n]*x^n
//     N : number of quotients to output
//
// 1.2599210498948731647672 ...
func lagrange(a []int, p int) []int {
	ans := make([]int, p)
	shift := func() {
		for k := range a {
			for j := len(a) - 2; j > k-1; j-- {
				a[j] += a[j+1]
			}
		}
	}
	for i := 0; i < p; i++ {
		q := int(1)
		shift()

		for sum(a) < 0 {
			q++
			shift()
		}
		ans[i] = q

		revNeg(a)
	}
	return ans
}

func revNeg(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = -a[j], -a[i]
	}
}

func sum(a []int) int {
	var x int
	for _, v := range a {
		x += v
	}
	return x
}
