package cmd

import (
	"fmt"
	"github.com/kc1116/perch-interactive-challenge/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var websocketCmd = &cobra.Command{
	Use:   "websocket",
	Short: "Will run websocket server to stream events to clients",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := core.NewStore(host, database)
		if err != nil {
			return fmt.Errorf("unable to start websocket err when connecting to event store %s", err)
		}

		err = store.Init()
		if err != nil {
			return fmt.Errorf("error initializing store %s", err)
		}

		store.StartWSProxy()
		return nil
	},
}

func init()  {
	websocketCmd.PersistentFlags().StringVarP(&host, "rethinkdb", "H", "127.0.0.1:28015", "Full endpoint to rethinkdb server")
	websocketCmd.PersistentFlags().StringVarP(&database, "database", "D", "interactions", "Name of rethinkdb database to store events")
	_ = viper.BindPFlag("rethinkdb", websocketCmd.PersistentFlags().Lookup("rethinkdb"))
	_ = viper.BindPFlag("database", websocketCmd.PersistentFlags().Lookup("database"))
}