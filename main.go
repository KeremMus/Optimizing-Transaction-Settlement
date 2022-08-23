package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Transaction struct {// Contains parsed transaction data
	from   string
	to     string
	amount float32
}

func check(e error) { //error handling function
	if e != nil {
		panic(e)
	}
}

func stringInSlice(a Transaction, list []Transaction) bool { // Used to determine transactions between two people 
	for _, b := range list {
		if ((b.from == a.from) && (b.to == a.to)) || ((b.from == a.to) && (b.to == a.from)) {
			return true
		}
	}
	return false
}

func remove(slice []Transaction, s int) []Transaction {// Removes an element with given index from a slice and returns the resulting slice 
	return append(slice[:s], slice[s+1:]...)
}

func removefromstring(slice []string, s int) []string {// Removes an element with given index from a slice and returns the resulting slice 
	slice = append(slice[:s], slice[s+1:]...)
	if len(slice) > 0 {
		slice = slice[:len(slice)-1]
	}
	return slice
}

func gettext(transactions []Transaction) []Transaction { // Gets text from the .txt file and parses it into an empty transaction slice  and returns the non-empty transaction slice
	file, ferr := os.Open("/Users/keremcanbkr/Internship/transactions.txt")
	check(ferr)
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseInt(strings.Split(line, ",")[2], 0, 64)
		if err != nil {
			panic(err)
		}
		temptx := Transaction{
			from:   strings.Split(line, ",")[0],
			to:     strings.Split(line, ",")[1],
			amount: float32(value),
		}
		transactions = append(transactions, temptx)
		count++
	}

	return transactions
}

func mergemultiedges(transactions []Transaction, processed []Transaction) []Transaction {
	var currfrom string
	var currto string
	lastindexx := -1
	for i := 0; i < len(transactions); i++ {
		if !stringInSlice(transactions[i], processed) {
			processed = append(processed, transactions[i])
			lastindexx += 1
			currfrom = transactions[i].from
			currto = transactions[i].to
			for j := i + 1; j < len(transactions); j++ {
				if (transactions[j].from == currfrom) && (transactions[j].to == currto) {
					processed[lastindexx].amount = processed[lastindexx].amount + transactions[j].amount
				}
				if (transactions[j].from == currto) && (transactions[j].to == currfrom) {
					processed[lastindexx].amount = processed[lastindexx].amount - transactions[j].amount
					if processed[lastindexx].amount < 0 {
						tempfrom := processed[lastindexx].from
						processed[lastindexx].from = processed[lastindexx].to
						processed[lastindexx].to = tempfrom
						processed[lastindexx].amount = -processed[lastindexx].amount
					}
				}
			}

		}
	}
	for k := 0; k < len(processed); k++ {
		if processed[k].amount == 0 {
			remove(processed, k)
			if len(processed) > 0 {
				processed = processed[:len(processed)-1]
			}
		}
		
	}	
	return processed
}

func getadjlist(transactions []Transaction) map[string][]string {
	adjList := make(map[string][]string)
	for _, tx := range transactions {
		adjList[tx.from] = append(adjList[tx.from], tx.to)
	}
	return adjList
}



func checkcycle(from string, adjList map[string][]string, visited map[string]bool, dfsvisited map[string]bool,cycles map[string][]string) (int, map[string][]string) {
	dfsvisited[from] = true
	visited[from] = true
	if adjList[from] == nil {
		dfsvisited[from] = false
		return 2, cycles
	}
	for _, to := range adjList[from] {
		if !visited[to] {
			cycles[from] = append(cycles[from], to)
			flag := 0
			flag, cycles = checkcycle(to, adjList, visited, dfsvisited,cycles)
			if flag == 1 {
				return 1 , cycles
			}else if flag == 2 {
				cycles[from] = nil
				continue
			}else if dfsvisited[to] {
				return 1 , cycles
			}
		}
		dfsvisited[from] = false
		if len(cycles[from]) == 1 {
			cycles[from] = nil
			return 0, cycles
		}
	}
	return 0, cycles
}

