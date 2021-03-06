package cmd

import (
	"fmt"
	"github.com/kc1116/perch-interactive-challenge/core"
	"github.com/spf13/cobra"
	"sync"
)

//@TODO: make google certs path configurable
//@TODO: configurable device keys

var (
	sessions, iterations int
)

var sessionCmd = &cobra.Command{
	Use:   "simulator",
	Short: "Start a simulation that attempts to mimick a real perch session with a device",
	Long: `The simulator attempts to do random actions during sessions just like a real person. A session
		can be considered a person coming into contact with our IOT device. Every action this person takes is considered 
		an event and published through GCP IOT Core MQTT Bridge. Each session will not last as long as others this duration 
		is randomized between 5s & 180. A Simulated device acts like it's own independent entity. Thus sessions are independent
		from each other. The frequency in which events occur in a simulated session is a randomized number between 5s and 30s. 
		You can increase the amount of concurrent sessions (default: 2).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if sessions < 1 {
			return fmt.Errorf("invalid value for sessions %s", sessions)
		} else if iterations < 1 {
			return fmt.Errorf("invalid value for iterations %s", iterations)
		}

		if sessions > 5 {
			return fmt.Errorf("for demonstration purposes simulator cannot create more than 5 sessions at a time")
		} else if iterations > 20 {
			return fmt.Errorf("for demonstration purposes simulator cannot do more than 20 iterations %s", iterations)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return StartDeviceSimulation()
	},
}

func StartDeviceSimulation() error {
	for i := 0; i < iterations; i++ {
		wg := &sync.WaitGroup{}
		wg.Add(sessions)
		for i := 0; i < sessions; i++ {
			go StartSimulation(wg)
		}

		wg.Wait()
	}

	return nil
}

func StartSimulation(wg *sync.WaitGroup) {
	defer wg.Done()
	registry, err := core.NewDeviceRegistry(projectID, region, registryID, topicID).Init(true)
	if err != nil {
		logger.Fatalln(err)
	}

	device, err := core.NewDevice(projectID, region, registryID, registry.RegistryName()).Init()
	if err != nil {
		logger.Fatalln(err)
	}

	err = device.ConnectMQTT()
	if err != nil {
		logger.Fatalln(err)
	}

	device.StartSession(wg)

	err = device.CleanUp()
	if err != nil {
		logger.Fatalln(err)
	}
}

func init() {
	sessionCmd.PersistentFlags().IntVarP(&sessions, "sessions", "S", 2, "Number of device simulations to start in parallel")
	sessionCmd.PersistentFlags().IntVarP(&iterations, "iterations", "I", 1, "How many iterations of simulation should device make")
}
