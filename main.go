package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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

func main() {
	data := readData()
	for {
		fmt.Println(data[rand.Intn(len(data))])
		time.Sleep(1 * time.Second)
	}
}
