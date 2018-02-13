package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kreativka/gostempel/egothor"
	"github.com/kreativka/gostempel/javautf"
)

func main() {
	f, err := os.Open("./testData")
	// f, err := os.Open("./output.txt")
	if err != nil {
		log.Fatalln(err)
	}

	fi, err := os.Open("./stemmer_20000.tbl")
	if err != nil {
		log.Fatalln(err)
	}

	in := bufio.NewReader(fi)
	stemmer, err := loadStemmer(in)

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalln("cannot load benchmark data")
		}
	}()

	defer func() {
		err := fi.Close()
		if err != nil {
			log.Fatalln("cannot load stemmer table")
		}
	}()

	var words []string
	i := bufio.NewScanner(f)
	for i.Scan() {
		t := strings.Split(i.Text(), "\n")
		words = append(words, t[0])
	}

	startTime := time.Now()
	var r string
	for _, w := range words {
		r = stem(stemmer, w)
	}
	stemDuration := time.Since(startTime)
	stemDurationSeconds := float64(stemDuration) / float64(time.Second)
	timePerWord := float64(stemDuration) / float64(len(words))
	log.Printf("Stemmed %d words, in %.2fs (average %.2fµs/word)",
		len(words),
		stemDurationSeconds,
		timePerWord/float64(time.Microsecond),
	)
	_ = r
}

func stem(stem *egothor.MultiTrie, token string) string {
	// min length for testing set to 0
	if token == "" || utf8.RuneCountInString(token) <= 0 {
		return token
	}

	cmd := stem.GetLastOnPath(token)
	if cmd == nil {
		return token
	}

	res := egothor.DiffApply(token, cmd)
	if len(res) > 0 {
		return res
	}

	return token
}

func loadStemmer(in *bufio.Reader) (*egothor.MultiTrie, error) {
	// Read method from file
	// Do nothing, as we only read standard table
	// , and we know how to create Trie
	// method, err := javautf.ReadUTF(in)
	_, err := javautf.ReadUTF(in)
	if err != nil {
		log.Fatalln("failed to read method from file")
	}
	mt := egothor.NewMultiTrie(in)
	return mt, nil
	// if method[0] != 'M' {
	// }
	// return nil, errors.New("error reading table")
}
