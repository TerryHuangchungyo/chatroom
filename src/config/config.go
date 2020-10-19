package config

import "github.com/BurntSushi/toml"

const (
	MEMBER_MODERATOR = iota
	MEMBER_GENERAL
)

var DATABASE database
var REDIS redis

func init() {
	var c config
	filePath := "./config.toml"
	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		panic(err)
	}
	DATABASE = c.Db
	REDIS = c.Redis
}

type config struct {
	Redis redis    `toml:"redis"`
	Db    database `toml:"database"`
}

type database struct {
	Connection string
	Host       string
	Port       int
	User       string
	Password   string
	Dbname     string
	Charset    string
	Collation  string
}

type redis struct {
	Host                string
	Port                int
	Password            string
	Db                  int
	ChannelKeyPrefix    string
	HubHistoryKeyPrefix string
}
