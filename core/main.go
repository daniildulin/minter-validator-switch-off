package core

import (
	"fmt"
	"github.com/MinterTeam/minter-go-sdk/api"
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"github.com/MinterTeam/minter-go-sdk/wallet"
	"io/ioutil"
	"os"
	"strconv"
)

type MinterValidatorSwitchOffService struct {
}

func New() *MinterValidatorSwitchOffService {
	return &MinterValidatorSwitchOffService{}
}

func (s MinterValidatorSwitchOffService) Run() {
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

	minterClient := api.NewApi(os.Getenv("NODE_URL"))

	data, err := transaction.NewSetCandidateOffData().SetPubKey(os.Getenv("PUB_KEY"))
	if err != nil {
		return "", err
	}

	tx, err := transaction.NewBuilder(chainId).NewTransaction(data)
	if err != nil {
		return "", err
	}

	nonce, err := minterClient.Nonce(os.Getenv("ADDRESS"))
	if err != nil {
		return "", err
	}

	gp, err := minterClient.MinGasPrice()
	gasPrice, err := strconv.ParseInt(gp, 10, 8)
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
