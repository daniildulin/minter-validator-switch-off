package core

import (
	"fmt"
	"github.com/MinterTeam/minter-go-sdk/v2/api/grpc_client"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"github.com/sirupsen/logrus"
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
	var chainId transaction.ChainID
	var nonce uint64
	var gp *api_pb.MinGasPriceResponse
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

	tx.SetNonce(nonce).SetGasPrice(uint8(gp.MinGasPrice)).SetGasCoin(0)
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

func (s MinterValidatorSwitchOffService) getNoesList() []*grpc_client.Client {
	urlList := strings.Split(os.Getenv("NODES_LIST"), " ")
	var clientsList []*grpc_client.Client

	for _, url := range urlList {
		nodeApi, err := grpc_client.New(url)
		if err != nil {
			logrus.Fatal(err)
		}

		clientsList = append(clientsList, nodeApi)
	}
	return clientsList
}

func (s MinterValidatorSwitchOffService) checkMissedBlocks(clients []*grpc_client.Client) bool {
	var results []int64

	maxMissedBlocks, err := strconv.ParseInt(os.Getenv("MISSED_BLOCKS"), 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	for _, client := range clients {
		r, err := client.MissedBlocks(os.Getenv("PUB_KEY"))
		if err != nil {
			fmt.Println(err)
		} else {
			results = append(results, r.MissedBlocksCount)
		}
	}

	for _, bc := range results {
		if bc >= maxMissedBlocks {
			return true
		}
	}

	return false
}

func (s MinterValidatorSwitchOffService) sendSwitchOffTx(clients []*grpc_client.Client, tx transaction.Signed) {
	for _, client := range clients {
		txString, err := tx.Encode()
		result, err := client.SendTransaction(txString)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
			break
		}
	}
}
