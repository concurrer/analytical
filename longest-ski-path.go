/* Author: concurrer */

/*
Problem Description:

You're given a height map for a mountain, represented as a 2D array (matrix). When skiing, a skier can only move to adjacent (vertically/horizontally) tiles and go from higher elevation to lower elevation. Determine the longest possible ski path based on this rule. (Ties may be broken in arbitrary ways.)

Example input:

3 5 3

2 3 4

1 0 1

Output: path: 5 -> 3 -> 2 -> 1 -> 0, length: 5
The path 4->3->2->1->0 would also have been valid.



The following solution is NOT complete. This is only an attempt to see how concurrency can simplify a DP solution and also make the code easy to read.

OUTPUT
-------------------------------------:
END01:  5
END11:  5  ->  3
END21:  5  ->  3  ->  0
END00:  5  ->  3
END02:  5  ->  3
END10:  5  ->  3  ->  2
END20:  5  ->  3  ->  2  ->  1
END21:  5  ->  3  ->  2  ->  1  ->  0
END10:  5  ->  3  ->  2
END20:  5  ->  3  ->  2  ->  1
END21:  5  ->  3  ->  2  ->  1  ->  0

Pending work: 
    * to output only the longest path
    * to ensure multiple goroutines are not launched for the same node from different neighbors

*/

package main

import (
    "fmt"
    "sync"
)



func main() {
    matrix := [][]int{{3,5,3},{2,3,4},{1,0,1}}

    var wg sync.WaitGroup

    // for every node in matrix, launch go ski with path nil
    x, y := findTallest(matrix)
    wg.Add(1)
    go ski(matrix, x,y, "", &wg)

    // wait until all the goroutines are done
    wg.Wait()
}

/* ski takes the coordinates of each node, then
 * check if its having an adjacent node with lower elevation
 * if yes launch go ski on each such node sending along the path so far
 */
func ski(matrix [][]int, x int, y int, path string, wg *sync.WaitGroup){
    defer func(){wg.Done()}()

    currentElev := matrix[x][y]
    path = fmt.Sprintf("%s %d ",path, currentElev)

    for _,neighbor := range neighbors(x,y,len(matrix),len(matrix[0])){
        nx := neighbor[0]
        ny := neighbor[1]
        nElev := matrix[nx][ny]

        //if neighbor's elevation is lower, launch a go ski
        if nElev < currentElev {
            wg.Add(1)
            tpath:= fmt.Sprintf("%s -> ", path)

           go ski(matrix,nx,ny, tpath, wg)  
        } 
    }
    // end of a route
    fmt.Printf("\nEND%d%d: %s",x,y,path)
}

/*
findTallest is a convenience function to get the tallest peak to start with., to reduce the number of goroutines launched assuming every node will reach down of course :)
*/
func findTallest( matrix [][]int) (x int, y int){
    height := len(matrix)
    width := len(matrix[0])
    tall, x, y := 0,0,0
    for i:=0; i<height; i++ {
        for j:=0; j<width; j++  {
            if matrix[i][j] > tall {
                tall,x,y = matrix[i][j], i, j
            }
        }
    }
   return x,y 
}
/*
neighbors gives all the adjacent nodes given the coordinates of a single node
*/
func neighbors(x int, y int, height int, width int) ([][2]int){
    adjacent := make([][2]int,0)
    //check left of row
    if y-1 >= 0 {
        adjacent = append(adjacent, [2]int{x,y-1})
    }
    //check right of row
    if y+1 < width {
        adjacent = append(adjacent, [2]int{x,y+1})
    }
    //check the column above
    if x-1 >= 0  {
        adjacent = append(adjacent, [2]int{x-1,y})
    }
    //check the column below
    if x+1 < height  {
        adjacent = append(adjacent, [2]int{x+1,y})
    }
    return adjacent
}
