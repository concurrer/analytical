/* Author: Concurrer */

/* 
Question courtesy: https://leetcode.com/problems/minimum-cost-to-merge-stones/

There are N piles of stones arranged in a row.  The i-th pile has stones[i] stones.

A move consists of merging exactly K consecutive piles into one pile, and the cost of this move is equal to the total number of stones in these K piles.

Find the minimum cost to merge all piles of stones into one pile.  If it is impossible, return -1.

Example-1:
Input: stones = [3,2,4,1], K = 2
Output: 20
Explanation:
We start with [3, 2, 4, 1].
We merge [3, 2] for a cost of 5, and we are left with [5, 4, 1].
We merge [4, 1] for a cost of 5, and we are left with [5, 5].
We merge [5, 5] for a cost of 10, and we are left with [10].
The total cost was 20, and this is the minimum possible.

Example-2:
Input: stones = [3,2,4,1], K = 3
Output: -1
Explanation: After any merge operation, there are 2 piles left, and we can't merge anymore.  So the task is impossible.

Example-3:
Input: stones = [3,5,1,2,6], K = 3
Output: 25
Explanation:
We start with [3, 5, 1, 2, 6].
We merge [5, 1, 2] for a cost of 8, and we are left with [3, 8, 6].
We merge [3, 8, 6] for a cost of 17, and we are left with [17].
The total cost was 25, and this is the minimum possible.
*/

/* 
Solution:

1. possible solution exists if and only if N%(K-1) equals 1 (at the end of last-but-one-merge, you must have K piles)
2. merging the 3 least piles gives the least cost for any merge.. global min is nothing but a chain of local mins 

In that sense, the example 3 given is wrong. 
Reducing [3,5,1,2,6] 
            -> [[1,2,3],5,6] (which means [[6],5,6]) // cost 6
            -> [17] // cost 17
.. is an overall cost of 6+17=23... so 25 is not the minimum as given in the example.  This also clarifies summing the mins in each iteration will lead to the overall min. 
*/

package main

import (
    "fmt"
    "sort"
)
func mergePiles(stones[]int, k int)([]int, int){
    sort.Ints(stones)
    total:=0

    for j:=0; j<k; j++ {
        total+=stones[j]
    }
    stones = stones[k-1:] // drop the first k-1  elements, index starts from 0 
    stones[0]=total // put the total in the 3rd slot
    return stones, stones[0]
}
func mergeStones(stones []int, k int)(int){
    totalStones := 0
    globalCost := 0
    for eachPile := range stones {
        totalStones+=eachPile
    }

    if len(stones)<k {  // num piles less than k.. ex:s={1},k=3
        return -1
    }

    if len(stones)%(k-1)!=1 && k>2{ // solution check-1
        return -1
    }
    for { // loop until the slice burns out
        if len(stones)==1 {
           break 
        } else {
            cost:=0
            stones, cost = mergePiles(stones,k)
            globalCost += cost
        }
    }
    return globalCost
}
func main() {
    stones := []int{3,2,4,1}
    k:=2
    fmt.Println(mergeStones(stones,k))
}
