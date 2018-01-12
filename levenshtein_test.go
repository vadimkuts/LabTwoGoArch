package main

import (
	"log"
	"os"
	"testing"
)

type testpairMin struct {
	values []int
	min    int
}

type testpairRandom struct {
	length         int
	expectedLength int
}

type testpairLines struct {
	count         int
	expectedCount int
}

type testpairLevenshteinDistance struct {
	from     string
	to       string
	distance int
}

var testsMinOfThree = []testpairMin{
	{[]int{1, 2}, 1},
	{[]int{50, 100}, 50},
	{[]int{-1, 1}, -1},
}

var testsRand = []testpairRandom{
	{1, 1},
	{10, 10},
	{3, 3},
}

var testsGenerate = []testpairLines{
	{10, 10},
	{15, 15},
	{52, 52},
}

var words = []Words{
	{{"test", 4},
		{"testa", 5},
		{"testab", 6},
		{"testabc", 7},
		{"testabcd", 8},
		{"testabcde", 9}},
}

var testLevenshteinDistance = []testpairLevenshteinDistance{
	{"test", "test", 0},
	{"test", "testa", 1},
	{"test", "testab", 2},
	{"test", "testabc", 3},
	{"test", "testabcd", 4},
	{"test", "testabcde", 5},
}

var testRun = Words{
	{"test", 0},
	{"testa", 1},
	{"testab", 2},
	{"testabc", 3},
	{"testabcd", 4},
	{"testabcde", 5},
}

func TestMinOfThree(t *testing.T) {
	for _, pair := range testsMinOfThree {
		result := minOfThree(pair.values[0], pair.values[1], pair.min)
		if result != pair.min {
			t.Error(
				"For", pair.values,
				"expected", pair.min,
				"got", result,
			)
		}
	}
}

func TestRandomWord(t *testing.T) {
	for _, pair := range testsRand {
		result := randomWord(pair.length)
		if len(result) != pair.expectedLength {
			t.Error(
				"For", pair.length,
				"expected", pair.expectedLength,
				"got", result,
			)
		}
	}
}

func TestGenerateTestFileWithLength(t *testing.T) {
	for _, pair := range testsGenerate {
		fileName := generateTestFileWithLength(pair.count)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		defer os.Remove(fileName)

		result, _ := lineCounter(file)
		if result != pair.expectedCount {
			t.Error(
				"For", pair.count,
				"expected", pair.expectedCount,
				"got", result,
			)
		}
	}
}

func TestLength(t *testing.T) {
	for _, pair := range words {
		result := pair.Len()
		if result != len(pair) {
			t.Error(
				"For", len(pair),
				"expected", len(pair),
				"got", result,
			)
		}
	}
}

func TestLess(t *testing.T) {
	for _, pair := range words {
		for i := range pair {
			result := pair.Less(0, i)

			if result != (pair[0].Distance < pair[i].Distance) {
				t.Error(
					"For", pair[0].Distance < pair[i].Distance,
					"expected", pair[0].Distance < pair[i].Distance,
					"got", result,
				)
			}
		}
	}
}

func TestLevensteinDistance(t *testing.T) {
	c := make(chan Word)
	for _, pair := range testLevenshteinDistance {
		concurrencyLimiter <- struct{}{}
		go LevenshteinDistance(pair.from, pair.to, c)
		result := <- c
		if result.Distance != pair.distance {
			t.Error(
				"For", pair.from, pair.to,
				"expected", pair.distance,
				"got", result,
			)
		}
	}
}

func TestRun(t *testing.T) {
	// Create test file
	result := run("test", "test.txt", 6)
	for i := range testLevenshteinDistance {

		if result[i].Distance != testRun[i].Distance || result[i].Text != testRun[i].Text {
			t.Error(
				"For", result[i].Text, testRun[i].Text,
				"expected", testRun[i].Distance,
				"got", result[i].Distance,
			)
		}
	}
}
