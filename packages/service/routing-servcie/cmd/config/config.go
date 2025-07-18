package config

import "os"


type Config struct{
	Routing_host string

}
func NewConfig()*Config{
	return  &Config{
		Routing_host: os.Getenv("ROUTING_HOST"),
	}
}