package main

import (
	"context"
	"log"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
)

func watcher() {

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
				// ethereum-data map
				m := GetEthereumData()
				// Create configmap
				configMap := &corev1.ConfigMap{
					Metadata: &metav1.ObjectMeta{
						Name:      k8s.String("ethereum-data"),
						Namespace: k8s.String(*n.Metadata.Name),
					},
					Data: m,
				}
				if err := client.Create(context.Background(), configMap); err != nil {
					//  if strings.Contains(err.Error(), "409") {
					//      if err := client.Update(context.Background(), configMap); err != nil {
					//          log.Printf("ConfigMap Update ERROR: %s\n", err.Error())
					//      }
					//  }

					//              } else {
					log.Printf("configMap Create ERROR: %s\n", err.Error())
				}

			}
		}
	}
}
