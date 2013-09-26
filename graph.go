package main

import (
	"fmt"
   "net/http"
	"strconv"
	"sync"
	)

type Page struct {
	 Title string
	 InvitInput string
	 Input  []byte
	 InvitOutput string
	 Output  []byte
}

var pd = Page {
	Title: "Finding the shortest path in a graph",
	InvitInput: "Input a matrix of the graph:",
	Input: []byte(""),
	InvitOutput: "Results:",
	Output: []byte("") }

var result bool = false

var message = []string{
	"Incorrect input: It's possible to input 1-9, -, <Enter> and <Space> only",
	"Incorrect input: '-' must be at the beginning of a number",
	"Incorrest input: A dimension of graph is over accessible" }

var mes_onload int = -1

const graph_MaxDimen = 20
var graph [graph_MaxDimen][graph_MaxDimen]int
var graph_dimen = 0

var wg sync.WaitGroup



func calc_ColCount(text string) int {
	col_count := 0
	temp := ""
	rang := -1
	text_len := len(text)
	for k := 0; k < text_len; k++ {
		 if ('0' <= text[k] && text[k] <= '9') || text[k] == '-' || text[k] == ' ' || text[k] == '\n' || text[k] == 13 {
			  if text[k] == '-' && rang != -1 {
					mes_onload = 1
					break
			  }

			  if text[k] == ' ' || text[k] == '\n' {
					if rang != -1 {
						 rang = -1
						 temp = ""
						 col_count++
						 if text[k] == '\n' {
							  break
						 }
					}
			  } else {
					if text[k] != 13 {
						 temp += string(text[k])
						 if text[k] != '-' {
							  rang++
						 }
					}
			  }

			  if k == (text_len-1) && rang != -1 {
						 col_count++
			  }

		 } else {
			  mes_onload = 0
			  break
		 }
	}

	return col_count;
}

func treat_row(k int, i int) {
	for j := 0; j < graph_dimen; j++ {
		 if graph[i][k] == -1 || graph[k][j] == -1 {
			  continue
		 }
		 temp := graph[i][k] + graph[k][j]
		 if graph[i][j] > temp || graph[i][j] == -1 {
			  graph[i][j] = temp
	//		  fmt.Printf("graph[%d][%d] = %d \n ", i, j, temp)
		 }
	}
	defer wg.Done()
}

func calcHandler(w http.ResponseWriter, r *http.Request) {

	text := r.FormValue("intext")

	graph_dimen = calc_ColCount(text)
	if graph_dimen > graph_MaxDimen {
		 mes_onload = 2
		 http.Redirect(w, r, "/graph/", http.StatusFound)
		 return
	}

//read graph
	temp := ""
	rang := -1
	cur_i := 0
	cur_j := 0
	text_len := len(text)
	for k := 0; k < text_len; k++ {
		 if ('0' <= text[k] && text[k] <= '9') || text[k] == '-' || text[k] == ' ' || text[k] == '\n' || text[k] == 13 {
			  if text[k] == '-' && rang != -1 {
					mes_onload = 1
					break
			  }

			  if text[k] == ' ' || text[k] == '\n' {
					if rang != -1 {
						 graph[cur_i][cur_j], _ = strconv.Atoi(temp)
						 rang = -1
						 temp = ""
						 cur_j++
						 if cur_j == graph_dimen {
							  cur_j = 0
							  cur_i++
						 }
					}
			  } else {
					if text[k] != 13 {
						 temp += string(text[k])
						 if text[k] != '-' {
							  rang++
						 }
					}
			  }

			  if k == (text_len-1) && rang != -1 {
					 graph[cur_i][cur_j], _ = strconv.Atoi(temp)
			  }

		 } else {
			  mes_onload = 0
			  break
		 }
	}
	if mes_onload != -1 {
		 http.Redirect(w, r, "/graph/", http.StatusFound)
		 return
	}
//

//calculate
	for k := 0; k < graph_dimen; k++ {
		 for i := 0; i < graph_dimen; i++ {
			  wg.Add(1)
			  go treat_row(k,i)
		 }
	}
	wg.Wait()
//

//output graph
	text = ""

	for i := 0; i < graph_dimen; i++ {
		 for j := 0; j < graph_dimen; j++ {
			  text += strconv.Itoa(graph[i][j])
			  text += " "
		 }
		 text += string('\n')
	}

	pd.Output = []byte(text)
//

	result = true
	http.Redirect(w, r, "/graph/", http.StatusFound)
}


func graphHandler(w http.ResponseWriter, r *http.Request) {
	 head := "<body>"
	 if mes_onload != -1 {
			head = "<head><script> function message() { alert(\" " + message[mes_onload] + " \"); } </script></head><body onload=\"message();\">"

			mes_onload = -1
	 }

    fmt.Fprintf(w,
		  head +
		  "<h1>%s</h1>"+
		  "<h2>%s</h2>"+
        "<form action=\"/calc/\" method=\"POST\">"+
        "<textarea style=\"width: 150px; height: 150px;\" name=\"intext\"></textarea><br>"+
        "<input type=\"submit\" value=\"Calculate\">"+
        "</form>" +
		  "</body>", pd.Title, pd.InvitInput)

	 if result {
			fmt.Fprintf(w, "<h2>%s</h2>"+
								"<textarea style=\"width: 150px; height: 150px;\" name=\"result\">%s</textarea>", pd.InvitOutput, pd.Output)
	 }
}


func main() {
  http.HandleFunc("/graph/", graphHandler)
  http.HandleFunc("/calc/", calcHandler)
  http.ListenAndServe(":8080", nil)
}
