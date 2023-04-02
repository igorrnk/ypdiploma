package configs

import (
	"flag"
	"os"
)

type ServerConfigType struct {
	ServerAddress string `env:"RUN_ADDRESS"`
}

type DBConfigType struct {
	DBDriverName string
	DBUri        string `env:"DATABASE_URI"`
}

type ConfigType struct {
	DBConfigType
	ServerConfigType
	AccrualSysAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func InitConfig() (*ConfigType, error) {
	address := flag.String("a", "", "Run address")
	databaseURI := flag.String("d", "", "Database URI")
	accrualSysAddress := flag.String("r", "", "Accrual system address")
	config := &ConfigType{
		DBConfigType: DBConfigType{
			DBDriverName: defaultDBConfig.DBDriverName,
			DBUri:        *databaseURI,
		},
		ServerConfigType: ServerConfigType{
			ServerAddress: *address,
		},
		AccrualSysAddress: *accrualSysAddress,
	}
	config.ServerAddress = os.Getenv("RUN_ADDRESS")
	config.DBUri = os.Getenv("DATABASE_URI")
	config.AccrualSysAddress = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	return config, nil
}
