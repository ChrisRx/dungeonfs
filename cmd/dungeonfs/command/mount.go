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
	"github.com/kardianos/osext"

	"github.com/ChrisRx/dungeonfs/pkg/eval"
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
	sc.Flags().CountP("debug", "v", "")
	sc.Flags().BoolP("daemon", "d", false, "")
	sc.Flags().StringP("assets", "a", "examples/simplelevel", "")
	viper.BindPFlags(sc.Flags())
	return sc
}

func runMountCommand(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Need mountpoint")
	}
	if _, err := os.Stat(viper.GetString("assets")); os.IsNotExist(err) {
		log.Fatal("Assets not found")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		os.Mkdir(args[0], 0755)
	}
	if viper.GetBool("daemon") {
		opts := []string{
			"mount",
			args[0],
			"-a",
			viper.GetString("assets"),
		}
		if viper.GetBool("debug") {
			opts = append(opts, "-v")
		}
		exe, err := osext.Executable()
		if err != nil {
			log.Fatal("Unable to determine exectuable path, can't start daemon")
		}
		cmd := exec.Command(exe, opts...)
		cmd.Start()
		pid := fmt.Sprintf("%d\n", cmd.Process.Pid)
		ioutil.WriteFile("/tmp/dungeonfs.pid", []byte(pid), 0755)
		os.Exit(0)
	}
	if viper.GetInt("debug") > 0 {
		fs.PkgLogger = &logging.DefaultLogger{}
		engine.PkgLogger = &logging.DefaultLogger{}
		assets.PkgLogger = &logging.DefaultLogger{}
		if viper.GetInt("debug") > 1 {
			eval.PkgLogger = &logging.DefaultLogger{}
		}
	}
	r := assets.New()
	d, err := r.LoadDir(viper.GetString("assets"))
	if err != nil {
		log.Fatal(err)
	}

	fs.GameEngine = engine.NewEngine(r)

	f, err := fs.NewFS(d)
	if err != nil {
		log.Fatal(err)
	}
	if err = f.MountAndServe(args[0], false); err != nil {
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
