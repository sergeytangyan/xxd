// https://github.com/vim/vim/blob/master/src/xxd/xxd.c
// https://www.analyticsvidhya.com/blog/2024/06/xxd-command-in-linux/
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

var (
	postscript bool
	buffSize   int
)

func main() {
	flag.BoolVar(&postscript, "p", false, "postscript")
	flag.IntVar(&buffSize, "c", 16, "column size")
	flag.Parse()

	filepath := flag.Arg(0)
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := make([]byte, buffSize)

	offset := 0
	for {
		n, err := f.Read(buf)
		if err != nil {
			break
		}

		if n < buffSize {
			buf = buf[0:n]
		}

		printLine(offset, buf, strippedBuf(buf))
		offset += n
	}
}

func strippedBuf(buf []byte) []byte {
	newBuf := make([]byte, len(buf))

	for i, b := range buf {
		if b < 32 || b > 126 {
			b = 46 // ascii period
		}

		newBuf[i] = b
	}

	return newBuf
}

func printLine(offset int, bufLeft, bufRight []byte) {
	// FIX column width!!!
	// if postscript {
	// 	fmt.Printf("%-32x", bufLeft)
	// 	return
	// }

	// rightpad────────┐
	//  leftpad─┐      │
	fmt.Printf("%08x: %-32x  %s \n", offset, bufLeft, string(bufRight))
}
