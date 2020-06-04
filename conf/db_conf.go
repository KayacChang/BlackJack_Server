package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	MysqlConf    *mysqlConf
	Debug        bool
	Singular     bool
	MaxIdleConns int
	MaxOpenConns int
)

func init() {
	viper.SetConfigName("db") // name of config file (without extension)
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/conf.d") // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	MysqlConf = &mysqlConf{
		Host:     viper.GetString("host"),
		Port:     viper.GetString("port"),
		Database: viper.GetString("db"),
		// MysqlTimeZone : viper.GetString("timezone"),
		Charset:      viper.GetString("charset"),
		User:         viper.GetString("user"),
		Passwd:       viper.GetString("passwd"),
		ParseTime:    viper.GetString("parse_time"),
		Debug:        viper.GetBool("debug"),
		Singular:     viper.GetBool("singular"),
		MaxIdleConns: viper.GetInt("max_idle_conns"),
		MaxOpenConns: viper.GetInt("max_open_conns"),
	}

	fmt.Println(MysqlConf.ConverToPath())
}

type mysqlConf struct {
	Host              string
	Port              string
	User              string
	Passwd            string
	Database          string
	Charset           string
	ColumnsWithAlias  bool
	InterpolateParams bool
	ParseTime         string
	Timeout           string
	ReadTimeout       string
	WriteTimeout      string

	Singular     bool
	Debug        bool
	MaxIdleConns int
	MaxOpenConns int
}

func (mc *mysqlConf) ConverToPath() string {
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s",
		mc.User, mc.Passwd, mc.Host, mc.Port, mc.Database, mc.Charset, mc.ParseTime)

	return path
}
