package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"log"
	"os"
	"sort"
	"time"
)

// For limiting number of simultaniously running goroutines we use concurrencyLimiter
// as semaphore, that permits 70000 connection in one time (to prevent OS freezes)
const concurrency int = 70000

var concurrencyLimiter = make(chan struct{}, concurrency)

// A little bit optimized Levenshtein algorithm, so it uses O(min(m,n))
// space instead of O(mn), where m and n - lengths of compared strings.
// The key observation is that we only need to access the contents
// of the previous column when filling the matrix column-by-column.
// Hence, we can re-use a single column over and over, overwriting its contents as we proceed.
func LevenshteinDistance(from, to string, c chan Word) {
	var prevDiagonalValue, buffer int
	fromLength := len(from)
	toLength := len(to)

	// Initialize column
	var curColumn = make([]int, fromLength+1)
	for i := 0; i < fromLength; i++ {
		curColumn[i] = i
	}

	// Fill matrix column by column
	for i := 1; i <= toLength; i++ {
		curColumn[0] = i
		prevDiagonalValue = i - 1
		for j := 1; j <= fromLength; j++ {
			// Set operation cost (all operations except match(M) has value 1)
			operationCost := 1
			if from[j-1] == to[i-1] {
				operationCost = 0
			}
			buffer = curColumn[j]
			curColumn[j] = minOfThree(curColumn[j]+1, curColumn[j-1]+1, prevDiagonalValue+operationCost)
			prevDiagonalValue = buffer
		}
	}

	// Return value
	c <- Word{to, curColumn[fromLength]}
	// Decrease number of running goroutines
	<-concurrencyLimiter
}

func run(startWord string, fileName string, quantity int) Words {
	var words Words
	c := make(chan Word)

	go func() {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Read file word by word
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			concurrencyLimiter <- struct{}{}
			go LevenshteinDistance(startWord, scanner.Text(), c)
		}
		// Handle scanner errors
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()


	for i := 0; i < quantity; i++ {
		words = append(words, <- c)
	}
	sort.Sort(words)
	return words
}

func main() {
	wordsQuantity := flag.Int("n", 1000000, "Number of random words in test set.")
	startWord := flag.String("word", "test", "Program will calculate Levenshtein distance from this word.")
	flag.Parse()

	// Create test set, if it was not previously created
	fileName := fmt.Sprintf("test%v.txt", *wordsQuantity)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		generateTestFileWithLength(*wordsQuantity)
	}

	// Run profiler to measure time and memory consumption
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	start := time.Now()
	run(*startWord, fileName, *wordsQuantity)
	elapsedTime := time.Since(start)
	fmt.Println(*wordsQuantity, "words has been sorted succesfully. It took", elapsedTime)
}
