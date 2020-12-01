package main

import (
	"context"
	"log"
	"os"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
)

// watcher is a simple controller that detects any kubernetes namespace object
// with the ethereum_address label e.g
//
//   apiVersion: v1
//     kind: Namespace
//     metadata:
//       labels:
//       ethereum_address: "0xf71e4b6c2cdfcd83435f357aecd32994f0c69bc3"
//       name: test1
//
// On detection, the controller dynamically populates the namespace's configmap
// with ethereum data. This ethereum data config is now avaialble to be consumed
// by Developer's without implementing any API calls from their apps.
func watcher() {

	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Printf("Client ERROR: %s\n", err.Error())
		os.Exit(1)
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

			if n.Metadata.Labels["ethereum_address"] != "" {
				address := n.Metadata.Labels["ethereum_address"]
				log.Println(*n.Metadata.Name, address)

				//  Get ethereum-data
				ethereumData := GetEthereumData(address)

				// Create configmap
				configMap := &corev1.ConfigMap{
					Metadata: &metav1.ObjectMeta{
						Name:      k8s.String("ethereum-data"),
						Namespace: k8s.String(*n.Metadata.Name),
					},
					Data: ethereumData,
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
