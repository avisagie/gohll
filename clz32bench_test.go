// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package gohll

import (
	"testing"
)

func clz32_a(x uint32) uint8 {
	// http://en.wikipedia.org/wiki/Find_first_set#Algorithms
	if x == 0 {
		return 32
	}
	n := uint8(0)
	if (x & 0xFFFF0000) == 0 {
		n = n + 16
		x = x << 16
	}
	if (x & 0xFF000000) == 0 {
		n = n + 8
		x = x << 8
	}
	if (x & 0xF0000000) == 0 {
		n = n + 4
		x = x << 4
	}
	if (x & 0xC0000000) == 0 {
		n = n + 2
		x = x << 2
	}
	if (x & 0x80000000) == 0 {
		n = n + 1
	}
	return n
}

func clz32_b(x uint32) uint8 {
	// http://aggregate.org/MAGIC/#Leading Zero Count
	x |= (x >> 1)
	x |= (x >> 2)
	x |= (x >> 4)
	x |= (x >> 8)
	x |= (x >> 16)
	// http://aggregate.org/MAGIC/#Population Count (Ones Count)
	x -= ((x >> 1) & 0x55555555)
	x = (((x >> 2) & 0x33333333) + (x & 0x33333333))
	x = (((x >> 4) + x) & 0x0f0f0f0f)
	x += (x >> 8)
	x += (x >> 16)
	//x &= 0x3f
	// rest of Leading Zero Count
	return 32 - uint8(x)
}

func clz32_c(x uint32) uint8 {
	// Hacker's Delight (1st ed), page 80, figure 5-10
	y := -(x >> 16)
	m := (y >> 16) & 16
	n := 16 - m
	x = x >> m
	y = x - 0x100
	m = (y >> 16) & 4
	n += m
	x <<= m
	y = x - 0x1000
	m = (y >> 16) & 8
	n += m
	x <<= m
	y = x - 0x4000
	m = (y >> 16) & 2
	n += m
	x <<= m
	y = x >> 14
	m = y & ^(y >> 1)
	return uint8(n + 2 - m)
}

func BenchmarkClzWikipedia(b *testing.B) {
	for i := uint32(b.N); i != 0; i-- {
		_ = clz32_a(i)
	}
}

func BenchmarkClzAggregate(b *testing.B) {
	for i := uint32(b.N); i != 0; i-- {
		_ = clz32_b(i)
	}
}

func BenchmarkClzHackersDelight(b *testing.B) {
	for i := uint32(b.N); i != 0; i-- {
		_ = clz32_c(i)
	}
}
