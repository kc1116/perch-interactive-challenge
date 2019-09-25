package cmd

import (
	"github.com/kc1116/perch-interactive-challenge/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	projectID, region, registryID, topicID, googleCloudAuth string
	logger                                                  = core.Logger()
)

var RootCmd = &cobra.Command{
	Use:   "perch-iot-pubsub",
	Short: "CLI tool for running perch iot pubsub aggregator, or simulated device interaction session",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&projectID, "projectID", "p", "perch-challenge", "Google cloud project ID")
	RootCmd.PersistentFlags().StringVarP(&registryID, "registryID", "r", "test-registry", "Google cloud IOT core device registry ID")
	RootCmd.PersistentFlags().StringVarP(&topicID, "topicID", "t", "test-registry-topic", "Google cloud Pubsub topic ID")
	RootCmd.PersistentFlags().StringVarP(&region, "region", "R", "us-central1", "Google cloud region")

	_ = viper.BindPFlag("projectID", RootCmd.PersistentFlags().Lookup("projectID"))
	_ = viper.BindPFlag("registryID", RootCmd.PersistentFlags().Lookup("registryID"))
	_ = viper.BindPFlag("topicID", RootCmd.PersistentFlags().Lookup("topicID"))
	_ = viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))

	RootCmd.AddCommand(aggregateCmd, sessionCmd, websocketCmd)
}
