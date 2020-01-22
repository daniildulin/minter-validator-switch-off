package main

import (
	"flag"
	"github.com/daniildulin/minter-validator-switch-off/core"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	var generateTxMode = flag.Bool(`gen_tx`, false, `Generate file with switch off transaction`)
	var mnemonic = flag.String("m", "", "Mnemonic phrase")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	switcher := core.New()

	if *generateTxMode {
		err := switcher.CreateFileWithTx(*mnemonic)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	switcher.Run()
	os.Exit(0)
}
