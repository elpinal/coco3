package main

import (
	"bufio"
	"fmt"
	"os"
)

var prompt = "Σ> "

func main() {
	var err error
	for err == nil {
		err = loop()
	}
	fmt.Fprintln(os.Stderr, err)
}

func loop() error {
	fmt.Print(prompt)
	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Print(line)
	return nil
}
