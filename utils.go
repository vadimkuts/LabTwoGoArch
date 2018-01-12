package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randomWord(n int) string {
	// Got code from here: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func minOfThree(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		} else {
			return c
		}
	} else {
		if b < c {
			return b
		} else {
			return c
		}
	}
}

// Generate file filled with n random words.
func generateTestFileWithLength(n int) string {

	fileName := fmt.Sprintf("test%v.txt", n)
	file, err := os.Create(fileName)
	// Handle file open error
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var randowWordLength int
	for i := 0; i < n; i++ {
		randowWordLength = rand.Intn(10) + 4
		file.WriteString(randomWord(randowWordLength+1) + "\n")
		if math.Mod(float64(i), 1000000) == 0 && i != 0 {
			fmt.Println(i, "words has been created")
			// To check that program don't stuck
		}
	}
	fmt.Println(n, "words has been created")

	return fileName
}

// Got from here: http://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently
// Used for tests
func lineCounter(file io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSeparator := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSeparator)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
