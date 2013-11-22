// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

// Implementation of HyperLogLog from Flajolet et al, 2007. This is
// the 32-bit version from the paper only. A similar analysis would be
// required to get the constants right for 64-bit.
//
// This does not supply the 32-bit hash function. Check out hash/fnv
// from the go standard library.
//
package gohll

import (
	"fmt"
	"math"
)

const (
	two32 = 1 << 32
)

type HLL struct {
	alpha float64
	b     uint
	m     uint
	mask  uint32
	eM    []uint8
}

func alpha(b uint) float64 {
	switch {
	case b == 4:
		return 0.673
	case b == 5:
		return 0.697
	case b == 6:
		return 0.709
	case b > 6 && b < 17:
		return 0.7213 / (1 + 1.079/float64(uint(1)<<b))
	default:
		panic(fmt.Sprint("Value for b must be in [4,16]:", b))
	}
}

// Create a new HLL counter. 2**b is the number of registers. Good
// values include 4 to 16. Assumes 32-bit hashes which are good for
// counting sets with up to 1e9 unique items.
func NewHLL(b uint) *HLL {
	m := uint(1) << b
	return &HLL{alpha: alpha(b), b: b, m: m, eM: make([]uint8, m), mask: uint32(0xffffffff) >> b}
}

// Count leading zeros. From http://en.wikipedia.org/wiki/Find_first_set#Algorithms
func clz32(x uint32) uint {
	if x == 0 {
		return 32
	}
	n := uint(0)
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

// Add the hash value of an item to be counted to the counter.
func (h *HLL) Add(hashValue uint32) {
	j := hashValue >> (32 - h.b)
	emj := h.eM[j]
	w := uint8(clz32(hashValue&h.mask) - h.b + 1)
	if emj < w {
		h.eM[j] = w
	}
}

// Calculate the actual count.
func (h *HLL) Count() uint32 {
	acc := 0.0
	for _, eMj := range h.eM {
		acc += math.Pow(2.0, -float64(eMj))
	}

	fm := float64(h.m)
	E := h.alpha * fm * fm / acc

	if E <= 2.5*fm {
		count := 0
		for _, x := range h.eM {
			if x == 0 {
				count++
			}
		}
		if count == 0 {
			return uint32(E)
		}
		return uint32(fm * math.Log(fm/float64(count)))
	} else if E < 1/30.0*float64(math.MaxUint32) {
		return uint32(E)
	} else {
		return uint32(-two32 * math.Log(1.0-E/two32))
	}
}

// Work out the union between this HLL and another. This one is
// altered. Their b values must be the same.
func (h *HLL) Union(other *HLL) {
	if h.b != other.b {
		panic(fmt.Sprintf("HLLs must use the same number of bits for buckets: this=%d, other=%d", h.b, other.b))
	}
	for i, c := range other.eM {
		if h.eM[i] < c {
			h.eM[i] = c
		}
	}
}
