package main

import (
	"fmt"

	"github.com/mmcdole/wordnet"
)

func main() {
	// Initialize WordNet
	err := wordnet.Initialize("/usr/share/wordnet")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get input word
	fmt.Print("Enter a word: ")
	var input string
	fmt.Scanln(&input)

	// Find similar words
	similarWords := make([]string, 0)
	synsets, err := wordnet.LookupSynsets(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, synset := range synsets {
		for _, lemma := range synset.Lemmas {
			if lemma.Lemma != input {
				similarWords = append(similarWords, lemma.Lemma)
			}
		}
	}

	// Print similar words
	fmt.Println("Similar words:")
	for _, word := range similarWords {
		fmt.Println(word)
	}
}
