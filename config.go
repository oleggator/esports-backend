package main

import "github.com/oleggator/esports-backend/db"

type HTTPConfig struct {
	Method	string	`yaml:"method"`
	Address	string	`yaml:"address"`
}

type Config struct {
	HTTP	HTTPConfig	`yaml:"http"`
	DB		db.DBConfig	`yaml:"db"`
}
