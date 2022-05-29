package main

import (
	"math/big"
	"testing"

	"github.com/ericlagergren/decimal"
)

var sink []big.Int

const str = "1.694889244410333714141783611437197494892623622551650491315726964531624162040e28"

func BenchmarkCFDec(b *testing.B) {
	x, _ := new(decimal.Big).SetString(str)
	ctx := decimal.Context{
		Precision: 5,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = cf(ctx, x, 31)
	}
}

func BenchmarkCFRat(b *testing.B) {
	x, _ := new(big.Rat).SetString(str)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = cf2(x, 31)
	}
}
