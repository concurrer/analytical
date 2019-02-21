/* Author: concurrer */
/* Problem: 

You are given an integer array nums and you have to return a new counts array. The counts array has the property where counts[i] is the number of smaller elements to the right of nums[i].

Input: [5,2,6,1]
Output: [2,1,1,0]
Explanation:
To the right of 5 there are 2 smaller elements (2 and 1).
To the right of 2 there is only 1 smaller element (1).
To the right of 6 there is 1 smaller element (1).
To the right of 1 there is 0 smaller element.

Solution: 
[5 7 2 6 1]
[2 3 1 1 0]
*/


package main

import (
	"fmt"
	"sync"
)

type pair [2]int

func main() {
	input := []int{5, 7, 2, 6, 1}
	l := len(input)
	result := make([]int, l)
	parentChan := make(chan pair, l) // for closing defer
	cellChans := make([]chan pair, 0)
	var wg sync.WaitGroup

	//create individual channel for each index in the input
	for i := 0; i < len(input); i++ {
		tempChan := make(chan pair, l)
		cellChans = append(cellChans, tempChan)
	}

	wg.Add(1)
	go initAction(input, parentChan, cellChans, result, &wg)
	countSmaller(input, parentChan)

	//display input and result slices

	wg.Wait()
	fmt.Printf("%v\n", input)
	fmt.Printf("%v\n", result)
}

func initAction(input []int, parentChan chan pair, cellChans []chan pair, result []int, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()
	// launch the parent Channel
	wg.Add(1)
	go parentListener(parentChan, cellChans, wg)

	// launch the cell Channels
	for i, v := range input {
		wg.Add(1)
		cell := pair{i, v}
		go cellListener(cell, cellChans, result, wg)
	}
}

/*
   parentListener listens for the incoming pair and broadcasts that pair to all the cellChans
*/
func parentListener(parentChan chan pair, cellChans []chan pair, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()
	for {
		select {
		case incoming := <-parentChan:
			// broadcast the cell to all cellChans
			for i, _ := range cellChans {
				cellChans[i] <- incoming
			}
			if incoming[0] == -1 {
				return
			}

		default: // just wait
		}
	}
}

/*
   each cellChan waits on the parentChan for a pair. For each incoming pair
   * check if the index in the pair is 'after' the self
        * if yes, check if the number is less than the self's value
            * if yes, increment the value in result[cell] slice
   * if a sig pair {-1,-1} received, then return
*/
func cellListener(cell pair, cellChans []chan pair, result []int, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()
	cellIndex, cellValue := cell[0], cell[1]
	for {
		select {
		case incoming := <-cellChans[cellIndex]:
			if incoming[0] > cellIndex && incoming[1] < cellValue { //the incoming pair is 'after' the self AND its value is less than self
				result[cellIndex]++
			} else if incoming[0] == -1 {
				return
			}
		default: // just wait
		}
	}

}

// countSmaller sends all the (index,value) pairs to the parentChan for further broadcasting
func countSmaller(input []int, parentChan chan pair) {
	for i, v := range input {
		parentChan <- pair{i, v}
	}
	//close parent Channel
	defer func() { c := pair{-1, -1}; parentChan <- c }()
}
