package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type MySql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Timeout  string `yaml:"timeout"`
	Dbname   string `yaml:"dbname"`
}

type GoAdminConfig struct {
	MySql MySql `yaml:"mysql"`
}

var Config *GoAdminConfig

func InitConfig(ctx context.Context, configDir string) error {
	// 使用 Viper 库设置配置文件的类型
	viper.SetConfigType("yaml")
	if configDir == "" {
		// .代表项目根目录
		configDir = "."
	}
	// 使用 Viper 库设置要读取的配置文件路径
	viper.SetConfigFile(fmt.Sprintf("%s/config.yaml", configDir))

	// 使用 Viper 库来读取和解析配置文件
	err := viper.ReadInConfig()
	// 使用 errors 包中的 Wrap 函数来包装错误，提供更多的错误信息和上下文信息
	if err != nil {
		return errors.Wrap(err, "ReadInConfig")
	}

	// 使用 Viper 库将解析后的配置数据反序列化到 Config 变量中
	err = viper.Unmarshal(&Config)
	if err != nil {
		return errors.Wrap(err, "viper.Unmarshal")
	}

	// 使用 Go 标准库中的 json.Marshal() 函数将 Config 对象序列化为 JSON 格式的字节流
	configBuf, _ := json.Marshal(Config)
	// 使用 klog 库中的函数来记录日志信息
	klog.CtxInfof(ctx, "config: %s", string(configBuf))
	return nil
}
