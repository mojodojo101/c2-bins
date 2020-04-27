package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	DBName   string
	Host     string
	Port     int
	Username string
	Password string
	SSLMode  string
	Timezone string
}

type ResourceConfig struct {
	TargetsPath string
	BeaconsPath string
}
type PortConfig struct {
	BeaconPort string
	ClientPort string
}
type ClientConfig struct {
	Username string
	Password string
}

func GetResourceConfig() *ResourceConfig {
	return &ResourceConfig{
		TargetsPath: "targets/",
		BeaconsPath: "beacons/",
	}
}

func GetDBConfig() *Config {
	dbname := os.Getenv("C2DBName")
	host := os.Getenv("C2DBHost")
	port, _ := strconv.Atoi(os.Getenv("C2DBPort"))
	username := os.Getenv("C2DBUsername")
	password := os.Getenv("C2DBPassword")
	tz:= os.Getenv("C2DBTimezone")
	return &Config{
		DB: &DBConfig{
			DBName:   dbname,
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			SSLMode:  "require",
			Timezone: tz,
		},
	}
}

func GetServerPorts() *PortConfig {
	bport, err := strconv.Atoi(os.Getenv("C2BeaconPort"))
	if err != nil {
		fmt.Printf("Not a valid BeaconPort check the /config/config.go")
		os.Exit(1)
	}
	bp := fmt.Sprintf(":%v", bport)
	cport, err := strconv.Atoi(os.Getenv("C2ClientPort"))
	if err != nil {
		fmt.Printf("Not a valid ClientPort check the /config/config.go")
		os.Exit(1)
	}
	cp := fmt.Sprintf(":%v", cport)

	return &PortConfig{
		BeaconPort: bp,
		ClientPort: cp,
	}
}

func GetClient() *ClientConfig {
	cu := os.Getenv("C2ClientUsername")
	cp := os.Getenv("C2ClientPassword")
	return &ClientConfig{
		Username: cu,
		Password: cp,
	}
}
