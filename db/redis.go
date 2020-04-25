package db

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"server/config"
	"strconv"
)

var Redis *redis.Client

func Start() {
	// Get the env variables needed to connect to the db
	host := config.EnvVar("REDIS_HOST")
	port := config.EnvVar("REDIS_PORT")
	pass := config.EnvVar("REDIS_PASS")
	db, err := strconv.Atoi(config.EnvVar("REDIS_DB"))
	if err != nil {
		log.Fatalf("error while getting the redis db value from the env variables!")
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       db,
	})
}

// Return the list of keys that start with the given prefix
func ScanKeys(keyPrefix string, keySuffix string) ([]string, error) {
	var foundKeys []string
	var cursor uint64
	var err error
	for {
		var keys []string
		keys, cursor, err = Redis.Scan(cursor, keyPrefix+"*"+keySuffix, 10).Result()
		foundKeys = append(foundKeys, keys...)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning the Redis: %v", err)
		}
		if cursor == 0 {
			break
		}
	}

	return foundKeys, nil
}