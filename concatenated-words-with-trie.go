/* Author: Concurrer */
/*
Question courtesy:  https://leetcode.com/problems/concatenated-words/

Given a list of words (without duplicates), please write a program that returns all concatenated words in the given list of words.

A concatenated word is defined as a string that is comprised entirely of at least two shorter words in the given array.

Example:

Input: ["cat","cats","catsdogcats","dog","dogcatsdog","hippopotamuses","rat","ratcatdogcat"]

Output: ["catsdogcats","dogcatsdog","ratcatdogcat"]

Explanation: "catsdogcats" can be concatenated by "cats", "dog" and "cats";
 "dogcatsdog" can be concatenated by "dog", "cats" and "dog";
"ratcatdogcat" can be concatenated by "rat", "cat", "dog" and "cat".

Note:

    The number of elements of the given array will not exceed 10,000
    The length sum of elements in the given array will not exceed 600,000.
    All the input string will only include lower case letters.
    The returned elements order does not matter.
*/

/*
Solution:
1. build a trie with all the chars in all words in the input
2. mark each valid word with a bool=true
3. search function returns if it can find the input as a valid word (with ending bool=true)
4. compund function splits each word into all possible pieces and checks if those pieces exist as valid words in the trie
*/

package main

import "fmt"

type node struct {
	value      string
	downstream map[string]node
	end        bool
}

func makeNode(s string) (n node) {
	return node{value: s, end: false, downstream: make(map[string]node, 0)}
}

func isValidWord(head node, word string) (isAvailable bool) {
	//isAvailable = false // this is the default value anyway
	if len(word) == 0 {
		return
	}

	for i := 0; i < len(word); i++ {
		downstream := head.downstream // head is  at the top null in the beginning. get the trie and tread down with each iteration
		// check if the char exists and traverse down the path if yes
		char := word[i : i+1]              // get one char of the word
		if _, ok := downstream[char]; ok { // if the letter is found then continue to next letter, if not return
			head = downstream[char]
		} else {
			return
		}
	}
	if head.end == true {
		isAvailable = true
	}
	return // naked return
}

func searchCompoundWord(head node, word string) (isCompound bool) {
	//isCompound = false // this is the default value anyway
	wl := len(word)
	if wl == 0 {
		return
	}

	p := head // save the head

	for i := 0; i < wl; i++ {
		if isValidWord(head, word[0:i+1]) { // if the prefix is a valid word then check if the remaining string is a valid word too
			head = p
			if i < wl-1 && isValidWord(head, word[i+1:]) { // not at the end of a valid word already and there is a remaining string which is also valid
				isCompound = true
				return
			} // can't have else here as head gets reset
			head = p
			if i < wl-1 && searchCompoundWord(head, word[i+1:]) { // its a single valid word so far, lets check if the remaining is a compound word
				isCompound = true
				return
			}
		}
		head = p // reset head to the top
	}
	return // naked return
}

func findAllConcatenatedWordsInADict(words []string) (output []string) {
	head := node{value: "", downstream: make(map[string]node, 0), end: false}
	p := head                    // save head
	for _, word := range words { // get one word
		for i := 0; i < len(word); i++ {
			char := word[i : i+1]               // get one letter of the word
			downstream := head.downstream       // get the downstream trie
			if n, ok := downstream[char]; !ok { // if the letter is not already there, make a new node and add it
				n = makeNode(char)
				if i == len(word)-1 { // mark the end of a valid word
					n.end = true
				}
				downstream[char] = n // put the new node into the trie
			} else { // the char exists.. we need to mark the end of word if 'i' at last char
				if i == len(word)-1 {
					n.end = true
					downstream[char] = n
				}
			}
			head = downstream[char] // move the node to next if the letter is already in the downstream
		}
		head = p // reset node for new word
	}

	// for each word in 'input' check if it exists as a compound word in the trie
	for _, word := range words {
		head = p // reset head
		if isCompound := searchCompoundWord(head, word); isCompound {
			output = append(output, word)
		}
	}
	return output
}

func main() {
	words := []string{"cat", "cats", "catsdogcats", "dog", "dogcatsdog", "hippopotamuses", "rat", "ratcatdogcat"}
	fmt.Printf("\n%+v\n", findAllConcatenatedWordsInADict(words))
}
