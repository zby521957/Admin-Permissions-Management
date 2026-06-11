package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 总配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   //Server：服务启动相关配置（如端口）
	Database DatabaseConfig `mapstructure:"database"` //Database：数据库连接相关配置。
	JWT      JWTConfig      `mapstructure:"jwt"`      //JWT：JWT 签名密钥等认证配置。
}

type ServerConfig struct {
	Port string `mapstructure:"port"` //控制 HTTP 服务器的监听端口。
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"` //存储 JWT 签名所需的密钥
}

// 声明一个包级别变量（全局变量），类型为 *Config
var AppConfig *Config

func Load() error {
	viper.SetConfigName("config") // 文件名（不含扩展名）
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 搜索路径（当前目录）
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	//log.Println("配置文件读取失败")

	//创建一个新的 Config 结构体实例，并取它的指针，赋值给全局变量 AppConfig
	AppConfig = &Config{}
	//把 Viper 从 config.yaml 中读取到的所有配置数据，自动填充到 AppConfig 结构体中
	if err := viper.Unmarshal(AppConfig); err != nil {
		return err
	}
	log.Println("配置加载成功")
	return nil
}
