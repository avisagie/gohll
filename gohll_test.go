// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package gohll

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"testing"
)

func TestGohll1(t *testing.T) {
	hll := NewHLL(5)
	if len(hll.eM) != 32 {
		t.Fatal("1")
	}

	hll.Add(1)
	if hll.eM[0] != 27 {
		t.Fatal("2", hll.eM)
	}
	hll.Add(0)
	if hll.eM[0] != 28 {
		t.Fatal("3")
	}

	x := uint32(17) << 27
	hll.Add(1 | x)
	if hll.eM[17] != 27 {
		t.Fatal("4")
	}
	hll.Add(0 | x)
	if hll.eM[17] != 28 {
		t.Fatal("5")
	}
}

func TestGohll2(t *testing.T) {
	for N := 10; N <= 10000000; N *= 100 {
		for b := uint(4); b <= 16; b += 2 {
			if e := math.Abs(exp(t, b, N)); (b < 10 && e > 30.0) || (b >= 10 && e > 7.0) {
				t.Fatal("Error out of bounds:", e)
			}
		}
	}
}

func exp(t *testing.T, b uint, N int) (errorPerc float64) {
	hll := NewHLL(b)
	for ii := 0; ii < N; ii++ {
		hll.Add(rand.Uint32())
	}
	C := hll.Count()
	errorPerc = 100.0 * (float64(N) - float64(C)) / float64(N)
	t.Logf("b = %d, expect %d, got %d, error %0.02f%%", b, N, C, errorPerc)
	return
}

func TestClz(t *testing.T) {
	x := uint32(0xffffffff)
	for ii := uint(0); ii < 32; ii++ {
		if z := clz32(x); z != ii {
			t.Fatalf("expected %d, clz says %d. x = 0x%08x", ii, z, x)
		}
		x = x >> 1
	}
}

func add(hll *HLL) {
	f := fnv.New32()
	for i := 0; i < 1000; i++ {
		x := rand.Int31n(10000)
		f.Write([]byte(fmt.Sprintf("%d%d", x, x)))
		hll.Add(f.Sum32())
		f.Reset()
	}
}

func TestHLLUnion(t *testing.T) {
	hll := NewHLL(10)

	add(hll)
	c := hll.Count()
	t.Logf("Expect 1000, got %d", c)
	for i := 0; i < 10; i++ {
		prev := c
		add(hll)
		c := hll.Count()
		t.Logf("Got %d", c)

		// A slightly simplistic test: at least make sure it
		// keeps getting bigger.
		if c < prev {
			t.Fatalf("Count should not get smaller: %d < %d", c, prev)
		}
	}
}
