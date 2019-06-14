# gostempel

[![Build Status](https://cloud.drone.io/api/badges/kreativka/gostempel/status.svg)](https://cloud.drone.io/kreativka/gostempel)
[![GoDoc](https://godoc.org/github.com/kreativka/gostempel?status.svg)](https://godoc.org/github.com/kreativka/gostempel)
[![Go Report Card](https://goreportcard.com/badge/github.com/kreativka/gostempel)](https://goreportcard.com/report/github.com/kreativka/gostempel)

Stempel port to golang.
It's the fastest golang implementation you can get.

## Usage

To play with gostempel:

stemmer, _ := gostempel.LoadStemmer("./stemmer_20000.tbl")  
stemmedWord := gostempel.Stem(stemmer, "polishWordToStem")

## License

This product includes software developed by the Egothor Project. http://egothor.sf.net/
