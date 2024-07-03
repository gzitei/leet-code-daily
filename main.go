package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"leet/api"
)

func Usage() {
	fmt.Println("Usage: leet [options]")
	fmt.Println("Options:")
	fmt.Println("  --dir string       defines destination directory (required)")
	fmt.Println("  --random           selects a random problem (optional)")
	fmt.Println("  --slug string      selects a problem according to provided slug (optional)")
	fmt.Println("Note: You must provide the --dir flag. If no problem flag is provided, it defaults to the daily challenge.")
}

func main() {
	flag.Usage = Usage

	directory := flag.String("dir", "", "defines destination directory")
	randomProblem := flag.Bool("random", false, "selects a random problem")
	slug := flag.String("slug", "", "selects a problem according to provided slug")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Error: You must inform a root directory using --dir flag.")
		flag.Usage()
		return
	}

	dir, err := filepath.Abs(*directory)
	if err != nil {
		fmt.Println("Error getting specified directory:", err)
		return
	}

	var question api.Question
	switch {
	case *randomProblem:
		question, err = api.GetRandomQuestion()
	case *slug != "":
		question, err = api.GetQuestionBySlug(*slug)
	default:
		question, err = api.GetDailyChallenge()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	var d string
	if d, err = api.SetUpEnv(dir, question); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d)
}
