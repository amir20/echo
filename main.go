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

type LogEntry struct {
	timestamp time.Time
	message   string
}

// parseTimestamp attempts to parse RFC3339 timestamp from the beginning of a line
func parseTimestamp(line string) (time.Time, string, bool) {
	// Try to find a timestamp at the beginning of the line
	// Support common formats like RFC3339, RFC3339Nano
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return time.Time{}, line, false
	}

	// Try parsing first field as timestamp
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05.999999999 -0700 MST",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, parts[0]); err == nil {
			// Remove timestamp from message
			message := strings.TrimSpace(strings.TrimPrefix(line, parts[0]))
			return t, message, true
		}
	}

	return time.Time{}, line, false
}

// replay reads timestamped logs from stdin and replays them with original timing
func replay(speedFactor float64) {
	scanner := bufio.NewScanner(os.Stdin)
	var entries []LogEntry

	// Read all entries
	for scanner.Scan() {
		line := scanner.Text()
		timestamp, message, ok := parseTimestamp(line)
		if ok {
			entries = append(entries, LogEntry{timestamp: timestamp, message: message})
		} else {
			// If no timestamp found, just echo the line immediately
			entries = append(entries, LogEntry{timestamp: time.Now(), message: line})
		}
	}

	if len(entries) == 0 {
		return
	}

	// Print first entry immediately
	fmt.Fprintln(os.Stderr, entries[0].message)

	// Replay subsequent entries with original timing adjusted by speed factor
	for i := 1; i < len(entries); i++ {
		duration := entries[i].timestamp.Sub(entries[i-1].timestamp)
		if duration > 0 {
			adjustedDuration := time.Duration(float64(duration) / speedFactor)
			// Cap wait time at 10 seconds maximum
			adjustedDuration = min(adjustedDuration, 10*time.Second)
			time.Sleep(adjustedDuration)
		}
		fmt.Fprintf(os.Stderr, "(%d) %s\n", i, entries[i].message)
	}
}

const randomData = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst. Praesent sapien turpis, fermentum vel, eleifend faucibus, vehicula eu, lacus. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst. Praesent sapien turpis, fermentum vel, eleifend faucibus, vehicula eu, lacus. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst. Praesent sapien turpis, fermentum vel, eleifend faucibus, vehicula eu, lacus.`

func main() {
	random := flag.Bool("r", false, "generate random data")
	burst := flag.Int64("b", -1, "generate large burst of data")
	sleep := flag.Int64("s", 1000, "sleep time")
	shuffle := flag.Bool("x", false, "shuffle data")
	numbers := flag.Bool("n", false, "show number")
	all := flag.Bool("a", false, "print all data and pause")
	playback := flag.Float64("p", 0, "replay logs with original timing (speed factor: 1=normal, 10=10x faster, 0=disabled)")
	flag.Parse()

	// Handle replay mode
	if *playback > 0 {
		replay(*playback)
		return
	}

	var data []string
	if *random {
		data = append(data, strings.Split(randomData, ". ")...)
	} else if *numbers {
		data = make([]string, 10000)
		for i := range data {
			data[i] = fmt.Sprintf("line %d", i)
		}
	} else {
		data = readData()
	}

	if *shuffle {
		rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
	}

	if *burst > 0 {
		go func() {
			for {
				time.Sleep(time.Millisecond * time.Duration(*burst))
				for range 5_000 {
					fmt.Fprintln(os.Stderr, data[rand.Intn(len(data))])
				}
			}
		}()
	}

	if *all {
		for _, line := range data {
			time.Sleep(5 * time.Millisecond)
			fmt.Println(line)
		}
		time.Sleep(time.Hour)
	} else {

		for i := 0; ; i = (i + 1) % len(data) {
			fmt.Fprintln(os.Stderr, data[i])
			time.Sleep(time.Millisecond * time.Duration(*sleep))
		}
	}
}
