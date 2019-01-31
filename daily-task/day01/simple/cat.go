package simple

import (
	"os"
	"bufio"
	"fmt"
	"io"
)

func catFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

func main() {

	filename := ""
	arguments := os.Args
	if len(arguments) == 1 {
		io.Copy(os.Stdout, os.Stdin)
		os.Exit(0)
	}

	filename = arguments[1]
	err := catFile(filename)
	if err != nil {
		fmt.Println(err)
	}

}
