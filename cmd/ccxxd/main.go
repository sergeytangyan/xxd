// https://codingchallenges.fyi/challenges/challenge-xxd
// https://github.com/vim/vim/blob/master/src/xxd/xxd.c
// https://www.analyticsvidhya.com/blog/2024/06/xxd-command-in-linux/
package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

var (
	// postscript    bool
	byteBuffSize  int
	byteGroupSize int
	hexPad        int
)

const (
	DEFAULT_G_SIZE = 2
	DEFAULT_C_SIZE = 16
)

func main() {
	// postscript = *flag.Bool("p", false, "postscript")
	flag.IntVar(&byteBuffSize, "c", DEFAULT_C_SIZE, "column size")
	flag.IntVar(&byteGroupSize, "g", DEFAULT_G_SIZE, "group size")
	flag.Parse()

	if byteBuffSize <= 0 {
		byteBuffSize = DEFAULT_C_SIZE
	}
	if byteGroupSize < 0 {
		byteGroupSize = DEFAULT_G_SIZE
	}

	hexPad = byteBuffSize*2 + byteBuffSize/byteGroupSize

	filepath := flag.Arg(0)
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	offset := 0
	buf := make([]byte, byteBuffSize)

	for {
		n, err := f.Read(buf)
		if err != nil {
			break
		}

		if n < byteBuffSize {
			buf = buf[0:n]
		}

		printLine(offset, buf)
		offset += n
	}
}

func toStrippedString(buf []byte) string {
	newBuf := make([]byte, len(buf))

	for i, b := range buf {
		if b < 32 || b > 126 {
			b = 46 // ascii period
		}

		newBuf[i] = b
	}

	return string(newBuf)
}

func toChunkedHexString(bytes []byte) string {
	var str string

	for b := range slices.Chunk(bytes, byteGroupSize) {
		str += fmt.Sprintf("%x ", b)
	}

	return str
}

func printLine(offset int, bytes []byte) {
	// FIX column width!!!
	// if postscript {
	// 	fmt.Printf("%-32x", bufLeft)
	// 	return
	// }

	fmt.Printf(
		"%08x: %-*s %s\n",
		offset,                            // leftpadded
		hexPad, toChunkedHexString(bytes), // rightpaded
		toStrippedString(bytes), // replace non-ascii with '.'
	)
}
