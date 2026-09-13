package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) < 1 {
		usage()
	}

	switch os.Args[1] {
	case "avalanche":
		avalanche([]byte(os.Args[2]))
	case "birthday":
		num, err := strconv.Atoi(os.Args[2])
		if err != nil {
			usage()
		}
		birthday(num)
	default:
		usage()
	case "verify":
		hashed := sha256.Sum256([]byte(os.Args[2]))
		if hex.EncodeToString(hashed[:]) == os.Args[3] {
			fmt.Println("OK")
		}
		fmt.Printf("MISMATCH\n  expected: %s\n  actual:   %s\n", os.Args[3], hashed)
		os.Exit(1)
	case "hash":
		hashed := sha256.Sum256([]byte(os.Args[2]))
		fmt.Printf("hashed data: %x \n", hex.EncodeToString(hashed[:]))
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  hashlab hash <string>
  hashlab verify <string> <expected-hex>`)
	os.Exit(2)
}

func hamming(a, b []byte) int {
	n := 0
	for i := range a {
		n += bits.OnesCount8(a[i] ^ b[i])
	}
	return n
}

func avalanche(input []byte) {
	base := sha256.Sum256(input)
	total, count, min, max := 0, 0, 256, 0

	for i := range input {
		for bit := 0; bit < 8; bit++ {
			m := make([]byte, len(input))
			copy(m, input)
			m[i] ^= 1 << bit // XOR, flipping a bit
			h := sha256.Sum256(m)
			d := hamming(base[:], h[:])
			total += d
			count++
			if d < min {
				min = d
			}
			if d > max {
				max = d
			}
		}
	}
	fmt.Printf("flipped %d input bit\n", count)
	fmt.Printf("output changed bit number: on average of  %.1f / 256(%.1f%%)\n",
		float64(total)/float64(count), 100*float64(total)/float64(count)/256)
	fmt.Printf("least %d,most %d\n", min, max)
}

func truncateBits(h []byte, nbits int) string {
	nb := (nbits + 7) / 8
	out := make([]byte, nb)
	copy(out, h[:nb])
	if rem := nbits % 8; rem != 0 {
		out[nb-1] &= byte(0xFF) << (8 - rem)
	}
	return string(out)
}

func birthday(nbits int) {
	seen := make(map[string]string)
	for i := 0; ; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		h := sha256.Sum256([]byte(msg))
		key := truncateBits(h[:], nbits)

		if prev, dup := seen[key]; dup {
			fmt.Printf("在 %d 次尝试后找到碰撞\n", i)
			fmt.Printf("  理论期望约 2^(%d/2) = %.0f 次\n", nbits, math.Pow(2, float64(nbits)/2))
			fmt.Printf("  %q\n  %q\n", prev, msg)
			ph := sha256.Sum256([]byte(prev))
			fmt.Printf("  截断到 %d bit 后都是 %s\n", nbits,
				hex.EncodeToString([]byte(key)))
			fmt.Printf("  完整 hash 其实不同:\n    %s\n    %s\n",
				hex.EncodeToString(ph[:8]), hex.EncodeToString(h[:8]))
			return
		}
		seen[key] = msg
	}
}
