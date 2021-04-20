package core

import (
	"fmt"
	"github.com/MinterTeam/minter-go-sdk/v2/api/grpc_client"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
	"github.com/MinterTeam/minter-go-sdk/v2/wallet"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"github.com/daniildulin/minter-validator-switch-off/bot"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	Vs string
)

type MinterValidatorSwitchOffService struct {
	nodeClients []*grpc_client.Client
	log         *logrus.Entry
	tx          transaction.Signed
	TgBot       *bot.TgBot
}

func New() *MinterValidatorSwitchOffService {
	var err error
	var nonce uint64
	var tx transaction.Signed

	//Init Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	log := logger.WithFields(logrus.Fields{
		"version": "1.2.0",
		"app":     "Minter Validator Protector",
	})

	if os.Getenv("NODES_LIST") == "" {
		log.Fatal("Empty nodes list")
	}

	clients := getNoesList()

	if Vs == "" {
		tx, err = getTxFromFile()
		if err != nil {
			log.Fatal(err)
		}

		for _, client := range clients {
			nonce, err = client.Nonce(os.Getenv("ADDRESS"))
			if err != nil {
				log.Error(err)
			} else {
				break
			}
		}
		if nonce == 0 {
			log.Fatal("Looks like all servers from the list is unreachable")
		}
		if tx.GetTransaction().Nonce != nonce {
			log.Fatal(`Transaction is not valid! Please run command ./switch -gen_tx -m="mnemonic phrase" before`)
		}
	} else {
		log.Info(Vs)
	}

	var tgBot *bot.TgBot
	if os.Getenv("TG_TOKEN") != "" && os.Getenv("TG_CHANNEL_ID") != "" {
		tgBot = bot.New()
	}

	return &MinterValidatorSwitchOffService{
		nodeClients: clients,
		log:         log,
		tx:          tx,
		TgBot:       tgBot,
	}
}

func (s MinterValidatorSwitchOffService) Run() {
	for {
		fmt.Println("Start watching...")
		if s.checkMissedBlocks() && s.isValidatorEnabled() {
			s.sendSwitchOffTx()
			msg := fmt.Sprintf("The validator has been stopped at %s", time.Now().Format("2006.01.02-15:04:05"))
			s.log.Warn(msg)
			if s.TgBot != nil {
				s.TgBot.SendMsg(msg)
			}
			if Vs == "" {
				msg = "Please update the file with transaction"
				s.log.Warn(msg)
				if s.TgBot != nil {
					s.TgBot.SendMsg(msg)
				}
				os.Exit(0)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (s MinterValidatorSwitchOffService) GenerateTx(mnemonic string) (transaction.Signed, error) {
	var gp *api_pb.MinGasPriceResponse
	var chainId transaction.ChainID
	var nonce uint64
	var err error

	for _, client := range s.nodeClients {
		nonce, err = client.Nonce(os.Getenv("ADDRESS"))
		if err != nil {
			s.log.Println(err)
			continue
		}
		gp, err = client.MinGasPrice()
		if err != nil {
			s.log.Println(err)
		} else {
			break
		}
	}

	data, err := transaction.NewSetCandidateOffData().SetPubKey(os.Getenv("PUB_KEY"))
	if err != nil {
		return nil, err
	}

	tx, err := transaction.NewBuilder(chainId).NewTransaction(data)
	if err != nil {
		return nil, err
	}

	seed, err := wallet.Seed(mnemonic)
	if err != nil {
		return nil, err
	}

	privateKey, err := wallet.PrivateKeyBySeed(seed)
	if err != nil {
		return nil, err
	}

	tx.SetNonce(nonce).SetGasPrice(uint8(gp.MinGasPrice)).SetGasCoin(0)
	signedTx, err := tx.Sign(privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (s MinterValidatorSwitchOffService) CreateFileWithTx(mnemonic string) error {
	fmt.Println("Generate file with Tx")
	tx, err := s.GenerateTx(mnemonic)
	if err != nil {
		return err
	}
	hash, err := tx.Encode()
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

func (s MinterValidatorSwitchOffService) isValidatorEnabled() bool {
	for _, client := range s.nodeClients {
		r, err := client.Candidate(os.Getenv("PUB_KEY"))
		if err != nil {
			fmt.Println(err)
		} else {
			if r.Status == 2 {
				return true
			}
		}
	}
	return false
}
func (s MinterValidatorSwitchOffService) checkMissedBlocks() bool {
	var results []int64

	maxMissedBlocks, err := strconv.ParseInt(os.Getenv("MISSED_BLOCKS"), 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	for _, client := range s.nodeClients {
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

func (s MinterValidatorSwitchOffService) sendSwitchOffTx() {
	var err error
	var tx transaction.Signed

	if Vs == "" {
		tx = s.tx
	} else {
		tx, err = s.GenerateTx(Vs)
	}

	if err != nil {
		s.log.Error(err)
		return
	}

	for _, client := range s.nodeClients {
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

func getTxFromFile() (transaction.Signed, error) {
	data, err := ioutil.ReadFile("./tx.list")
	if err != nil {
		return nil, err
	}
	tx, err := transaction.Decode(string(data))
	if err != nil {
		return nil, err
	}
	return tx, err
}

func getNoesList() []*grpc_client.Client {
	urlList := strings.Split(os.Getenv("NODES_LIST"), " ")
	var clientsList []*grpc_client.Client
	for _, url := range urlList {
		nodeApi, err := grpc_client.New(url)
		if err != nil {
			panic(err)
		}
		clientsList = append(clientsList, nodeApi)
	}
	return clientsList
}
