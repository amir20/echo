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

// parseTimestamp attempts to parse a timestamp from the beginning of a line.
func parseTimestamp(line string) (time.Time, string, bool) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return time.Time{}, line, false
	}

	// Each layout is paired with the number of whitespace-separated tokens it
	// spans. RFC3339 timestamps are a single token; "date time" formats (such
	// as Python's logging output) span two.
	candidates := []struct {
		layout string
		tokens int
	}{
		{time.RFC3339Nano, 1},
		{time.RFC3339, 1},
		{"2006-01-02 15:04:05.999999999", 2},
		{"2006-01-02 15:04:05", 2},
	}

	for _, c := range candidates {
		if len(fields) < c.tokens {
			continue
		}
		prefix := strings.Join(fields[:c.tokens], " ")
		// Python's logging module uses a comma as the fractional-second
		// separator (e.g. "11:48:21,235"); Go's time.Parse only understands a
		// period, so normalize it before parsing.
		normalized := strings.Replace(prefix, ",", ".", 1)
		t, err := time.Parse(c.layout, normalized)
		if err != nil {
			continue
		}
		// Strip the consumed tokens from the original line to get the message.
		rest := line
		for _, tok := range fields[:c.tokens] {
			rest = strings.TrimPrefix(strings.TrimSpace(rest), tok)
		}
		return t, strings.TrimSpace(rest), true
	}

	return time.Time{}, line, false
}

// isNeedle reports whether the count-th haystack line (1-based) should be the
// rare NEEDLE marker rather than noise. every<=0 disables needles entirely.
func isNeedle(count, every int) bool {
	return every > 0 && count%every == 0
}

// haystack floods noise lines at a target rate, injecting a rare unique NEEDLE
// line every `every` lines. This reproduces the needle-in-a-haystack search
// case: a search for "NEEDLE" finds only sparse matches (or none, for an absent
// term) and forces Dozzle to walk the whole log backward, which is slow.
func haystack(rate, every int) {
	noise := strings.Split(randomData, ". ")
	// emit in 10ms ticks so the rate stays smooth instead of one big burst
	perTick := rate / 100
	if perTick < 1 {
		perTick = 1
	}
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	count := 0
	for range ticker.C {
		for range perTick {
			count++
			if isNeedle(count, every) {
				fmt.Fprintf(os.Stderr, "NEEDLE rare-event-%d at line %d\n", count/every, count)
			} else {
				fmt.Fprintln(os.Stderr, noise[rand.Intn(len(noise))])
			}
		}
	}
}

// replay reads timestamped logs from stdin and replays them with original timing.
// When showNumbers is true each replayed line is prefixed with its index; it is
// off by default so the output stays byte-for-byte identical to the input (and
// therefore still valid JSON when the messages are JSON).
func replay(speedFactor float64, showNumbers bool) {
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
			adjustedDuration = min(adjustedDuration, 3*time.Second)
			time.Sleep(adjustedDuration)
		}
		if showNumbers {
			fmt.Fprintf(os.Stderr, "(%d) %s\n", i, entries[i].message)
		} else {
			fmt.Fprintln(os.Stderr, entries[i].message)
		}
	}

	time.Sleep(100 * time.Minute)
}

const randomData = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst. Praesent sapien turpis, fermentum vel, eleifend faucibus, vehicula eu, lacus. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst. Praesent sapien turpis, fermentum vel, eleifend faucibus, vehicula eu, lacus. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Curabitur pretium tincidunt lacus. Nulla gravida orci a odio. Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris. Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula. Donec lobortis risus a elit. Etiam tempor. Ut ullamcorper, ligula eu tempor congue, eros est euismod turpis, id tincidunt sapien risus a quam. Maecenas fermentum consequat mi. Donec fermentum. Pellentesque malesuada nulla a mi. Duis sapien sem, aliquet nec, commodo eget, consequat quis, neque. Aliquam faucibus, elit ut dictum aliquet, felis nisl adipiscing sapien, sed malesuada diam lacus eget erat. Cras mollis scelerisque nunc. Nullam arcu. Aliquam consequat. Curabitur augue lorem, dapibus quis, laoreet et, pretium ac, nisi. Aenean magna nisl, mollis quis, molestie eu, feugiat in, orci. In hac habitasse platea dictumst. Praesent sapien turpis, fermentum vel, eleifend faucibus, vehicula eu, lacus.`

func main() {
	random := flag.Bool("r", false, "generate random data")
	burst := flag.Int64("b", -1, "generate large burst of data")
	sleep := flag.Int64("s", 1000, "sleep time")
	shuffle := flag.Bool("x", false, "shuffle data")
	numbers := flag.Bool("n", false, "show line numbers (also prefixes replayed lines in -p mode)")
	all := flag.Bool("a", false, "print all data and pause")
	playback := flag.Float64("p", 0, "replay logs with original timing (speed factor: 1=normal, 10=10x faster, 0=disabled)")
	hay := flag.Bool("haystack", false, "flood noise logs with a rare NEEDLE line (needle-in-a-haystack search test)")
	rate := flag.Int("rate", 20000, "with -haystack, noise lines per second")
	every := flag.Int("every", 50000, "with -haystack, emit one NEEDLE line every N lines")
	flag.Parse()

	// Handle replay mode
	if *playback > 0 {
		replay(*playback, *numbers)
		return
	}

	// Handle needle-in-a-haystack mode
	if *hay {
		haystack(*rate, *every)
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
