package command

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUnmountCommand() *cobra.Command {
	sc := &cobra.Command{
		Use: "unmount [mountpoint]",
		Run: runUnmountCommand,
	}
	// TODO: need to handle overlapping flag names
	//sc.Flags().BoolP("debug", "v", false, "")
	viper.BindPFlags(sc.Flags())
	return sc
}

func runUnmountCommand(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Need mountpoint")
	}
	if _, err := os.Stat("/tmp/dungeonfs.pid"); os.IsNotExist(err) {
		log.Fatalf("DungeonFS does not appear to be running\n")
	}
	data, err := ioutil.ReadFile("/tmp/dungeonfs.pid")
	if err != nil {
		log.Fatal(err)
	}
	pid, err := strconv.Atoi(string(bytes.TrimRight(data, "\n")))
	if err != nil {
		log.Fatal(err)
	}
	syscall.Kill(pid, syscall.SIGINT)
	time.Sleep(1 * time.Second)
	if err := os.Remove(args[0]); err != nil {
		log.Fatal(err)
	}
}
