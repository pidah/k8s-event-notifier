package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

//type Ethereum struct {
//	EthGetTransactionByBlockNumberAndIndex []byte
//}

func GetEthereumData() map[string]string {
	client, err := rpc.Dial("https://mainnet.infura.io/v3/70ca5e48e7cb47079e0377062b461ec2")
	if err != nil {
		fmt.Println("rpc.Dial err", err)
	}

	var result map[string]string
	err = client.Call(&result, "eth_getTransactionByBlockNumberAndIndex", "latest", "0x0")

	if err != nil {
		fmt.Println("client.Call err", err)
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		fmt.Println(err.Error())
	}

	// eth_getBlockByNumber
	var result2 map[string]interface{}
	err = client.Call(&result2, "eth_getBlockByNumber", "latest", false)

	if err != nil {
		fmt.Println("client.Call err", err)
	}

	empData2, err := json.Marshal(result2)
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonStr2 := string(empData2)
	fmt.Println("eth_getBlockByNumber JSON data is:")
	fmt.Println(jsonStr2)

	// eth_getBalance
	var balance string
	err = client.Call(&balance, "eth_getBalance", "0xC09e427f7F282172Cc84Ec48eFb30C2F7576D303", "latest")
	// err = ethClient.Call(&balance, "eth_getBalance", n.Metadata.Labels["ethereum_address"]  , "latest")

	if err != nil {
		fmt.Println("client.Call err", err)
	}

	fmt.Printf("%+v", balance)
	output, err := strconv.ParseUint(hexaNumberToInteger(balance), 16, 64)

	if err != nil {
		fmt.Println(err)
	}

	var Ethereum = map[string]string{
		"eth_getBalance": strconv.FormatUint(output, 10),
		"eth_getTransactionByBlockNumberAndIndex": string(jsonData),
		"eth_getBlockByNumber":                    jsonStr2,
	}

	//	log.Printf("public gists for %s: %s\n", "githubUser", string(jsonData))

	fmt.Println(Ethereum["eth_getTransactionByBlockNumberAndIndex"])

	return Ethereum
}

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}
