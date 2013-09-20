package main

import "fmt"

var graph = [][]int{
		  {0,3,10,-1,-1},
		  {3,0,-1,5,-1},
		  {10,-1,0,6,15},
		  {-1,5,6,0,4},
		  {-1,-1,-1,4,0},
}

func out_graph() {
	for i := range graph {
		 for j := range graph[i] {
			 fmt.Printf("%d\t", graph[i][j])
		 }
		 fmt.Println()
	}
}

func treat_row(k int, i int) {
	for j := range graph[i] {
		 if graph[i][k] == -1 || graph[k][j] == -1 {
			  continue
		 }
		 temp := graph[i][k] + graph[k][j]
		 if graph[i][j] > temp || graph[i][j] == -1 {
			  graph[i][j] = temp
		 }
	}

}

func main() {
	 fmt.Println("The inputted graph:")
	 out_graph()



	for k := range graph {
		 for i := range graph {
			  go treat_row(k,i)
		 }
	}

	 fmt.Println("The shortest path in the graph:")
	 out_graph()
}



