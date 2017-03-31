package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ChrisRx/dungeonfs/cmd/dungeonfs/command"
)

var RootCmd = &cobra.Command{
	Use:   "dungeonfs",
	Short: "Action-packed dungeon crawling file system",
}

func main() {
	RootCmd.AddCommand(
		command.NewMountCommand(),
		command.NewUnmountCommand(),
	)
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
