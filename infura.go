package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// GetEthereumData retreives eth_getBalance for an ethereum_address,
// eth_getBlockByNumber and eth_getTransactionByBlockNumberAndIndex
func GetEthereumData(ethereum_address string) map[string]string {

	start := time.Now()

	client, err := rpc.Dial("https://mainnet.infura.io/v3/70ca5e48e7cb47079e0377062b461ec2")
	if err != nil {
		log.Println("rpc.Dial err", err)
	}

	// eth_getTransactionByBlockNumberAndIndex
	var transaction map[string]string
	err = client.Call(&transaction, "eth_getTransactionByBlockNumberAndIndex", "latest", "0x0")
	if err != nil {
		log.Println("client.Call err", err)
	}
	transactionJson, err := json.Marshal(transaction)
	if err != nil {
		log.Println(err.Error())
	}
	transactionJsonStr := string(transactionJson)

	// eth_getBlockByNumber
	var blockByNumber map[string]interface{}
	err = client.Call(&blockByNumber, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		log.Println("eth_getBlockByNumber client.Call err", err)
	}
	blockByNumberJson, err := json.Marshal(blockByNumber)
	if err != nil {
		log.Println(err.Error())
	}
	blockByNumberJsonStr := string(blockByNumberJson)

	// eth_getBalance
	var balance string
	err = client.Call(&balance, "eth_getBalance", ethereum_address, "latest")
	if err != nil {
		log.Println("eth_getBalance client.Call err", err)
	}
	balanceParsed, err := strconv.ParseUint(hexaNumberToInteger(balance), 16, 64)
	if err != nil {
		fmt.Println(err)
	}

	elapsed := time.Since(start).Seconds()

	var Ethereum = map[string]string{
		"eth_getBalance": strconv.FormatUint(balanceParsed, 10),
		"eth_getTransactionByBlockNumberAndIndex": transactionJsonStr,
		"eth_getBlockByNumber":                    blockByNumberJsonStr,
		"duration":                                fmt.Sprint(elapsed),
	}

	return Ethereum
}

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}
