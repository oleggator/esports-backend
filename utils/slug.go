package utils

import (
	"strings"
	"strconv"
	"github.com/oleggator/esports-backend/db"
	"errors"
)

func GetShardIndex(slug *string) (shardIndex int, err error) {
	if slug == nil {
		return -1, errors.New("Slug is nil")
	}

	i := strings.IndexRune(*slug, '-')
	if i == -1 {
		return -1, errors.New("Wrong slug")
	}

	shardIndex, err = strconv.Atoi((*slug)[:i])
	if err != nil {
		return -1, errors.New("Wrong slug")
	}

	if shardIndex >= db.GetShardsCount() {
		return -1, errors.New("Wrong slug")
	}

	return shardIndex, nil
}