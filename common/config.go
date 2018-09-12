package common

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Configs 全局配置信息
type Configs struct {
	Listen                 string
}

// Config 全局配置信息
var Config *Configs

// InitConfig 加载配置
func InitConfig(path string) {
	config, err := loadConfig(path)
	if err != nil {
		panic(err)
	}
	Config = config
}

func loadConfig(path string) (*Configs, error) {
	config := new(Configs)
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}
	return config, nil
}

// Duration 配置中使用的时长
type Duration struct {
	time.Duration
}

// UnmarshalText 将字符串形式的时长信息转换为Duration类型
func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// D 从Duration struct中取出time.Duration类型的值
func (d *Duration) D() time.Duration {
	return d.Duration
}

