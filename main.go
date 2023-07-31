package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// read data from stdin
func readData() []string {
	var data []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return data
}

const randomData = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Sed egestas egestas fringilla phasellus faucibus. Egestas sed sed risus pretium quam vulputate. Amet nisl suscipit adipiscing bibendum est ultricies integer quis auctor. Morbi non arcu risus quis varius quam quisque id diam.`

func main() {
	random := flag.Bool("r", false, "generate random data")
	flag.Parse()

	var data []string
	if *random {
		data = append(data, strings.Split(randomData, ". ")...)
	} else {
		data = readData()
	}

	for {
		fmt.Fprintln(os.Stderr, data[rand.Intn(len(data))])
		time.Sleep(1 * time.Second)
	}
}
