// https://github.com/vim/vim/blob/master/src/xxd/xxd.c
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	flag.Parse()
	filepath := flag.Arg(0)

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(f)
	// reader := bufio.NewScanner(f)
	buf1 := make([]byte, 16)
	buf2 := make([]byte, 16)

	offset := 0
	for i := 0; i < 16; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			if i > 0 {
				printLine(offset, buf1[0:i], buf2[0:i])
			}
			break
		}

		buf1[i] = b

		if b < 32 || b > 126 {
			b = 46 // ascii period
		}

		buf2[i] = b

		if i == 15 {
			i = -1
			printLine(offset, buf1, buf2)
			offset++
		}
	}
}

func printLine(offset int, bufLeft, bufRight []byte) {
	// rightpad────────┐
	//  leftpad─┐      │
	log.Printf("%06x0: %-32x  %s \n", offset, bufLeft, string(bufRight))
	// https://www.analyticsvidhya.com/blog/2024/06/xxd-command-in-linux/
}
