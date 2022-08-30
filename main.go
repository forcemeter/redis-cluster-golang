package main

import (
	"context"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"time"

	//If you are using Redis 6, install go-redis/v8
	//"github.com/go-redis/redis/v8"

	//If you are using Redis 7, install go-redis/v9:
	"github.com/go-redis/redis/v9"
)

func main() {
	var ctx = context.Background()
	standalone(ctx)
	cluster(ctx)
}

func UniversalClient(ctx context.Context, rdb redis.UniversalClient) {
	glog.Info(ctx, strings.Repeat("=", 20))

	if r := rdb.Info(ctx, "cluster"); strings.Contains(r.String(), "cluster_enabled:1") {
		glog.Info(ctx, "cluster mode")
	} else {
		glog.Info(ctx, "standalone cluster mode")
	}

	timeStr := gtime.Now().Format("[Y-m-d,H:i:s]")

	glog.Info(ctx, rdb.ClientGetName(ctx))
	glog.Info(ctx, rdb.ClientID(ctx))
	glog.Info(ctx, rdb.ClientList(ctx))

	glog.Info(ctx, rdb.Set(ctx, "a", timeStr, -1))
	glog.Info(ctx, rdb.Get(ctx, "a"))

	glog.Info(ctx, rdb.MSet(ctx, "{hashTag1}:a", timeStr, "{hashTag1}:b", timeStr))
	glog.Info(ctx, rdb.MGet(ctx, "{hashTag1}:a", "{hashTag1}:b"))

	// SET key value EX 10 NX
	set, err := rdb.SetNX(ctx, "key", "value", 10*time.Second).Result()
	checkErr(ctx, err, set)

	// SET key value keepttl NX
	set, err = rdb.SetNX(ctx, "key", "value", redis.KeepTTL).Result()
	checkErr(ctx, err, set)

	// SORT list LIMIT 0 2 ASC
	vals1, err := rdb.Sort(ctx, "list", &redis.Sort{Offset: 0, Count: 2, Order: "ASC"}).Result()
	checkErr(ctx, err, vals1)

	// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
	vals2, err := rdb.ZRangeByScoreWithScores(ctx, "zset", &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  2,
	}).Result()
	checkErr(ctx, err, vals2)

	// ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
	vals3, err := rdb.ZInterStore(ctx, "{hashTag1}:out", &redis.ZStore{
		Keys:    []string{"{hashTag1}:zset1", "{hashTag1}:zset2"},
		Weights: []float64{2, 3},
	}).Result()
	checkErr(ctx, err, vals3)

	// EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	vals4, err := rdb.Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()
	checkErr(ctx, err, vals4)

	// custom command
	res, err := rdb.Do(ctx, "set", "key", "value").Result()
	checkErr(ctx, err, res)
}

func standalone(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	UniversalClient(ctx, rdb)
}

func cluster(ctx context.Context) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"172.29.0.2:6379",
			"172.29.0.3:6379",
			"172.29.0.4:6379",
		},
		Password: "bitnami", // no password set
	})

	UniversalClient(ctx, rdb)
}

func checkErr(ctx context.Context, err error, opt ...interface{}) {
	glog.Info(ctx, opt)
	if err != nil {
		glog.Error(ctx, err)
	}
}
