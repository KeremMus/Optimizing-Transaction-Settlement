package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sort"
)

type Transaction struct {// Contains parsed transaction data from the .txt file
	from   string
	to     string
	amount float32
}

type Txer struct { // Contains the name and net balance of a transactioner. "Balance can be a negative number"
	name string
	balance float32
}

func check(e error) { //error handling function
	if e != nil {
		panic(e)
	}
}

func gettext(transactions []Transaction) []Transaction { // Gets text from the .txt file and parses it into an empty transaction slice
	// and returns the non-empty transaction slice
	file, ferr := os.Open("/Users/keremcanbkr/Internship/Optimizing Tx Settlement/csgo_kda_data.txt")
	check(ferr)
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		value, err := strconv.ParseInt(strings.Split(line, " ")[2], 0, 64)
		if err != nil {
			panic(err)
		}
		temptx := Transaction{
			from:   strings.Split(line, " ")[0],
			to:     strings.Split(line, " ")[1],
			amount: float32(value),
		}
		transactions = append(transactions, temptx)
		count++
	}

	return transactions
}

func getadjlist(transactions []Transaction) map[string][]string { // Gets the adjacency list of the transaction graph.
	// The adjacency list is a map of transactioners and their transactioners (accounts that they have transacted to)
	adjList := make(map[string][]string)
	for _, tx := range transactions {
		adjList[tx.from] = append(adjList[tx.from], tx.to)
	}
	return adjList
}

func isatransactioner(frommm string, list []string) bool { // Helper function of the  getuniquetxers function.
	// Checks if a transactioner is already in the list
	for _, b := range list {
		if (frommm == b) {
			return true
		}
	}
	return false
}

func getuniquetxers(processed []Transaction) []string { // Gets the unique transactioners from the transaction slice. 
	// Returns a string slice of unique transactioners
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

func getbalances (transactions []Transaction) (positive_balances []Txer, negative_balances []Txer) { // Gets the net balances of the transactioners. Takes all transactions as input
	// and returns two slices of Txer structs according to their net balance.
	txersnames := getuniquetxers(transactions)
	positive_balances_result := make([]Txer, 0)
	negative_balances_result := make([]Txer, 0)
	for i:=0; i<len(txersnames); i++ {
		temptxer := Txer{
			name: txersnames[i],
			balance: 0,
		}
		for j:=0; j<len(transactions); j++ {
			if transactions[j].from == txersnames[i] {
				temptxer.balance -= transactions[j].amount
			}else if transactions[j].to == txersnames[i] {
				temptxer.balance += transactions[j].amount
			}
		}
		if (temptxer.name != "" && temptxer.balance < 0) {
		negative_balances_result = append(negative_balances_result, temptxer)
		} else if (temptxer.name != "" && temptxer.balance > 0) {
			positive_balances_result = append(positive_balances_result, temptxer)
		}
	}
	return positive_balances_result, negative_balances_result
}

func optimizer(transactions []Transaction, adjList map[string][]string, posi []Txer, nega []Txer) []Transaction { // Main optimization function.
	// It takes all transactions, the adjacency list, the sorted positive balances and the sorted negative balances lists as input.
	// It contains 2 for loops that iterate through the positive and negative balances lists. It distributes negative balances to positive balances.
	var resulttxes []Transaction
	for i:=0; i<len(nega); i++ {
		for j:=0; j<len(posi); j++ {
			if posi[j].balance != 0 && nega[i].balance != 0 {
				tempbalance := nega[i].balance
				nega[i].balance += posi[j].balance
				if (nega[i].balance >= 0) {
					posi[j].balance = nega[i].balance
					nega[i].balance = 0
					resulttxes = append(resulttxes, Transaction{
						from: nega[i].name,
						to: posi[j].name,
						amount: -tempbalance,
					})
					break
				}
				if (nega[i].balance < 0) {
					resulttxes = append(resulttxes, Transaction{
						from: nega[i].name,
						to: posi[j].name,
						amount: posi[j].balance,
					})
					posi[j].balance = 0
					j--
				}
			}
		}
	}
	return resulttxes
}


func main() {
	var transactions []Transaction
	transactions = gettext(transactions)
	for i:=0; i<len(transactions); i++ { // It inverts the transactions to convert kill counts to token counts to transfer.
		tempname := transactions[i].from
		transactions[i].from = transactions[i].to
		transactions[i].to = tempname
	}
	positive_balances, negative_balances := getbalances(transactions)
	sort.SliceStable(positive_balances, func(i, j int) bool { // Sorts the positive balances slice according to the balance values.
		return positive_balances[i].balance < positive_balances[j].balance
	})
	sort.SliceStable(negative_balances, func(i, j int) bool { // Sorts the negative balances slice according to the balance values.
		return negative_balances[i].balance < negative_balances[j].balance
	})
	transactions = optimizer(transactions, getadjlist(transactions), positive_balances, negative_balances)
	for j:=0; j<len(transactions); j++ { // Inverts the transactions back to kill counts.
		fmt.Println(transactions[j].from, transactions[j].to, -transactions[j].amount)
	}
}