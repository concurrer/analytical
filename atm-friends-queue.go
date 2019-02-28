/* Author: concurrer */

/* Question courtesy codechef: https://www.codechef.com/problems/CHN16D
  (modified inputs below to make the explanation simpler)

An entire class of students has gone to the ATM to withdraw cash. Let's name them {S1, S2, .., SN}. Obviously they have to stand in a queue. They want to do so, such that close friends are standing not too far from each other. Specifically, if Si and Sj are close friends, then there should not be more than 2 other students standing in between them. Given the close-friendship information, you need to tell whether it is possible to form such a queue, or not.
Formally, you should state whether there exists a permutation (Sp1, Sp2, .., SpN) of {S1, S2, .., SN}, such that for all pairs (Spi, Spj) where Spi and Spj are close friends, |i-j| ≤ 3 holds true.
Input
The first line contains two integers N, M, denoting the number of students, and the number of close-friendship pairs, respectively.
The next M lines contain two integers each, separated by single spaces: i j, denoting that Si and Sj are close friends.
Output
Output a single line for each testcase, containing the answer, which should be 1, if such a queue is possible, and 0 otherwise.
Constraints
1 ≤ T ≤ 3
1 ≤ N ≤ 30
1 ≤ M ≤ 200
Example
Input:
5 6
1 2
3 1
4 1
2 5
2 4
5 4
Output:
1
Explanation
Consider the permutation (S1, S3, S2, S4, S5). No pair of close friends have more than 2 students between them. So this is a valid queue, and hence the answer is 1.

/*

/* Solution:
lets quickly assume few variables here, say:
  S1 is friends with S2 and S3.
  S2 is friends with S3 and S4.
If we write it as a map of 'int' to a 'set', its roughly represented as :
1 -> {2,3} and 2 -> {1,3,4}

now, for the purpose of a solution, for every friend of '1', if the friend-of-friend is insidie the 1's own set, then lets call him 'insider'.. otherwise outsider.

so lets take the first friend of 1, i.e. 2.
2's friend 3 is an 'insider' for 1's set (because 3 is present in {2,3})... but 4 is an 'outsider'.


with this assumption in place, the way we form the result queue MUST follow these premises as per the problem description.

	1.  a student can't have more than 6 friends, in which case we bail out an impossible... because the distance must be less than or equal to 3 either side

	// now imagine a student in the center and all 6 positions are taken by his friends (3 on each side)
	// the friends in the two positions to his immediate next can only have 1 outsider and 5 insiders in their own sets at max. Lets call them o1i5
	// the friends in the two positions to his immediate next-to-next can only have 2 outsiders and 4 insiders in their own sets at max. Lets call them o2i4
	// the friends in the two outer positions can only have 3 outsiders and 3 insiders in their own sets at max. Lets call them o3i3

	// because if the inner friends have outsiders, then the outer friends will have lesser positions, so
	2. o2i4 <=  2 - o1i5
	3. o3i3 <=  2 - o2i4  - o1i5

*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// IntSet acts as a set of integers
type IntSet struct {
	set map[int]bool
}

func (intSet IntSet) Add(i int) {
	intSet.set[i] = true
}

func checkError(err error, s string) {
	if err != nil {
		fmt.Println(s)
		panic(err)
	}
}

func noSolution(s string) {
	panic(errors.New(s))

}
func main() {
	fname := "atm-input.txt"

	file, err := os.Open(fname)
	checkError(err, "ERROR: file open")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	checkError(err, "ERROR: file read")
	scanner.Scan()
	checkError(scanner.Err(), "ERROR: scanning line")

	// first line contains the number of students and the number of forthcoming pairs of friendships
	line := scanner.Text()
	s := strings.Split(line, " ")
	numStudents, err := strconv.Atoi(s[0])
	checkError(err, "ERROR: numStudents conversion")
	//numPairs, err := strconv.Atoi(s[1]) // this var is useless until we do multiple test cases.. but this code is only a poc, testing a single case
	//checkError(err, "ERROR: numPairs conversion")

	// map for holding the friendship pairs
	m := make(map[int]IntSet, 0)
	// initiate the map for numStudents and blank pairs
	// students start from 1->N... so to ignore 0 index with <=
	for i := 0; i <= numStudents; i++ {
		set := make(map[int]bool, 0)
		intSet := IntSet{set: set}
		m[i] = intSet
	}

	// read the rest of the lines for friendship pairs
	for scanner.Scan() {
		line = scanner.Text()
		s = strings.Split(line, " ")
		f1, err := strconv.Atoi(s[0])
		checkError(err, "ERROR: f1 conversion")
		f2, err := strconv.Atoi(s[1])
		checkError(err, "ERROR: f2 conversion")
		intSet := m[f1]
		intSet.Add(f2)
		// add the pair in reverse too
		intSet = m[f2]
		intSet.Add(f1)
	}
	// check scanner err
	checkError(scanner.Err(), "ERROR: scanning pairs")

	/* for every friend in the map 'm':
	* loop over the friend set
		* check if there are more than 6 friends, if yes, bail out
		* loop over the friends-of-friend set and calculate 'o1i5', 'o2i4' and 'o3i3'.
		* check ( o2i4 <=  2 - o1i5 ) AND (o3i3 <=  2 - o2i4  - o1i5)
			* if yes, return 1, else return 0
	*/

	// first for-loop gives the O(N)
	for student, intSet := range m {
		// if more than 6 friends return 0
		if len(intSet.set) > 6 {
			noSolution(fmt.Sprintf("ERROR: friend length > 6 for %d", student))
		}

		// loop over friends-of-friend, i.e. the "intSet of each friend" (NOT each student) and see how many are insiders and outsiders of the initial student intSet and increment the corresponding 'oNiN' variable
		o1i5, o2i4, o3i3 := 0, 0, 0

		// second for-loop is only for a max 6 (as per length check above).. hence its still O(N*constant)
		for friend := range intSet.set {
			insiders, outsiders := 0, 0
			// friends-of-friend
			fof_intSet := m[friend]
			// for each friend-of-friend, check the insider/outsider status
			// third for-loop is again for a max 6.. hence its still O(N*constant*constant)
			for fof := range fof_intSet.set {
				if intSet.set[fof] == true || fof == student {
					insiders++
				} else {
					outsiders++
				}
			}

			// fill the 'oNiN' variables
			//     4 3 2 1 5 6 7
			switch {
			case outsiders == 0 && insiders <= 5:
				// all good, continue
			case outsiders == 1 && insiders <= 5:
				o1i5++
			case outsiders == 2 && insiders <= 4:
				o2i4++
			case outsiders == 3 && insiders <= 3:
				o3i3++
			default: // no other case can be accommodated
				noSolution(fmt.Sprintf("ERROR: insiders/outsiders bail out for student=%d, friend=%d, insiders=%d, outsiders=%d", student, friend, insiders, outsiders))
			}

			// oNiN checks
			if !((o1i5 <= 2) && (o2i4 <= 2-o1i5) && (o3i3 <= 2-o1i5-o2i4)) {
				noSolution(fmt.Sprintf("ERROR: oNiN bail out for student=%d, friend=%d, o1i5=%d, o2i4=%d, o3i3=%d", student, friend, o1i5, o2i4, o3i3))
			}
		}
	}
	// all good so far.. its possible to form a queue
	fmt.Printf("%d", 1)
}
