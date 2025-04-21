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
	postscript bool
	buffSize   int
	groupSize  int
	seek       int
	len        int
	hexPad     int
	in         io.ReadSeekCloser
	out        io.WriteCloser
}

const (
	DEFAULT_G   = 2
	DEFAULT_C   = 16
	DEFAULT_P_C = 30
	DEFAULT_S   = 0
	DEFAULT_L   = 0
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
	totalBytesRead := 0
	buf := make([]byte, cmd.buffSize)

	if cmd.seek > 0 {
		cmd.in.Seek(int64(cmd.seek), io.SeekStart)
		offset += cmd.seek
	}

	for {
		if cmd.len > 0 && totalBytesRead >= cmd.len {
			break
		}

		n, err := cmd.in.Read(buf)
		if err != nil {
			break
		}

		if n < cmd.buffSize {
			buf = buf[0:n] // truncate dangling data from previous iteration
		}

		if cmd.len > 0 && totalBytesRead+n > cmd.len {
			n = n - (totalBytesRead + n - cmd.len)
			buf = buf[0:n]
		}

		printLine(cmd, offset, buf)
		offset += n
		totalBytesRead += n
	}
}

func parseCmd() *command {
	cmd := &command{}

	flag.BoolFunc("p", "postscript", func(str string) error {
		cmd.postscript = true
		cmd.buffSize = DEFAULT_P_C
		return nil
	})
	flag.IntVar(&(cmd.buffSize), "c", DEFAULT_C, "cols: Format number bytes per output line")
	flag.IntVar(&(cmd.groupSize), "g", DEFAULT_G, "groupsize: Separate the output of number bytes per group in the hex dump")
	flag.IntVar(&(cmd.seek), "s", DEFAULT_S, "seek: Start at offset bytes from the beginning of the input file")
	flag.IntVar(&(cmd.len), "l", DEFAULT_L, "length: Stop after length bytes of the input file")
	flag.Parse()

	// -c
	if cmd.buffSize > MAX_C {
		dieAndDump(fmt.Errorf("invalid number of columns (max. %d)", MAX_C))
	} else if cmd.buffSize <= 0 {
		cmd.buffSize = DEFAULT_C
	}

	// -g
	if cmd.groupSize < 0 {
		cmd.groupSize = DEFAULT_G
	} else if cmd.groupSize == 0 {
		cmd.groupSize = cmd.buffSize * 2
	}

	// -s
	if cmd.seek < 0 {
		cmd.seek = DEFAULT_S
	}

	// -l
	if cmd.len < 0 {
		cmd.len = DEFAULT_L
	}

	cmd.hexPad = cmd.buffSize*2 + cmd.buffSize/cmd.groupSize

	parseInput(cmd)
	parseOutput(cmd)

	return cmd
}

func parseInput(cmd *command) {
	inputPath := flag.Arg(0)

	if inputPath != "" {
		f, err := os.Open(inputPath)
		dieAndDump(err)

		cmd.in = f
	} else {
		cmd.in = os.Stdin
	}
}

func parseOutput(cmd *command) {
	outputPath := flag.Arg(1)

	if outputPath != "" {
		outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
		dieAndDump(err)

		cmd.out = outFile
	} else {
		cmd.out = os.Stdout
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

func toChunkedHexString(bytes []byte, chunkSize int) string {
	var str string

	for b := range slices.Chunk(bytes, chunkSize) {
		str += fmt.Sprintf("%x ", b)
	}

	return str
}

func dieAndDump(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "xxd: %s\n", err)
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
		offset,                                               // leftpadded
		cmd.hexPad, toChunkedHexString(bytes, cmd.groupSize), // rightpaded
		toStrippedString(bytes), // replace non-ascii with '.'
	)
}
