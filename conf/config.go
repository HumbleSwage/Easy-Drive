package conf

import (
	"easy-drive/types"
	"github.com/spf13/viper"
	"os"
)

var Conf *Config

type Config struct {
	Server     *types.Server           `yaml:"server"`
	Mysql      map[string]*types.Mysql `yaml:"mysql"`
	Redis      *types.Redis            `yaml:"redis"`
	Email      *types.Email            `yaml:"email"`
	UploadPath *types.UploadPath       `yaml:"uploadPath"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/conf/local/")
	viper.AddConfigPath(workDir)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
