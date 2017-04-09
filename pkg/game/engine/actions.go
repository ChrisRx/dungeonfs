package engine

import (
	"strings"

	sh "github.com/ChrisRx/dungeonfs/pkg/exec/template"
	"github.com/ChrisRx/dungeonfs/pkg/game/fs"
)

func createAction(name string, node fs.Node) fs.Node {
	// onCreate:
	if node.Name() == "door" && name == "wall" {
		node.MetaData().Set("Description", "You found a small switch on the wall and it opened up a path to the east.")
	}
	return node.New(fs.FileNode, name)
}

func lookupAction(name string, node fs.Node) fs.Node {
	switch parseArgs(name) {
	case "look":
		desc := node.MetaData().GetString("Description")
		f := node.New(fs.TempFileNode, "look")
		f.MetaData().Set("Content", sh.Script(sh.Echo(desc)))
		return f
	case "sword":
		var target string
		for _, f := range node.Children().Files() {
			target = f.Name()
			break
		}
		if target == "" {
			target = "unknown"
		}
		commands := []string{
			sh.Echo("you swing your sword mightily at the %s ...", target),
			sh.Command("sleep 1"),
			sh.Echo("It appeared to have no effect."),
		}
		f := node.New(fs.TempFileNode, "sword")
		f.MetaData().Set("Content", sh.Script(commands...))
		return f
	}
	return nil
}

func parseArgs(s string) string {
	args := strings.Fields(s)
	if len(args) > 0 {
		return args[0]
	}
	return s
}
