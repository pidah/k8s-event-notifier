package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

func main() {

	log.Print("Starting ethereum-data-fetcher...")

	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Printf("Client ERROR: %s\n", err.Error())
	}

	var namespace corev1.Namespace

	for {

		watcher, err := client.Watch(context.Background(), "", &namespace)

		if err != nil {
			log.Printf("Watch ERROR: %s\n", err.Error())
		}

		defer watcher.Close()

	WatchLoop:

		for {
			n := new(corev1.Namespace)

			_, err := watcher.Next(n)
			if err != nil {
				log.Printf("Watch ERROR: %s\n", err.Error())
				watcher.Close()
				break WatchLoop
			}

			log.Println(*n.Metadata.Name, n.Metadata.Labels["ethereum_address"])
			if n.Metadata.Labels["ethereum_address"] != "" {
				ethClient, err := ethrpc.Dial("https://mainnet.infura.io/v3/70ca5e48e7cb47079e0377062b461ec2")
				if err != nil {
					fmt.Println("rpc.Dial err", err)
				}

				// eth_getTransactionByBlockNumberAndIndex
				var result map[string]string
				err = ethClient.Call(&result, "eth_getTransactionByBlockNumberAndIndex", "latest", "0x0")

				empData, err := json.Marshal(result)
				if err != nil {
					fmt.Println(err.Error())
				}

				jsonStr := string(empData)
				fmt.Println("eth_getTransactionByBlockNumberAndIndex JSON data is:")
				fmt.Println(jsonStr)

				if err != nil {
					fmt.Println("client.Call err", err)
				}

				// eth_getBlockByNumber
				var result2 map[string]interface{}
				//	var result hexutil.Big
				err = ethClient.Call(&result2, "eth_getBlockByNumber", "latest", false)

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
				//	err = ethClient.Call(&balance, "eth_getBalance", "0xC09e427f7F282172Cc84Ec48eFb30C2F7576D303", "latest")
				err = ethClient.Call(&balance, "eth_getBalance", n.Metadata.Labels["ethereum_address"], "latest")

				if err != nil {
					fmt.Println("client.Call err", err)
				}

				fmt.Printf("%+v", balance)
				output, err := strconv.ParseUint(hexaNumberToInteger(balance), 16, 64)
				if err != nil {
					fmt.Println(err)
				}

				// ethereum-data map
				var m = map[string]string{
					"eth_getTransactionByBlockNumberAndIndex": jsonStr,
					"eth_getBlockByNumber":                    jsonStr2,
					"eth_getBalance":                          strconv.FormatUint(output, 10),
				}

				// Create configmap
				configMap := &corev1.ConfigMap{
					Metadata: &metav1.ObjectMeta{
						Name:      k8s.String("ethereum-data"),
						Namespace: k8s.String(*n.Metadata.Name),
					},
					Data: m,
				}
				if err := client.Create(context.Background(), configMap); err != nil {
					//	if strings.Contains(err.Error(), "409") {
					//		if err := client.Update(context.Background(), configMap); err != nil {
					//			log.Printf("ConfigMap Update ERROR: %s\n", err.Error())
					//		}
					//	}

					//				} else {
					log.Printf("configMap Create ERROR: %s\n", err.Error())
				}

			}
		}
	}
}

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}
