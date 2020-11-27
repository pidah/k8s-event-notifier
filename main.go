package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
				eth_client, err := ethclient.Dial("https://mainnet.infura.io/v3/70ca5e48e7cb47079e0377062b461ec2")
				if err != nil {
					log.Fatal(err)
				}

				//	account := common.HexToAddress("0xC09e427f7F282172Cc84Ec48eFb30C2F7576D303")
				account := common.HexToAddress(n.Metadata.Labels["ethereum_address"])
				balance, err := eth_client.BalanceAt(context.Background(), account, nil)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(balance)
				configMap := &corev1.ConfigMap{
					Metadata: &metav1.ObjectMeta{
						Name:      k8s.String("ethereum-data"),
						Namespace: k8s.String(*n.Metadata.Name),
					},
					Data: map[string]string{"ethereum_address": n.Metadata.Labels["ethereum_address"], "account_balance": "66666666"},
				}
				if err := client.Create(context.Background(), configMap); err != nil {
					if strings.Contains(err.Error(), "409") {
						if err := client.Update(context.Background(), configMap); err != nil {
							log.Printf("ConfigMap Update ERROR: %s\n", err.Error())
						}
					}

				} else {
					log.Printf("ConfigMap Create ERROR: %s\n", err.Error())
				}
			}
		}
	}
}
