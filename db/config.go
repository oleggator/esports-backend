package db

type DB struct {
	Host				string	`yaml:"host"`
	Port				uint16	`yaml:"port"`
	Database			string	`yaml:"database"`
	User				string	`yaml:"user"`
	Password			string	`yaml:"password"`
	ConnectionsCount	int		`yaml:"connections"`
}

type ShardConfig struct {
	Name	string	`yaml:"name"`
	Master	DB		`yaml:"master"`
	Slaves	[]DB	`yaml:"slaves"`
}

type DBConfig struct {
	Shards	[]ShardConfig	`yaml:"shards"`
}
