package main

import (
	"fmt"
	"time"

	"./backup"
)

func main() {
	// go backup.ReadDirectory("./", "")
	timeLayout := "Mon, 2 Jan 2006 15:04:05"
	fmt.Print("Keep patiences, we're so smart to finish things as fast as we can.\n\n")
	startTime := time.Now()
	fmt.Println("Started at", startTime.Format(timeLayout))

	backup.ReadConfig("./config.json")

	// var input string
	// fmt.Scanln(&input)
	endTime := time.Now()
	fmt.Print("\n\nDone at ", endTime.Format(timeLayout))
	duration := time.Since(startTime)
	fmt.Print("\nTotal time taken ", duration.Seconds(), " Seconds or, ", duration.Minutes(), " Minutes\n")

}
