package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	info := color.New(color.FgHiGreen).PrintfFunc()
	instruction := color.New(color.Bold, color.FgMagenta).PrintfFunc()
	question := color.New(color.FgHiGreen).PrintfFunc()
	
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 20, "the time limit for the quiz in seconds")
	flag.Parse()
	
	file, err := os.Open(*csvFilename)	
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFilename))
	}
	
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	
	instruction("Answer as many problems as possible within the time limit! (%d seconds)\n\n", *timeLimit)
	time.Sleep(time.Second)
	for i := 3; i > 0; i-- {
		time.Sleep(time.Second) 
		info("%d\n", i)
	}
	time.Sleep(time.Second) 
	instruction("Go !\n\n")
	time.Sleep(time.Second)

	problems := parseLines(lines)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	
	// using labels for breaking loop is not recommended for bigger app
	problemLoop:
		for i, p := range problems {
			question("Problem #%d: %s = ", i+1, p.q)
			answerCh := make(chan string)
			go func() {
				var answer string
				fmt.Scanf("%s\n", &answer)
				answerCh <- answer
			}()

			select {
			case <-timer.C:
				fmt.Print("\n")
				break problemLoop
			case answer := <- answerCh:
				if answer == p.a {
					correct++
				}
			}
	}

	instruction("\nYou scored %d out of %d.\n", correct, len(problems))
	if correct == len(problems) {
		info("\nPerfect Score! Congratulations!")
	}
}

func parseLines(lines [][]string) []problem {
	// when we know the length of a slice, there is no reason to let append() doing extra work to resize it 
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string)  {
	error := color.New(color.Bold, color.FgHiRed).PrintfFunc()
	error(msg + "\n")
	os.Exit(1)
}