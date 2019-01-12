package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("sdf")
	for i := 12; i > 7; i-- {
		time.Sleep(time.Second)
		fmt.Printf("\033[H\033[2J%v", i)
	}
	fmt.Println("")

}
