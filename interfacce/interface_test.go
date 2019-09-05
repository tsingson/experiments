package interfacce

import (
	"io"
	"testing"
)

func BenchmarkInterfaceCallSimple(b *testing.B) {
	// b.ReportAllocs()
	var z zero
	var g getter
	g = z
	b.ResetTimer()
	b.Run("via interfacce", func(b *testing.B) {
		b.ReportAllocs()
		total := 0
		for i := 0; i < b.N; i++ {
			total += g.get()
		}

		if total > 0 {
			b.Logf("total is %d", total)
		}
	})

	b.Run("direct", func(b *testing.B) {
		b.ReportAllocs()
		total := 0
		for i := 0; i < b.N; i++ {
			total += z.get()
		}

		if total > 0 {
			b.Logf("total is %d", total)
		}
	})
}

func BenchmarkInterfaceAlloc(b *testing.B) {
	var z zeroReader
	var r io.Reader
	r = z
	b.ResetTimer()
	b.Run("via interfacce", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			var buf [7]byte
			r.Read(buf[:])
		}
	})

	b.Run("direct", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			var buf [7]byte
			z.Read(buf[:])
		}
	})
}

func BenchmarkSliceConversion(b *testing.B) {
	numbers := make([]int, 100)
	for i := range numbers {
		numbers[i] = i
	}
	b.ResetTimer()
	b.Run("bad", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			numbersToStringsBad(numbers)
		}
	})

	b.Run("better", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			numbersToStringsBetter(numbers)
		}
	})
}
