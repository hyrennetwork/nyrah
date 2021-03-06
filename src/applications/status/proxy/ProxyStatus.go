package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net"
	"net/hyren/nyrah/cache/local"
	"time"

	NyrahProvider "net/hyren/nyrah/misc/providers"
)

func GetApplicationOnlinePlayers(application string) (int, error) {
	redisConnection := NyrahProvider.REDIS_MAIN.Provide().Get()

	defer redisConnection.Close()

	var bytes, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	if err != nil {
		return 0, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return 0, err
	}
	return int(data["online_players"].(float64)), nil
}

func GetApplicationAddress(application string) (string, error) {
	redisConnection := NyrahProvider.REDIS_MAIN.Provide().Get()

	defer redisConnection.Close()

	var bytes, err = redis.Bytes(
		redisConnection.Do("GET", fmt.Sprintf("applications:%s", application)),
	)

	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		return "", err
	}

	address := data["inet_socket_address"].(map[string]interface{})

	return fmt.Sprintf(
		"%s:%d",
		address["host"].(string),
		int(address["port"].(float64)),
	), nil
}

func IsProxyOnline(server string) bool {
	_, err := net.Dial("tcp", server)

	if err != nil {
		return false
	}

	return true
}

func GetOnlinePlayers() int {
	onlinePlayers, found := local.CACHE.Get("online_players")

	if !found {
		redisConnection := NyrahProvider.REDIS_MAIN.Provide().Get()

		defer redisConnection.Close()

		cursor := 0

		for ok := true; ok; ok = cursor != 0 {
			result, err := redis.Values(
				redisConnection.Do("SCAN", cursor, "MATCH", "users:*"),
			)

			if err != nil {
				panic(err)
			}

			cursor, _ = redis.Int(result[0], nil)
			keys, _ := redis.Strings(result[1], nil)

			if onlinePlayers == nil {
				onlinePlayers = 0 + len(keys)
			} else {
				onlinePlayers = onlinePlayers.(int) + len(keys)
			}
		}

		local.CACHE.Set("online_players", onlinePlayers, 1*time.Second)
	}

	return onlinePlayers.(int)
}
