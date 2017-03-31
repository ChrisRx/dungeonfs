package assets

import (
	"io/ioutil"
	filepath "path"

	"gopkg.in/yaml.v2"

	fs "github.com/ChrisRx/dungeonfs/pkg/game/fs/core"
)

var (
	Troll = `
                       _-------------------_
     _            _   /                     \
    | \  __--__  / |  |  Roar! I'm a troll! | 
    |  \/      \/  |  \                     /
    \__  \    / __ /   -_  ________________-
     /   O    -   \     / /
    |      oo      |    /
    \ -- ______ -- /  
     (_  |\||/|  _)
       --______--
`
	Key = []byte(`-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMjSi/S67mqFDGP0
+MRBkidG6zenktXMPMcwXCN1WQjHAo4Gcp9GqjH+IVTgBSO7tdz9VxxKAJgSFXqv
kWj95QjoT6OTjVk0UXH4t9PV7OU2jn/xnA8pqhg4sxVzbVq9LSB55CsqNpX5IpUu
9cUUgiq9rzDBmh4bg1X9fUHufAxbAgMBAAECgYAfARKGeA2y+FOPYxS9B/qOgc5y
yzZKN7vybK7s8oMKbd8hGjG8EWbZTQjMV8GzYJmVQq+eOHabA7+5Lz3d3cTsMu5I
DSCKGers7AHHQAYQv0P14CxBp5PFY5qewMYB/FITcD/Z25YXlg3ZjwCB3XrQwLnN
QB83C1r6lcFi2fdioQJBAPP/Wy39L+Vda7thcRPcQpb4FmpwU+v+Rzpxov5KRq9G
38u3BVd0GdSwQWLVDMmvSLy70doGgxJ5p54YPScxCfkCQQDSs31+ftL34+lSXtrx
5utBZc34q3UcCOa3twoHTzxGeM4BiYvAcVa+PjdWRaNXz71UBs2GQWWihApZoGk9
2X3zAkEAllFwI/ICauTV9Re/6UNeBtIKRUK0gQQjb58Ikm7CA0O/pio38TvGmiCH
99JXUX1aa2OukgpG/7/RAvXd3uI4SQJAdHJ2nP6CojX3sWpzHtY8lrwpBZHc+02A
FXC3vipwaZJCaF8YOZdqFWJVOvzptZI+VL4dwGFMRnErNzWMdH5LOQJABytqQ8eF
qM/JedHdd4l6ADmUA3A3JxF6eQwYUEd7V51f1dKHB7El2nCLPD1lLgWNewYLg5a6
eVH4jECLO2OH8Q==
-----END PRIVATE KEY----- `)
)

type DirectoryAsset struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Adjacent    []string `yaml:"adjacent"`
	Contains    []string `yaml:"contains"`
}

type FileAsset struct {
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

type Assets struct {
	Files       map[string]FileAsset
	Directories map[string]DirectoryAsset
}

func NewAssets() *Assets {
	a := &Assets{
		Files:       make(map[string]FileAsset),
		Directories: make(map[string]DirectoryAsset),
	}
	return a
}

func LoadAssetsFromFile(folder string) (*fs.Directory, error) {
	a := NewAssets()
	return a.LoadAssets(folder)
}

func (a *Assets) LoadAssets(folder string) (*fs.Directory, error) {
	data, err := ioutil.ReadFile(filepath.Join(folder, "directories.yaml"))
	if err != nil {
		return nil, err
	}
	var dirs []DirectoryAsset
	err = yaml.Unmarshal([]byte(data), &dirs)
	if err != nil {
		return nil, err
	}
	for _, d := range dirs {
		a.Directories[d.Name] = d
	}
	data, err = ioutil.ReadFile(filepath.Join(folder, "files.yaml"))
	if err != nil {
		return nil, err
	}
	var files []FileAsset
	err = yaml.Unmarshal([]byte(data), &files)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		a.Files[f.Name] = f
	}
	root := fs.NewDirectory("Root", "")
	a.buildLevel(root)
	return root, nil
}

func (a *Assets) buildLevel(d *fs.Directory) {
	if val, ok := a.Directories[d.Name()]; ok {
		for _, name := range val.Adjacent {
			if v, ok := a.Directories[name]; ok {
				newDir := d.NewDirectory(v.Name)
				newDir.MetaData().Set("Description", v.Description)
				a.buildLevel(newDir)
			}
		}
		for _, item := range val.Contains {
			if v, ok := a.Files[item]; ok {
				f := d.NewFile(v.Name)
				f.MetaData().Set("Content", []byte(v.Content))
			}
		}
	}
}
