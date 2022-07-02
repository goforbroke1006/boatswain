package main

import (
	"fmt"

	"github.com/goforbroke1006/boatswain/internal"
)

func main() {
	chain := internal.BlockChain{}
	chain.Start()

	const t_2022_07_16_54_04_gmt = uint64(1656780844)
	const t_2022_07_16_56_03_gmt = uint64(1656780963)
	const t_2022_07_17_02_55_gmt = uint64(1656781375)

	chain.Generate(t_2022_07_16_54_04_gmt, "Sat Jul 02 2022 16:54:04 GMT+0000")
	chain.Generate(t_2022_07_16_56_03_gmt, "Sat Jul 02 2022 16:56:03 GMT+0000")
	chain.Generate(t_2022_07_17_02_55_gmt, "Sat Jul 02 2022 17:02:55 GMT+0000")

	for _, b := range chain.GetBlocks() {
		fmt.Println(b)
	}
}
