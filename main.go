package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

const key = "leaderboard"
const TimeLength = 100_0000_0000

var client = redis.NewClient(&redis.Options{
	Addr:     "port:6379",
	Password: "", // no password set
	DB:       15, // use default DB
})

func main() {
	Add()
	//Modify()
	fmt.Println(GetRank())
	fmt.Println(GetUserRank("user48"))
}

// Add 增加排行榜
func Add() {
	for i := 0; i < 100; i++ {
		client.ZAdd(key, redis.Z{
			Score:  packScore(int64(i + 1)),
			Member: "user" + strconv.Itoa(i),
		})
	}
}

// Modify 修改排行榜
func Modify() {
	for i := 0; i < 100; i++ {
		client.ZAdd(key, redis.Z{
			Score:  float64(i + 100),
			Member: "user" + strconv.Itoa(i),
		})
	}
}

// GetRank 查询排行榜
func GetRank() []Rank {
	list, err := client.ZRevRangeWithScores(key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	result := make([]Rank, len(list))
	for index, value := range list {
		result[index] = Rank{
			UserId: value.Member.(string),
			Score:  parseScore(value.Score),
		}
	}
	return result
}

// GetUserRank 根据玩家id查看名次和分数
func GetUserRank(userId UserId) int64 {
	rank, err := client.ZRevRank(key, userId).Result()
	if err != nil {
		panic(err)
	}
	return rank + 1
}

//封装成分数
func packScore(score int64) float64 {
	now := time.Now().Unix()
	return float64(score*TimeLength + now)
	//如果排行榜支持负数
	//if score >= 0 {
	//	return float64(score*TimeLength + now)
	//} else {
	//	return float64(score*TimeLength + now - TimeLength)
	//}
}

//将分数解析
func parseScore(score float64) int64 {
	return int64(score / TimeLength)
}
