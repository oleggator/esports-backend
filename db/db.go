package db

import (
	"github.com/jackc/pgx"
	"log"
)

var shards []*Shard

func InitDB(config DBConfig) (err error) {
	for _, shardConfig := range config.Shards {
		shard, err := NewShard(shardConfig)
		if err != nil {
			return err
		}

		shards = append(shards, shard)
	}


	return nil
}

func Read(index int) (conn *pgx.ConnPool) {
	return shards[index].Read()
}

func Write(index int) (conn *pgx.ConnPool) {
	return shards[index].Write()
}

func GetShardsCount() int {
	return len(shards)
}

func Close() {
	log.Println("db connection closed")
	for i, _ := range shards {
		shards[i].Close()
	}
}

func prepareStatements(conn *pgx.Conn) error {
	prepareGames(conn)
	prepareTeams(conn)

	return nil
}
