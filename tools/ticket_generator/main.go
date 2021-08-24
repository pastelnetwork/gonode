package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/config"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/generator"
	"github.com/pastelnetwork/gonode/tools/ticket_generator/node"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	exampleConfig string = "./config/example.yaml"
	exampleNFT    string = "./data/test-image1.jpg"
)

type Config struct {
	Miner config.MinerNode  `json:"miner_node" mapstructure:"miner_node"`
	Walet config.WalletNode `json:"walet" mapstructure:"wallet_node"`
}

func main() {
	var configFile string
	var nft string

	flag.StringVar(&configFile, "config-file", exampleConfig, "Path to configuration file")
	flag.StringVar(&nft, "nft", exampleNFT, "Path to NFT file")
	flag.Parse()

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Errorf("read config from %s failed: %w", configFile, err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(errors.Errorf("parse config file %s failed: %w", configFile, err))
	}
	regNFTGenerator := generator.NewNFTRegGenerator(node.NewWallet(config.Walet), node.NewMiner(config.Miner))

	actTicket, err := regNFTGenerator.Gen(context.Background(), exampleNFT)
	if err != nil {
		panic(errors.Errorf("registe NFT %s failed: %w", nft, err))
	}

	js, err := json.MarshalIndent(actTicket, "", "  ")
	if err != nil {
		panic(errors.Errorf("marshal ticket %v failed: %w", actTicket, err))
	}

	fmt.Println(string(js))
}
