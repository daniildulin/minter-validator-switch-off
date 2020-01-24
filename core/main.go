package core

import (
	"fmt"
	"github.com/MinterTeam/minter-go-sdk/api"
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"github.com/MinterTeam/minter-go-sdk/wallet"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type MinterValidatorSwitchOffService struct {
}

func New() *MinterValidatorSwitchOffService {
	return &MinterValidatorSwitchOffService{}
}

func (s MinterValidatorSwitchOffService) Run() {
	if os.Getenv("NODES_LIST") == "" {
		log.Fatal("Empty nodes list")
	}

	data, err := ioutil.ReadFile("./tx.list")
	if err != nil {
		log.Fatal(err)
	}
	tx, err := transaction.Decode(string(data))
	if err != nil {
		log.Fatal(err)
	}

	var nonce uint64
	clients := s.getNoesList()

	for _, client := range clients {
		nonce, err = client.Nonce(os.Getenv("ADDRESS"))
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}

	if nonce == 0 {
		log.Fatal("Looks like all servers from the list unreachable")
	}

	txNonce := tx.GetTransaction().Nonce

	if txNonce != nonce {
		log.Fatal(`Transaction is not valid! Please run command ./switch -gen_tx -m="mnemonic phrase" before`)
	}

	for {
		fmt.Println("Start watching...")
		status := s.checkMissedBlocks(clients)
		if status {
			s.sendSwitchOffTx(clients, tx)
			fmt.Println("Please update the file with transaction")
			os.Exit(0)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s MinterValidatorSwitchOffService) GenerateTx(mnemonic string) (string, error) {
	var symbol string
	var chainId transaction.ChainID
	if os.Getenv("CHAIN_ID") == "1" {
		chainId = transaction.MainNetChainID
		symbol = "BIP"
	} else {
		chainId = transaction.TestNetChainID
		symbol = "MNT"
	}

	var nonce uint64
	var gp string
	var err error

	clients := s.getNoesList()

	for _, client := range clients {
		nonce, err = client.Nonce(os.Getenv("ADDRESS"))
		if err != nil {
			log.Println(err)
			continue
		}

		gp, err = client.MinGasPrice()
		if err != nil {
			log.Println(err)
		} else {
			break
		}
	}

	gasPrice, err := strconv.ParseInt(gp, 10, 8)
	if err != nil {
		return "", err
	}

	data, err := transaction.NewSetCandidateOffData().SetPubKey(os.Getenv("PUB_KEY"))
	if err != nil {
		return "", err
	}

	tx, err := transaction.NewBuilder(chainId).NewTransaction(data)
	if err != nil {
		return "", err
	}

	seed, err := wallet.Seed(mnemonic)
	if err != nil {
		return "", err
	}

	privateKey, err := wallet.PrivateKeyBySeed(seed)
	if err != nil {
		return "", err
	}

	tx.SetNonce(nonce).SetGasPrice(uint8(gasPrice)).SetGasCoin(symbol)
	signedTx, err := tx.Sign(privateKey)
	if err != nil {
		return "", err
	}

	return signedTx.Encode()
}

func (s MinterValidatorSwitchOffService) CreateFileWithTx(mnemonic string) error {
	fmt.Println("Generate file with Tx")
	hash, err := s.GenerateTx(mnemonic)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./tx.list", []byte(hash), 0644)
	if err != nil {
		return err
	}
	fmt.Println("File has been generated")
	return nil
}

func (s MinterValidatorSwitchOffService) getNoesList() []*api.Api {
	urlList := strings.Split(os.Getenv("NODES_LIST"), " ")
	var clientsList []*api.Api

	for _, url := range urlList {
		clientsList = append(clientsList, api.NewApi(url))
	}
	return clientsList
}

func (s MinterValidatorSwitchOffService) checkMissedBlocks(clients []*api.Api) bool {
	var results []uint64

	maxMissedBlocks, err := strconv.ParseUint(os.Getenv("MISSED_BLOCKS"), 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	for _, client := range clients {
		r, err := client.MissedBlocks(os.Getenv("PUB_KEY"), 0)
		if err != nil {
			fmt.Println(err)
		} else {
			count, err := strconv.ParseUint(r.MissedBlocksCount, 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			results = append(results, count)
		}
	}

	for _, bc := range results {
		if bc >= maxMissedBlocks {
			return true
		}
	}

	return false
}

func (s MinterValidatorSwitchOffService) sendSwitchOffTx(clients []*api.Api, tx transaction.SignedTransaction) {
	for _, client := range clients {
		result, err := client.SendTransaction(tx)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
			break
		}
	}
}
