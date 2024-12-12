package gorm

import (
	"fmt"
	"github.com/aiechoic/admin/core/ioc"
	"github.com/aiechoic/admin/core/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

const DefaultConfig = "gorm"

var initConfig = `
# GORM configuration file

# database driver (e.g., postgres, mysql, sqlite, sqlserver)
driver: "postgres"

# data source name
# e.g.,
# postgres dsn: "host=localhost user=USERNAME password=PASSWORD dbname=DBNAME port=5432 sslmode=disable TimeZone=Asia/Shanghai"
# mysql dsn: "USERNAME:PASSWORD@tcp(localhost:3306)/DBNAME?charset=utf8mb4&parseTime=True&loc=Local"
# sqlite dsn: "test.db"
# sqlserver dsn: "sqlserver://USERNAME:PASSWORD@localhost:9930?database=DBNAME"
dsn: ""

# settings for database connection pool
maxIdleConns: 2
maxOpenConns: 100
connMaxLifetime: -1 # seconds, -1 means forever

# log level for GORM logger, 1 - Silent, 2 - Error, 3 - Warn, 4 - Info
logLevel: 4

# table prefix
tablePrefix: ""

# singular table
singularTable: false

# automatically migrate models, be careful to use it in production environment
autoMigrate: true
`

type CloseAbleGormDB struct {
	*gorm.DB
	autoMigrate bool
}

func (c *CloseAbleGormDB) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

type Config struct {
	Driver          string `mapstructure:"driver"`
	DSN             string `mapstructure:"dsn"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	LogLevel        int    `mapstructure:"logLevel"`
	TablePrefix     string `mapstructure:"tablePrefix"`
	SingularTable   bool   `mapstructure:"singularTable"`
	AutoMigrate     bool   `mapstructure:"autoMigrate"`
}

// Providers defines a providers for gorm.DB, it can be redefined
var Providers = ioc.NewProviders(func(name string, args ...any) *ioc.Provider[*CloseAbleGormDB] {
	return ioc.NewProvider(func(c *ioc.Container) (*CloseAbleGormDB, error) {
		vp, err := viper.GetViper(name, initConfig, c)
		if err != nil {
			return nil, err
		}
		var cfg Config
		err = vp.Unmarshal(&cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal '%s' config: %w", name, err)
		}

		var dial gorm.Dialector

		switch cfg.Driver {
		case "postgres":
			dial = postgres.Open(cfg.DSN)
		case "mysql":
			dial = mysql.Open(cfg.DSN)
		case "sqlite":
			dial = sqlite.Open(cfg.DSN)
		case "sqlserver":
			dial = sqlserver.Open(cfg.DSN)
		default:
			return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
		}
		db, err := gorm.Open(dial, &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   cfg.TablePrefix,
				SingularTable: cfg.SingularTable,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect database: %w", err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get sql.DB: %w", err)
		}
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
		return &CloseAbleGormDB{DB: db, autoMigrate: cfg.AutoMigrate}, nil
	})
})

func GetDB(name string, c *ioc.Container, models ...any) (*gorm.DB, error) {
	db, err := Providers.GetProvider(name).Get(c)
	if err != nil {
		return nil, err
	}
	if db.autoMigrate {
		for _, model := range models {
			err = db.AutoMigrate(model)
			if err != nil {
				return nil, err
			}
		}
	}
	return db.DB, nil
}

func GetDefaultDB(c *ioc.Container, models ...any) (*gorm.DB, error) {
	return GetDB(DefaultConfig, c, models...)
}
