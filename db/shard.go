package db

import (
	"github.com/jackc/pgx"
	"sync"
	"fmt"
)

type Shard struct {
	master *pgx.ConnPool
	slaves []*pgx.ConnPool
	slavesCount int

	mu sync.Mutex
	readIndex int
}

func NewShard(config ShardConfig) (shard *Shard, err error) {
	shard = new(Shard)
	shard.master, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host: config.Master.Host,
			Port:     config.Master.Port,
			Database: config.Master.Database,
			User:     config.Master.User,
			Password: config.Master.Password,
		},
		MaxConnections: config.Master.ConnectionsCount,
		AfterConnect: prepareStatements,
	})
	if err != nil {
		return nil, err
	}

	shard.slavesCount = len(config.Slaves)
	shard.slaves = make([]*pgx.ConnPool, shard.slavesCount)
	for i, slaveConfig := range config.Slaves {
		slave, err := pgx.NewConnPool(pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host: slaveConfig.Host,
				Port:     slaveConfig.Port,
				Database: slaveConfig.Database,
				User:     slaveConfig.User,
				Password: slaveConfig.Password,
			},
			MaxConnections: slaveConfig.ConnectionsCount,
			AfterConnect: prepareStatements,
		})
		if err != nil {
			return nil, err
		}

		shard.slaves[i] = slave
	}

	return shard, nil
}

func (s *Shard) Read() (conn *pgx.ConnPool) {
	s.mu.Lock()
	if s.readIndex >= s.slavesCount {
		fmt.Println("read from master")
		conn = s.master
		s.readIndex = 0
	} else {
		fmt.Println("read from slaves")
		conn = s.slaves[s.readIndex]
		s.readIndex++
	}
	s.mu.Unlock()

	return conn
}

func (s *Shard) Write() (conn *pgx.ConnPool) {
	return s.master
}

func (s *Shard) Close() {
	s.master.Close()
	for i, _ := range s.slaves {
		s.slaves[i].Close()
	}
}