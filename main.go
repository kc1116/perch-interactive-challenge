package main

import "github.com/kc1116/perch-interactive-challenge/cmd"

func main() {
	cmd.Execute()
	//viper.AutomaticEnv()
	//projectID := "perch-challenge"
	//region := "us-central1"
	//registryID := "test-registry"
	//topicID := "test-registry-topic"
	//
	//registry, err := core.NewDeviceRegistry(projectID, region, registryID, topicID).Init(false)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//aggregator := core.NewEventListener(registry, 1)
	//err = aggregator.Start()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//device, err := core.NewDevice(projectID, region, registryID, registry.RegistryName()).Init()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(device)
	//
	//err = device.ConnectMQTT()
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//device.StartSession()
	//
	//err = device.CleanUp()
	//if err != nil {
	//	log.Fatal(err)
	//}

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
