package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var round int
	fmt.Scanf("%d", &round)
	var ignore int
	for i := 0; i < round - 1; i++ {
		fmt.Scanf("%d%d%d%d", &ignore, &ignore, &ignore, &ignore)
	}

	fmt.Println(round % 5 + 1, round % 5 + 1 + rand.Int() % 5 + 1)
}
