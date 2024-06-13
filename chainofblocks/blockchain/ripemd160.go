// ripemd160.go
package main

import (
	
	"errors"
	"math/bits"
)

// i had a problem importing ripemd so i created one thanks to AI

const (
	h0 = 0x67452301
	h1 = 0xEFCDAB89
	h2 = 0x98BADCFE
	h3 = 0x10325476
	h4 = 0xC3D2E1F0
)

var r = [80]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	7, 4, 13, 1, 10, 6, 15, 3, 12, 0, 9, 5, 2, 14, 11, 8,
	3, 10, 14, 4, 9, 15, 8, 1, 2, 7, 0, 6, 13, 11, 5, 12,
	1, 9, 11, 10, 0, 8, 12, 4, 13, 3, 7, 15, 14, 5, 6, 2,
	4, 0, 5, 9, 7, 12, 2, 10, 14, 1, 3, 8, 11, 6, 15, 13,
}

var rp = [80]int{
	5, 14, 7, 0, 9, 2, 11, 4, 13, 6, 15, 8, 1, 10, 3, 12,
	6, 11, 3, 7, 0, 13, 5, 10, 14, 15, 8, 12, 4, 9, 1, 2,
	15, 5, 1, 3, 7, 14, 6, 9, 11, 8, 12, 2, 10, 0, 4, 13,
	8, 6, 4, 1, 3, 11, 15, 0, 5, 12, 2, 13, 9, 7, 10, 14,
	12, 15, 10, 4, 1, 5, 8, 7, 6, 2, 13, 14, 0, 3, 9, 11,
}

var s = [80]int{
	11, 14, 15, 12, 5, 8, 7, 9, 11, 13, 14, 15, 6, 7, 9, 8,
	7, 6, 8, 13, 11, 9, 7, 15, 7, 12, 15, 9, 11, 7, 13, 12,
	11, 13, 6, 7, 14, 9, 13, 15, 14, 8, 13, 6, 5, 12, 7, 5,
	11, 12, 14, 15, 14, 15, 9, 8, 9, 14, 5, 6, 8, 6, 5, 12,
	9, 15, 5, 11, 6, 8, 13, 12, 5, 12, 13, 14, 11, 8, 5, 6,
}

var sp = [80]int{
	8, 9, 9, 11, 13, 15, 15, 5, 7, 7, 8, 11, 14, 14, 12, 6,
	9, 13, 15, 7, 12, 8, 9, 11, 7, 7, 12, 7, 6, 15, 13, 11,
	9, 7, 15, 11, 8, 6, 6, 14, 12, 13, 5, 14, 13, 13, 7, 5,
	15, 5, 8, 11, 14, 14, 6, 14, 6, 9, 12, 9, 12, 5, 15, 8,
	8, 5, 12, 9, 12, 5, 14, 6, 8, 13, 6, 5, 15, 13, 11, 11,
}

func ripemd160Block(digest []uint32, block []byte) {
	var x [16]uint32
	for i := range x {
		x[i] = uint32(block[i*4]) | uint32(block[i*4+1])<<8 | uint32(block[i*4+2])<<16 | uint32(block[i*4+3])<<24
	}

	a, b, c, d, e := digest[0], digest[1], digest[2], digest[3], digest[4]
	aa, bb, cc, dd, ee := a, b, c, d, e

	for i := 0; i < 80; i++ {
		j := i
		t := a + f(j, b, c, d) + x[r[j]] + k(j)
		a, b, c, d, e = e, bits.RotateLeft32(t, s[j]), b, bits.RotateLeft32(c, 10), d

		t = aa + f(79-j, bb, cc, dd) + x[rp[j]] + kp(j)
		aa, bb, cc, dd, ee = ee, bits.RotateLeft32(t, sp[j]), bb, bits.RotateLeft32(cc, 10), dd
	}

	t := digest[1] + c + dd
	digest[1] = digest[2] + d + ee
	digest[2] = digest[3] + e + aa
	digest[3] = digest[4] + a + bb
	digest[4] = digest[0] + b + cc
	digest[0] = t
}

func f(j int, x, y, z uint32) uint32 {
	switch {
	case j < 16:
		return x ^ y ^ z
	case j < 32:
		return (x & y) | (^x & z)
	case j < 48:
		return (x | ^y) ^ z
	case j < 64:
		return (x & z) | (y & ^z)
	default:
		return x ^ (y | ^z)
	}
}

func k(j int) uint32 {
	switch {
	case j < 16:
		return 0x00000000
	case j < 32:
		return 0x5A827999
	case j < 48:
		return 0x6ED9EBA1
	case j < 64:
		return 0x8F1BBCDC
	default:
		return 0xA953FD4E
	}
}

func kp(j int) uint32 {
	switch {
	case j < 16:
		return 0x50A28BE6
	case j < 32:
		return 0x5C4DD124
	case j < 48:
		return 0x6D703EF3
	case j < 64:
		return 0x7A6D76E9
	default:
		return 0x00000000
	}
}


type ripemd160 struct {
	h   [5]uint32
	x   [64]byte
	nx  int
	len uint64
}

func New() *ripemd160 {
	d := new(ripemd160)
	d.Reset()
	return d
}

func (d *ripemd160) Reset() {
	d.h[0] = h0
	d.h[1] = h1
	d.h[2] = h2
	d.h[3] = h3
	d.h[4] = h4
	d.nx = 0
	d.len = 0
}

func (d *ripemd160) Size() int { return 20 }

func (d *ripemd160) BlockSize() int { return 64 }

func (d *ripemd160) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == 64 {
			ripemd160Block(d.h[:], d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= 64 {
		n := len(p) &^ (64 - 1)
		for i := 0; i < n; i += 64 {
			ripemd160Block(d.h[:], p[i:])
		}
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d *ripemd160) Sum(in []byte) []byte {
	len := d.len
	var tmp [64]byte
	copy(tmp[:], d.x[:d.nx])
	tmp[d.nx] = 0x80
	if d.nx < 56 {
		for i := d.nx + 1; i < 56; i++ {
			tmp[i] = 0
		}
	} else {
		for i := d.nx + 1; i < 64; i++ {
			tmp[i] = 0
		}
		ripemd160Block(d.h[:], tmp[:])
		for i := 0; i < 56; i++ {
			tmp[i] = 0
		}
	}
	len <<= 3
	for i := 0; i < 8; i++ {
		tmp[56+i] = byte(len >> (8 * i))
	}
	ripemd160Block(d.h[:], tmp[:])

	var digest [20]byte
	for i, s := range d.h {
		digest[i*4] = byte(s)
		digest[i*4+1] = byte(s >> 8)
		digest[i*4+2] = byte(s >> 16)
		digest[i*4+3] = byte(s >> 24)
	}
	return append(in, digest[:]...)
}

func ripemd160Hash(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	h := New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
