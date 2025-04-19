// https://github.com/vim/vim/blob/master/src/xxd/xxd.c
// https://www.analyticsvidhya.com/blog/2024/06/xxd-command-in-linux/
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	postscript bool
	buffSize   int
	groupSize  int
	hexPad     int
)

func main() {
	flag.BoolVar(&postscript, "p", false, "postscript")
	flag.IntVar(&buffSize, "c", 16, "column size")
	flag.IntVar(&groupSize, "g", 4, "group size")
	flag.Parse()

	hexPad = buffSize*2 + buffSize*2/groupSize

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

func toChunkedString(str string) string {
	var newStrBuf []rune

	for i, r := range str {
		if (i+1)%groupSize == 0 {
			newStrBuf = append(newStrBuf, r, ' ')
		} else {
			newStrBuf = append(newStrBuf, r)
		}
	}

	return string(newStrBuf)
}

func printLine(offset int, hexBuffer []byte) {
	// FIX column width!!!
	// if postscript {
	// 	fmt.Printf("%-32x", bufLeft)
	// 	return
	// }

	hexStr := fmt.Sprintf("%x", hexBuffer)
	chunkedStr := toChunkedString(hexStr)

	fmt.Printf(
		"%s: %s %s\n",
		fmt.Sprintf("%08x", offset),             // leftpad
		fmt.Sprintf("%-*s", hexPad, chunkedStr), // dynamic rightpad
		toStrippedString(hexBuffer),             // replace non-ascii with '.'
	)
}
