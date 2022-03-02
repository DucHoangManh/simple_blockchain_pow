package main

import (
	"log"
	"simple_blockchain/blockchain"
)

func main() {
	bc, err := blockchain.NewBlockchain()
	if err != nil {
		log.Fatalf("Failed to initialize blockchain %s", err)
	}
	defer func() {
		 err = bc.CloseDbConn()
		 if err != nil {
			 log.Fatalf("Failed to close db connection %s",err)
		 }
	}()

	cli := blockchain.NewCLI(bc)
	cli.Run()
}
