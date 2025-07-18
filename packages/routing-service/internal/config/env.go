package config

import (
	"fmt"
	"log"
	"sync"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zenv"
)

type env struct {
	PORT       int
	DB_HOST    string
	DB_USER    string
	DB_PASS    string
	DB_NAME    string
	DB_PORT    string
	JWT_SECRET string
}

var (
	once sync.Once
	cfg  env
)

var schema = z.Struct(z.Shape{
	"PORT":       z.Int().GT(1000).LT(65535).Default(3000),
	"DB_HOST":    z.String().Default("localhost"),
	"DB_USER":    z.String().Default("postgres"),
	"DB_PASS":    z.String().Default("postgres"),
	"DB_NAME":    z.String().Default("irtiqaa"),
	"DB_PORT":    z.String().Default("5432"),
	"JWT_SECRET": z.String().Default("your-super-secret-jwt-key-change-in-production"),
})

func Get() *env {
	once.Do(func() {
		errs := schema.Parse(zenv.NewDataProvider(), &cfg)
		if errs != nil {
			fmt.Println("Failed to load environment config:")
			log.Fatal(z.Issues.SanitizeMap(errs))
		}
	})
	return &cfg
}

func (*env) DSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_USER,
		cfg.DB_PASS,
		cfg.DB_NAME,
		"disable",
	)
	return dsn
}