func iscyclic(transactions []Transaction, adjList map[string][]string) (bool, map[string][]string) {
	visited := make(map[string]bool)
	dfsvisited := make(map[string]bool)
	cycles := make(map[string][]string)
	for _, txer := range getuniquetxers(transactions) {
		visited[txer] = false
		dfsvisited[txer] = false
	}
	for i:=0; i<len(transactions); i++ {
		if !visited[transactions[i].from] {
			cycles[transactions[i].from] = nil
			flag := 0
			flag, cycles = checkcycle(transactions[i].from, adjList, visited, dfsvisited,cycles)
			if flag == 1 {
				return true , cycles
			}
		}
	}
	return false, cycles		
}



// func checkcycle(from string, adjList map[string][]string, visited map[string]bool, dfsvisited map[string]bool) bool {
// 	dfsvisited[from] = true
// 	visited[from] = true
// 	for _, to := range adjList[from] {
// 		if !visited[to] {
// 			if checkcycle(to, adjList, visited, dfsvisited) {
// 				return true
// 			}
// 		} else if dfsvisited[to] {
// 			return true
// 		}
// 	}
// 	dfsvisited[from] = false
// 	return false
// }
// func iscyclic(transactions []Transaction, adjList map[string][]string) bool {
// 	visited := make(map[string]bool)
// 	dfsvisited := make(map[string]bool)
// 	for _, txer := range getuniquetxers(transactions) {
// 		visited[txer] = false
// 		dfsvisited[txer] = false
// 	}
// 	for i:=0; i<len(transactions); i++ {
// 		if !visited[transactions[i].from] {
// 			if checkcycle(transactions[i].from, adjList, visited, dfsvisited) {
// 				return true
// 			}
// 		}
// 	}
// 	return false		
// }



func main() {
	var transactions []Transaction
	var processed []Transaction
	transactions = gettext(transactions)
	processed = mergemultiedges(transactions, processed)
	fmt.Println(processed)
	getadjlist(processed)
	//fmt.Println(processed)
	fmt.Println(iscyclic(processed, getadjlist(processed)))
}  

func isatransactioner(frommm string, list []string) bool {
	for _, b := range list {
		if (frommm == b) {
			return true
		}
	}
	return false
}

func isinlist(frommm string, list []string) bool {
	for _, b := range list {
		if frommm == b {
			return true
		}
	}
	return false
}

// func movecurrent(current string, source []string, destination []string) []string {
// 	destination = append(destination, current)
// 	for i:= 0; i < len(source); i++ {
// 		if (source[i] == current) {
// 			source = removefromstring(source, i)
// 		}	
// 	}
// 	return destination
// }

func getuniquetxers(processed []Transaction) []string {
	var txers []string
	for i := 0; i < len(processed); i++ {
		if !isatransactioner(processed[i].from, txers) {
			txers = append(txers, processed[i].from)
		}
		if !isatransactioner(processed[i].to, txers) {
			txers = append(txers, processed[i].to)
		}
	}
	return txers
}	

// func getneighbors(node string, processed []Transaction) []string {
// 	var neighbors []string
// 	for i := 0; i < len(processed); i++ {
// 		if processed[i].from == node {
// 			neighbors = append(neighbors, processed[i].to)
// 		}
// 	}
// 	return neighbors
// }

// func hasloops(transactions []Transaction, processed []Transaction) bool {
// 	var blist []string
// 	var wlist []string
// 	var glist []string
// 	wlist = getuniquetxers(processed)
// 	for i:= 0; i < len(wlist); i++ {
// 		current := wlist[i]
// 		if dfs(current, processed, wlist, blist, glist) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func dfs(current string, processed []Transaction, wlist []string, blist []string, glist []string) bool {
// 	glist = movecurrent(current, wlist, glist)
// 	neighbors := getneighbors(current, processed)
// 	for i := 0; i < len(neighbors); i++ {
// 		if isinlist(neighbors[i], blist) {
// 			continue
// 		}
// 		if isinlist(neighbors[i], glist) {
// 			return true
// 		}
// 		if dfs(neighbors[i], processed, wlist, blist, glist) {
// 			return true
// 		}
// 	}
// 	blist = movecurrent(current, glist, blist)
// 	return false
// }


