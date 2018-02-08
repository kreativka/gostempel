package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/kreativka/reader/egothor"
	"github.com/kreativka/reader/javautf"
)

func main() {
	f, err := os.Open("./stemmer_20000.tbl")
	if err != nil {
		log.Fatalln(err)
	}

	in := bufio.NewReader(f)
	stemmer := loadStemmer(in)

	err = f.Close()
	if err != nil {
		log.Fatalln(err)
	}

	terms := []string{"aachen", "wkręcający", "chodzący", "walnięty", "biegnący", "biegnę", "biegną", "kochają", "żółtość", "profilaktyczny"}

	for _, term := range terms {
		fmt.Println(stem(stemmer, term))
	}

}

func stem(stem *egothor.MultiTrie, token string) string {
	// min length for testing set to 0
	if token == "" || utf8.RuneCountInString(token) <= 0 {
		return token
	}

	cmd := stem.GetLastOnPath(token)
	if cmd == "" {
		return token
	}

	res := egothor.DiffApply(token, cmd)
	if len(res) > 0 {
		return res
	}

	return token
}

func loadStemmer(in *bufio.Reader) *egothor.MultiTrie {
	tries := egothor.NewMultiTrie()

	// Read method from file
	// Do nothing, as we only read, and stem
	// from given table
	_, err := javautf.ReadUTF(in)
	if err != nil {
		log.Fatalln("failed to read method from file")
	}

	tries.SetForward(readBool(in))
	// Read BY
	readInt(in)

	// Add tries
	for i := readInt(in); i > 0; i-- {
		// MultiTrie2 root trie
		t := egothor.NewTrie()
		t.SetForward(readBool(in))
		// Read root
		readInt(in)
		// Set commands
		for j := readInt(in); j > 0; j-- {
			cmd, err := javautf.ReadUTF(in)
			if err != nil {
				log.Println(err)
			}
			// Append cmd
			t.AddCmds(cmd)
		}
		// Add rows
		for j := readInt(in); j > 0; j-- {
			row := egothor.NewRow()
			// Add cells
			for k := readInt(in); k > 0; k-- {
				ch := readChar(in)
				cmd := readInt(in)
				cnt := readInt(in)
				ref := readInt(in)
				skip := readInt(in)
				cell := egothor.NewCell(ref, cmd, cnt, skip)
				row.AddCell(ch, cell)
			}
			// Append row
			t.AddRow(row)
		}
		// Append trie
		tries.AddTrie(t)
	}
	return tries
}
