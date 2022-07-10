package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/goforbroke1006/boatswain/internal/blockchain"
	"github.com/goforbroke1006/boatswain/internal/storage"
)

func main() {
	db, err := sql.Open("sqlite3", "./demo-blocks.db")
	if err != nil {
		panic(err)
	}
	schemaQuery, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(schemaQuery)); err != nil {
		panic(err)
	}
	blockStorage := storage.NewBlockStorage(db)

	chain := blockchain.NewBlockChain(blockStorage)
	if chain.Empty() {
		fmt.Println("fill empty blockchain")
		var (
			t1 = time.Date(2022, time.July, 16, 54, 14, 0, 0, time.UTC)
			t2 = time.Date(2022, time.July, 16, 56, 13, 0, 0, time.UTC)
			t3 = time.Date(2022, time.July, 17, 12, 55, 0, 0, time.UTC)
		)
		_ = chain.Generate(t1.Unix(), "Sat Jul 02 2022 16:54:04 GMT+0000")
		_ = chain.Generate(t2.Unix(), "Sat Jul 02 2022 16:56:03 GMT+0000")
		_ = chain.Generate(t3.Unix(), "Sat Jul 02 2022 17:02:55 GMT+0000")
	}

	for _, b := range chain.GetBlocks() {
		fmt.Println(b)
	}
}
