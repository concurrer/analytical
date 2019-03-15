/* Author: Concurrer */
/*
Assume we have a long list of words from English dictionary, like:

words: ["water", "chair", "slide", "paris", "amsterdam", "air", "plane", "bottle", "book", "lift", "hackathon" ...]

And another long list of strings to process, write a function to identify "compound words" and return them:

input:[ "paris", "waterslide", "airplane", "amsterdam", "chairlift", "apple", "planet" … ]

output: [ "waterslide", "airplane", "chairlift" …]

The definition of "compound word":

x is a compound word if-and-only-if
x does not exist in EnglishWords, and
x = w1 + w2 for some w1,w2 in EnglishWords

*/

/*
Solution:
    * build a trie with all words and mark an 'end=true' for each valid word
    * iterate over the input words and split the search every time we hit a valid 'end=true' (increment a counter 'numWords' for each split)
    * if a word exits with a clean end and multiple splits, then its a compound word. Ex: 'paris' is not but 'airplane' is
*/

package main

import (
	"fmt"
)

type node struct {
	value      string
	downstream map[string]node
	end        bool
}

func makeNode(s string) (n node) {
	return node{value: s, end: false, downstream: make(map[string]node, 0)}
}

func searchWord(head node, word string) (isPrefixWholeWord, hasTail bool, index int) {
	isPrefixWholeWord, hasTail, index = false, false, 0
	var i int
	for i = 0; i < len(word); i++ {
		// check what part of the word already exists in the trie
		char := word[i : i+1] // get one letter of the word
		downstream := head.downstream
		if _, ok := downstream[char]; ok { // if the letter is found then continue to next letter, if not check for the word end
			head = downstream[char]
			if head.end == true {
				isPrefixWholeWord = true
				if i < len(word)-1 { // we are in the middle of the string but the char for loop is out
					hasTail = true
					index = i + 1 // send the index back to search for the remaining word afresh.. string upto 'i' is already read, so send i+1
				}
				return
				//break
			}
			continue
		}
	}
	return // naked return
}

func searchCompoundWord(head node, word string) (isCompound bool) {
	isCompound = false
	numWords := 0
	isPrefixWholeWord, hasTail, index := false, false, 0
	for isPrefixWholeWord, hasTail, index = searchWord(head, word); isPrefixWholeWord && hasTail; {
		// if a part-word matches, then save the index and check to see if
		// the remaining string is another whole word
		numWords++
		isPrefixWholeWord, hasTail, index = searchWord(head, word[index:])
	}

	if isPrefixWholeWord && (!hasTail) && (numWords > 0) { // no more string and its a isPrefixWholeWord in every iteration and its not a single word but at least two
		fmt.Println("\nsetting isCompound to true\n")
		isCompound = true
	}
	return // naked return
}
func main() {
	words := []string{"water", "chair", "slide", "paris", "amsterdam", "air", "plane", "bottle", "book", "lift", "hackathon"}
	input := []string{"paris", "waterslide", "airplane", "amsterdam", "chairlift", "apple", "planet"}
	var output []string

	// build a trie with the words
	head := node{value: "", downstream: make(map[string]node, 0), end: false}
	// fill head with 26 letters in the alphabet

	// feed all 'words' into the trie
	p := head                    // iterator from head
	for _, word := range words { // get one word
		for i := 0; i < len(word); i++ {
			char := word[i : i+1]               // get one letter of the word
			downstream := p.downstream          // get the downstream trie
			if _, ok := downstream[char]; !ok { // if the letter is not already there, make a new node and add it
				n := makeNode(char)
				if i == len(word)-1 { // mark the end of a valid word
					n.end = true
				}
				downstream[char] = n // put the new node into the trie
			}
			p = downstream[char] // move the node to next if the letter is already in the downstream
		}
		p = head // reset node for new word
	}

	// for each word in 'input' check if it exists as a whole word in the trie
	for _, word := range input {
		if isCompound := searchCompoundWord(head, word); isCompound {
			output = append(output, word)
		}
	}

	fmt.Printf("\noutput=%+v", output)
}
