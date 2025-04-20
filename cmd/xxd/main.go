// https://codingchallenges.fyi/challenges/challenge-xxd
// https://github.com/vim/vim/blob/master/src/xxd/xxd.c
// https://www.mankier.com/1/xxd
// https://www.analyticsvidhya.com/blog/2024/06/xxd-command-in-linux/
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
)

type command struct {
	postscript    bool
	byteBuffSize  int
	byteGroupSize int
	hexPad        int
	out           io.WriteCloser
	in            io.ReadCloser
}

const (
	DEFAULT_G   = 2
	DEFAULT_C   = 16
	DEFAULT_P_C = 30
	MAX_C       = 256
)

func main() {
	cmd := parseCmd()
	xxd(cmd)
}

func xxd(cmd *command) {
	defer cmd.in.Close()
	defer cmd.out.Close()

	offset := 0
	buf := make([]byte, cmd.byteBuffSize)

	for {
		n, err := cmd.in.Read(buf)
		if err != nil {
			break
		}

		if n < cmd.byteBuffSize {
			buf = buf[0:n]
		}

		printLine(cmd, offset, buf)
		offset += n
	}
}

func parseCmd() *command {
	cmd := &command{}

	flag.BoolFunc("p", "postscript", func(str string) error {
		cmd.postscript = true
		cmd.byteBuffSize = DEFAULT_P_C
		return nil
	})
	flag.IntVar(&(cmd.byteBuffSize), "c", DEFAULT_C, "column size")
	flag.IntVar(&(cmd.byteGroupSize), "g", DEFAULT_G, "group size")
	flag.Parse()

	if cmd.byteBuffSize > MAX_C {
		dieAndDump(fmt.Errorf("invalid number of columns (max. %d)", MAX_C))
	} else if cmd.byteBuffSize <= 0 {
		cmd.byteBuffSize = DEFAULT_C
	}

	if cmd.byteGroupSize < 0 {
		cmd.byteGroupSize = DEFAULT_G
	} else if cmd.byteGroupSize == 0 {
		cmd.byteGroupSize = cmd.byteBuffSize * 2
	}

	cmd.hexPad = cmd.byteBuffSize*2 + cmd.byteBuffSize/cmd.byteGroupSize

	inputPath := flag.Arg(0)
	f, err := os.Open(inputPath)
	dieAndDump(err)
	cmd.in = f

	outputPath := flag.Arg(1)
	if outputPath != "" {
		// TODO: FIX WRITING TO FILE
		outFile, err := os.OpenFile(outputPath, os.O_CREATE, 0o666)
		dieAndDump(err)
		cmd.out = outFile
	} else {
		cmd.out = os.Stdout
	}

	return cmd
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

func toChunkedHexString(bytes []byte, chunkSize int) string {
	var str string

	for b := range slices.Chunk(bytes, chunkSize) {
		str += fmt.Sprintf("%x ", b)
	}

	return str
}

func dieAndDump(err error) {
	if err != nil {
		fmt.Printf("xxd: %s\n", err)
		os.Exit(1)
	}
}

func printLine(cmd *command, offset int, bytes []byte) {
	if cmd.postscript {
		fmt.Fprintf(cmd.out, "%-*x\n", cmd.hexPad, bytes)
		return
	}

	fmt.Fprintf(
		cmd.out,
		"%08x: %-*s %s\n",
		offset,                                                   // leftpadded
		cmd.hexPad, toChunkedHexString(bytes, cmd.byteGroupSize), // rightpaded
		toStrippedString(bytes), // replace non-ascii with '.'
	)
}
