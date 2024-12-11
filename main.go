package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "csv.file", "a csv file in the format of 'question,anwere' ")
	timeLimit := flag.Int("limit", 3, "the time limit for the quiz in second")
	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {
		myErrorStr := fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename)
		exit(myErrorStr)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("failed to parse the provide CSV file.")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemsloop:
	for i, p := range problems {
		fmt.Printf("problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			_, e := fmt.Scan(&answer)
			if e != nil {
				fmt.Println("error")
			}
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemsloop
		case answer := <-answerCh:
			fmt.Println()
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("you socred %d out of %d.\n", correct, len(problems))
}

func scanf(str *string) (int, error) {
	*str = "hello0"
	return len(*str), nil
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			q: line[0],
			a: line[1],
		}
	}

	return result
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
