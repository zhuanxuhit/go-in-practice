package main

import (
	"flag"
	"io"
	"bufio"
	"fmt"
	"os"
)

var (
	minusB = flag.Bool("b", false, "Number the non-blank output lines, starting at 1.")
	//minusE = flag.Bool("e",false,"Display non-printing characters (see the -v option), and display a dollar sign (`$') at the end of each line.")
	minusN = flag.Bool("n", false, "Number the output lines, starting at 1.")
	minusS = flag.Bool("s", false, "Squeeze multiple adjacent empty lines, causing the output to be single spaced.")
	//minusT = flag.Bool("t",false,"Display non-printing characters (see the -v option), and display tab characters as `^I'.")
	//minusU = flag.Bool("u",false,"Disable output buffering.")
	//minusV = flag.Bool("v",false,"Display non-printing characters so they are visible.  Control characters print as `^X' for control-X; the delete character (octal 0177) prints as `^?'.  Non-ASCII characters (with the high bit set) are printed as `M-' (for meta) followed by the character for the low 7 bits.")
)

func copy(dst io.Writer, src io.Reader) (written int64, err error) {
	var lastLine = "\n"
	nr := 0
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		// 优先判断空行
		if line == "" && *minusS && lastLine == "" {
			continue
		}
		// 再判断是否不包含空格
		if line == "" && *minusB {
			fmt.Fprintln(dst, line)
		} else if *minusB || *minusN {
			// 只要 -n -b 都要计数
			nr++
			fmt.Fprintf(dst, "%6d\t%s\n", nr, line)
		} else {
			fmt.Fprintln(dst, line)
		}
		lastLine = line
	}
	err = scanner.Err()
	return
}

func main() {
	flag.Parse()
	// cat 本质上是将 src copy to dst
	//io.Copy()
	myCopy := io.Copy
	if *minusB || *minusS || *minusN {
		myCopy = copy
	}
	if len(flag.Args()) == 0 {
		myCopy(os.Stdout, os.Stdin)
	}
	for _, fname := range flag.Args() {
		if fname == "-" {
			myCopy(os.Stdout, os.Stdin)
		} else {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			myCopy(os.Stdout, f)
			f.Close()
		}
	}
}
