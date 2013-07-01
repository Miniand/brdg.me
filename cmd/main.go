package main

import (
	"fmt"
	"os"
	
	)

func main() {
	play:=os.Args[1]
	fmt.Println(play)
	if play=="PASS"{
		fmt.Println("Player passed it 'round")
	}
}
