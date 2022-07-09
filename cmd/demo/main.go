package main

import (
	"fmt"
	"time"

	"github.com/goforbroke1006/boatswain/internal"
)

func main() {
	chain := internal.BlockChain{}
	chain.Start()

	var (
		t1 = time.Date(2022, time.July, 16, 54, 14, 0, 0, time.UTC)
		t2 = time.Date(2022, time.July, 16, 56, 13, 0, 0, time.UTC)
		t3 = time.Date(2022, time.July, 17, 12, 55, 0, 0, time.UTC)
	)
	chain.Generate(t1.Unix(), "Sat Jul 02 2022 16:54:04 GMT+0000")
	chain.Generate(t2.Unix(), "Sat Jul 02 2022 16:56:03 GMT+0000")
	chain.Generate(t3.Unix(), "Sat Jul 02 2022 17:02:55 GMT+0000")

	for _, b := range chain.GetBlocks() {
		fmt.Println(b)
	}
}
