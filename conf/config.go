package conf

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	KeyDbAddr           = "DB_ADDR"
	KeyDATABASE         = "DB_DATABASE"
	KeyDbUser           = "DB_USER"
	KeyDbPwd            = "DB_PASSWORD"
	KeyDbPoolLimit      = "DB_POOL_LIMIT"
	KeyAddressPrefix    = "ADDRESS_PREFIX"
	KeyCoinMinimalDenom = "COIN_MINIMAL_DENOM"

	KeyLcd              = "LCD"
	KeyCoinMarketApiKey = "COINMARKET_API_KEY"

	EnvironmentDevelop = ".env"
	DefaultEnvironment = EnvironmentDevelop
)

func Get() Config {
	addrs := strings.Split(getEnv(KeyDbAddr, DefaultEnvironment), ",")
	return Config{
		Addrs:            addrs,
		Database:         getEnv(KeyDATABASE, DefaultEnvironment),
		UserName:         getEnv(KeyDbUser, DefaultEnvironment),
		Password:         getEnv(KeyDbPwd, DefaultEnvironment),
		PoolLimit:        getEnvInt(KeyDbPoolLimit, DefaultEnvironment),
		LcdUrl:           getEnv(KeyLcd, DefaultEnvironment),
		MarketApiKey:     getEnv(KeyCoinMarketApiKey, EnvironmentDevelop),
		AddresPrefix:     getEnv(KeyAddressPrefix, DefaultEnvironment),
		CoinMinimalDenom: getEnv(KeyCoinMinimalDenom, DefaultEnvironment),
	}
}

type Config struct {
	Addrs            []string
	Database         string
	UserName         string
	Password         string
	PoolLimit        int
	CoinMinimalDenom string

	LcdUrl       string
	MarketApiKey string

	AddresPrefix string
}

func getEnv(key string, environment string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getEnvInt(key string, environment string) int {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Fatalf("Error convert %s to string", key)
	}
	return value
}
