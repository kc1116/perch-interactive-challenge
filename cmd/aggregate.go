package cmd

import (
	"fmt"
	"github.com/kc1116/perch-interactive-challenge/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	threads               int
	host, database, table string
)

var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Will run GCP pubsub event aggregator",
	Long: `The event aggregator is the consumer of events that are published from our IOT devices. 
		A subscription is created and and the specified number of worker threaders are started in background 
		that will continuously process events put on their shared event queue. Events are slightly massaged from protobuf
		serialized objects to plain json objects and then stored in rethinkdb.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if threads < 1 {
			return fmt.Errorf("invalid value for threads %s", threads)
		}

		if threads > 5 {
			return fmt.Errorf("for demonstration purposes aggregator cannot create more than 5 threads")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := aggregateRun()
		if err != nil {
			return err
		}

		return nil
	},
}

func aggregateRun() error {
	registry, err := core.NewDeviceRegistry(projectID, region, registryID, topicID).Init(false)
	if err != nil {
		return err
	}

	aggregator := core.NewEventListener(registry, threads)
	err = aggregator.Start(host, database)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	aggregateCmd.PersistentFlags().IntVarP(&threads, "threads", "W", 1, "Number of worker threaders the event aggregator will create default is 1")
	aggregateCmd.PersistentFlags().StringVarP(&host, "rethinkdb", "H", "127.0.0.1:28015", "Full endpoint to rethinkdb server")
	aggregateCmd.PersistentFlags().StringVarP(&database, "database", "D", "interactions", "Name of rethinkdb database to store events")

	_ = viper.BindPFlag("threads", aggregateCmd.PersistentFlags().Lookup("threads"))
	_ = viper.BindPFlag("rethinkdb", aggregateCmd.PersistentFlags().Lookup("rethinkdb"))
	_ = viper.BindPFlag("database", aggregateCmd.PersistentFlags().Lookup("database"))
}
