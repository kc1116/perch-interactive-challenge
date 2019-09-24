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

var rootCmd = &cobra.Command{
	Use:   "perch-iot-pubsub",
	Short: "CLI tool for running perch iot pubsub aggregator, or simulated device interaction session",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectID, "projectID", "p", "perch-challenge", "Google cloud project ID")
	rootCmd.PersistentFlags().StringVarP(&registryID, "registryID", "r", "test-registry", "Google cloud IOT core device registry ID")
	rootCmd.PersistentFlags().StringVarP(&topicID, "topicID", "t", "test-registry-topic", "Google cloud Pubsub topic ID")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "R", "us-central1", "Google cloud region")

	_ = viper.BindPFlag("projectID", rootCmd.PersistentFlags().Lookup("projectID"))
	_ = viper.BindPFlag("registryID", rootCmd.PersistentFlags().Lookup("registryID"))
	_ = viper.BindPFlag("topicID", rootCmd.PersistentFlags().Lookup("topicID"))
	_ = viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))

	rootCmd.AddCommand(aggregateCmd, sessionCmd, websocketCmd)
}
