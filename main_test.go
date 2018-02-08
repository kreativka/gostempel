package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/kreativka/reader/egothor"
)

type sTest struct {
	in  string
	out string
}

var words []sTest

func makeTests() {
	f, err := os.Open("./output.txt")
	if err != nil {
		log.Fatalln(err)
	}

	fi, err := os.Open("./stemmer_20000.tbl")
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err := fi.Close()
		if err != nil {
			log.Fatalln("cannot read stem table")
		}
	}()

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalln("cannot load test data")
		}
	}()

	i := bufio.NewScanner(f)
	for i.Scan() {
		t := strings.Split(i.Text(), "\t")
		words = append(words, sTest{in: t[0], out: t[1]})
	}

	in := bufio.NewReader(fi)
	stemmer = loadStemmer(in)
}

var stemmer *egothor.MultiTrie

func TestStem(t *testing.T) {
	makeTests()
	for _, tt := range words {
		s := stem(stemmer, tt.in)
		if s != tt.out {
			t.Errorf("stem(%q) => %q, want %q", tt.in, s, tt.out)
		}
	}
}
