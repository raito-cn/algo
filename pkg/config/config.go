package config

import (
	"algo/internal/util"
	"github.com/creasty/defaults"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type Config struct {
	Dir Dir `toml:"dir"`
}

type Dir struct {
	CodeDir     string `toml:"code_dir" default:"~/algo/code"`
	NotesDir    string `toml:"notes_dir" default:"~/algo/notes"`
	Datasource  string `toml:"datasource" default:"~/algo/db"`
	MarkdownDir string `toml:"markdown_dir" default:"~/algo/markdown"`
}

func (d *Dir) ExpandHome(home string) {
	v := reflect.ValueOf(d).Elem() // 获取结构体指针的值
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			val := field.String()
			if strings.HasPrefix(val, "~/") {
				field.SetString(filepath.Join(home, val[2:]))
			}
		}
	}
}

var cnf *Config
var once sync.Once

func initConfig() {
	cnf = &Config{}
	if err := defaults.Set(cnf); err != nil {
		util.GetLog().Error("failed to set default config", zap.Error(err))
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(home + "/algo/algo.toml")
	if err != nil {
		if err = toml.Unmarshal(data, cnf); err != nil {
			util.GetLog().Error("failed to unmarshal config file", zap.Error(err))
			panic("failed to unmarshal config file")
		}
	}

	cnf.Dir.ExpandHome(home)
}

func GetConfig() *Config {
	once.Do(initConfig)
	return cnf
}
