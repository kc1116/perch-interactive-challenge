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
	Long: `Our Store struct is responsible for all interactions to rethinkdb, we extend this functionality 
		and add a simple websocket server that can be started with this command. The server exposes a websocket endpoint 
		on port :8000 and acts like a proxy between our client web app and rethinkdb it'self because rethinkdb does not except 
		websocket connections. We make use of rethinkdbs change sets which allows us to watch all updates on our events table and 
		stream them to the ui via websocket in real time.`,
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