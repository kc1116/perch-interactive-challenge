package main

import (
	"fmt"
	"github.com/kc1116/perch-interactive-challenge/core"
	"github.com/spf13/viper"
	"log"
)

func main()  {
	viper.AutomaticEnv()
	projectID := "perch-challenge"
	region := "us-central1"
	registryID := "test-registry"
	topicID := "test-registry-topic"

	registry, err := core.NewDeviceRegistry(projectID, region, registryID, topicID).Init()
	if err != nil {
		log.Fatal(err)
	}

	device, err := core.NewDevice(projectID, region, registryID, registry.RegistryName()).Init()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(device)

	err = device.ConnectMQTT()
	if err != nil {
		log.Println(err)
	}

	err = device.CleanUp()
	if err != nil {
		log.Fatal(err)
	}

	//err = registry.CleanUp()
	//if err != nil {
	//	log.Fatal(err)
	//}err = device.CleanUp()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = registry.CleanUp()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
