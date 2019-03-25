/* Author: Concurrer */

/* Problem:
Question Courtesy: https://www.careercup.com/question?id=5750868554022912

There are N places in the 1-dimensional parking lot. There are N-1 cars parked; each parked car is denoted via a capital letter such as E, T, etc. The parking lot has an empty slot denoted via '_' character. Use '_' to swap between any parked car and empty space. Given an initial configuration and final configuration; what would be the optimal number of swaps needed and provide the moves. Also return -1 if final configuration is not possible to reach.

Initial Configuration
 {'E', 'A', 'C', 'Q', '_', 'W', 'T' }

Final Configuration
 {'A', 'Q', 'E', 'T', 'C', 'W', '_' }

*/

/*
Solution: By looking at the 'target', swap the current elements keeping track of the '_' index all the time
*/

package main

import "fmt"

func swap(current []rune, D int, T int) {
	temp := current[D]
	current[D] = current[T]
	current[T] = temp
}

func fix(current []rune, currentMap map[rune]int, D int, r rune) (swaps int) {
	// D=i --  fix the rune at current index .. Ex: at 0 index, we need A
	S := currentMap[r]   // but A is at '1' in the currentMap
	T := currentMap['_'] // index of '_' empty space

	// lets swap the currentMap elements in D and S using the empty space in T
	if S != D { // of course only if a rune is not already in right place
		if D != T { // move the current car to empty space Ex: E->_ .. no need to swap if current space itself is empty
			//swap(D,T)
			currentRune := current[D]   // E at current[0]
			currentMap[currentRune] = T // move E to empty space
			currentMap['_'] = D
			swap(current, D, T) // swap runes on current too
			swaps++
		}
		//swap(D,S)
		currentMap[r] = D
		currentMap['_'] = S
		swap(current, D, S) // move the intended car to its right place A->0
		swaps++
	}
	return
}

func main() {
	target := []rune{'A', 'Q', 'E', 'T', 'C', 'W', '_'}
	current := []rune{'E', 'A', 'C', 'Q', '_', 'W', 'T'}

	numberOfSwaps := 0

	// build the current map and also get the index of '_' in current
	T := -1
	currentMap := make(map[rune]int, len(current))
	for i := 0; i < len(current); i++ {
		currentMap[current[i]] = i // E->0, A->1, C->2 ...
		if current[i] == '_' {
			T = i // save the '_' index
		}
	}

	switch {
	case T == -1: // no empty space to initiate swaps
		fmt.Println(-1)
	case len(current) != len(target): // current and target lengths are not equal
		fmt.Println(-1)
	}

	fmt.Println("before")
	fmt.Println(target)
	fmt.Println(current)

	for i := 0; i < len(target); i++ { // iterate over the target until all runes are crosschecked in place
		numberOfSwaps += fix(current, currentMap, i, target[i]) // return number of swaps if target[i] is not already there at the index i
	}

	fmt.Println("after")
	fmt.Println(target)
	fmt.Println(current)
	fmt.Printf("\nnumberOfSwaps:%d\n", numberOfSwaps)
}
