package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"bazil.org/fuse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ChrisRx/dungeonfs/pkg/game/assets"
	"github.com/ChrisRx/dungeonfs/pkg/game/engine"
	fs "github.com/ChrisRx/dungeonfs/pkg/game/fs/core"
	"github.com/ChrisRx/dungeonfs/pkg/logging"
)

func NewMountCommand() *cobra.Command {
	sc := &cobra.Command{
		Use: "mount [mountpoint]",
		Run: runMountCommand,
	}
	sc.Flags().BoolP("readonly", "r", false, "")
	sc.Flags().BoolP("debug", "v", false, "")
	sc.Flags().BoolP("daemon", "d", false, "")
	sc.Flags().StringP("assets", "a", "assets/example", "")
	viper.BindPFlags(sc.Flags())
	return sc
}

func runMountCommand(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Need mountpoint")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		os.Mkdir(args[0], 0755)
	}
	if viper.GetBool("daemon") {
		opts := []string{
			"mount",
			args[0],
		}
		if viper.GetBool("debug") {
			opts = append(opts, "-v")
		}
		// TODO: search for executable in CWD and/or PATH
		// TODO: general improvements, platform-specific handling, etc
		command := "bin/dungeonfs"
		cmd := exec.Command(command, opts...)
		cmd.Start()
		pid := fmt.Sprintf("%d\n", cmd.Process.Pid)
		ioutil.WriteFile("/tmp/dungeonfs.pid", []byte(pid), 0755)
		os.Exit(0)
	}
	if viper.GetBool("debug") {
		fs.PkgLogger = &logging.DefaultLogger{}
	}
	fs.GameEngine = engine.NewEngine()
	d, err := assets.LoadAssetsFromFile(viper.GetString("assets"))
	if err != nil {
		log.Fatal(err)
	}
	f, err := fs.NewFS(d)
	if err != nil {
		log.Fatal(err)
	}
	fs.PkgLogger.Printf("READONLY: %t\n", viper.GetBool("readonly"))
	if err = f.MountAndServe(args[0], viper.GetBool("readonly")); err != nil {
		if _, ok := err.(*fuse.MountpointDoesNotExistError); !ok {
			err := fuse.Unmount(args[0])
			if err != nil {
				log.Fatal(err)
			}
			if err = f.MountAndServe(args[0], viper.GetBool("readonly")); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}
