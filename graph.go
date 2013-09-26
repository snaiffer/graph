package main

import (
	"os"
	"io"
	"html/template"
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

type ViewGraph struct {
	 Nodes []byte
	 Edges []byte
}

var ViewG ViewGraph


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
					}
					if text[k] == '\n' {
						 break
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
		 }
	}
	defer wg.Done()
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	result = false
	text := r.FormValue("intext")
	pd.Input = []byte(text)

	graph_dimen = calc_ColCount(text)
	if graph_dimen > graph_MaxDimen {
		 mes_onload = 2
		 http.Redirect(w, r, "/graph/", http.StatusFound)
		 return
	}
	if graph_dimen == 0 {
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

//Generate paths for a graphic graph
	text = ""
	for i := 0; i < graph_dimen; i++ {
		 text += "{ data: { id: '" + strconv.Itoa(i) + "', name: '" + strconv.Itoa(i) + "' } },"
	}
	text = text[:(len(text)-1)]
	ViewG.Nodes = []byte(text)

	text = ""
	for i := 0; i < graph_dimen; i++ {
		 for j := 0; j < graph_dimen; j++ {
			  if graph[i][j] != 0 {
					text += "{ data: { source: '" + strconv.Itoa(i) + "', target: '" + strconv.Itoa(j) + "' } },"
			  }

		 }
	}
	text = text[:(len(text)-1)]
	ViewG.Edges = []byte(text)
//

	result = true
	http.Redirect(w, r, "/graph/", http.StatusFound)
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
		 // http.ServeFile(w, r, "5.html")

	 if mes_onload != -1 {
		  head := "<script> function message() { alert(\" " + message[mes_onload] + " \"); } </script></head><body onload=\"message();\">"

		  fmt.Fprintf(w,
				head +
				"<h1>%s</h1>"+
				"<h2>%s</h2>"+
				"<form action=\"/calc/\" method=\"POST\">"+
				"<textarea style=\"width: 150px; height: 150px;\" name=\"intext\">%s</textarea><br>"+
				"<input type=\"submit\" value=\"Calculate\">"+
				"</form>" +
				"</body>", pd.Title, pd.InvitInput, pd.Input)

			mes_onload = -1
	 } else {

		  if result {

			  bottom := "<div style=\"height:700px;  width: 400px\" name=\"cy\" id=\"cy\"></div>" +
				 "" +
				 "<script>" +
				 "$('#cy').cytoscape({" +
				 "  style: cytoscape.stylesheet()" +
				 "    .selector('node')" +
				 "      .css({" +
				 "        'content': 'data(name)'," +
				 "        'text-valign': 'center'," +
				 "        'color': 'white'," +
				 "        'text-outline-width': 2," +
				 "        'text-outline-color': '#888'," +
				 "      })" +
				 "    .selector('edge')" +
				 "      .css({" +
				 "        'target-arrow-shape': 'triangle'" +
				 "      })" +
				 "    .selector(':selected')" +
				 "      .css({" +
				 "        'background-color': 'black'," +
				 "        'line-color': 'black'," +
				 "        'target-arrow-color': 'black'," +
				 "        'source-arrow-color': 'black'" +
				 "      })" +
				 "    .selector('.faded')" +
				 "      .css({" +
				 "        'opacity': 0.25," +
				 "        'text-opacity': 0" +
				 "      })," +
				 "" +
				 "  elements: {" +
				 "    nodes: ["

				 /*
				 ViewG.Nodes = []byte("{ data: { id: 'j', name: 'Jerry' } }," +
						"{ data: { id: 'e', name: 'Elaine' } }," +
						"{ data: { id: 'k', name: 'Kramer' } }," +
						"{ data: { id: 'g', name: 'George' } }," +
						"{ data: { id: 'a', name: 'Alex' } } ")
						*/
				 bottom += string(ViewG.Nodes)

				 bottom += "]," +
				 "    edges: ["

				 /*
				 ViewG.Edges = []byte("{ data: { source: 'j', target: 'e' } }," +
					 "{ data: { source: 'j', target: 'k' } }," +
					 "{ data: { source: 'j', target: 'g' } }," +
					 "{ data: { source: 'e', target: 'j' } }," +
					 "{ data: { source: 'e', target: 'k' } }," +
					 "{ data: { source: 'k', target: 'j' } }," +
					 "{ data: { source: 'k', target: 'e' } }," +
					 "{ data: { source: 'k', target: 'g' } }," +
					 "{ data: { source: 'g', target: 'j' } }," +
					 "{ data: { source: 'g', target: 'a' } }" )
					 */
				 bottom += string(ViewG.Edges)

				 bottom += "]" +
				 "  }," +
				 "" +
				 "  ready: function(){" +
				 "    window.cy = this;" +
				 "" +
				 "    cy.elements().unselectify();" +
				 "" +
				 "    cy.on('tap', 'node', function(e){" +
				 "      var node = e.cyTarget;" +
				 "      var neighborhood = node.neighborhood().add(node);" +
				 "" +
				 "      cy.elements().addClass('faded');" +
				 "      neighborhood.removeClass('faded');" +
				 "    });" +
				 "" +
				 "    cy.on('tap', function(e){" +
				 "      if( e.cyTarget === cy ){" +
				 "        cy.elements().removeClass('faded');" +
				 "      }" +
				 "    });" +
				 "  }" +
				 "});" +
				 "</script>" +
				 "</body>" +
				 "</html>"


	  src, _ := os.Open("orig.html")
	  dest, _ := os.Create("dest.html")
	  io.Copy(dest, src)
	  dest.WriteString("<body>" +
		  "<h1>" + string(pd.Title) + "</h1>"+
		  "<h2>" + string(pd.InvitInput) + "</h2>"+
		  "<form action=\"/calc/\" method=\"POST\">"+
		  "<textarea style=\"width: 150px; height: 150px;\" name=\"intext\">" + string(pd.Input) + "</textarea><br>"+
		  "<input type=\"submit\" value=\"Calculate\">"+
		  "</form>" +
		  "<h2>" + string(pd.InvitOutput) + "</h2>"+
		  "<textarea style=\"width: 150px; height: 150px;\" name=\"result\">" +
		  string(pd.Output) +
		  "</textarea>" +
		  bottom)

	 t, _ := template.ParseFiles("dest.html")
	 t.Execute(w, nil)
		  } else {
				fmt.Fprintf(w,
					 "<body>" +
					 "<h1>%s</h1>"+
					 "<h2>%s</h2>"+
					 "<form action=\"/calc/\" method=\"POST\">"+
					 "<textarea style=\"width: 150px; height: 150px;\" name=\"intext\"></textarea><br>"+
					 "<input type=\"submit\" value=\"Calculate\">"+
					 "</form>" +
					 "</body>", pd.Title, pd.InvitInput)
		  }

	 }

}


func main() {
  http.HandleFunc("/graph/", graphHandler)
  http.HandleFunc("/calc/", calcHandler)
  http.ListenAndServe(":8080", nil)
}
