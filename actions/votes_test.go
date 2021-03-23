package actions

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func Test_ArticleVote(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	article := "article:130000"
	articleTime := time.Now().Unix()

	// Setup
	rdb.ZAdd(ctx, "time:", &redis.Z{Score: float64(articleTime), Member: article})

	err := ArticleVote(rdb, "user:230000", article)
	if err != nil {
		t.Errorf("Error found ==> %v\n", err.Error())
	}

	// ---
	userVotes, err := rdb.SMembers(ctx, "voted:130000").Result()
	if err != nil {
		t.Errorf("Error found ==> %v\n", err.Error())
	}

	if len(userVotes) == 0 {
		t.Error("Error, userVotes is zero")
	}

	uv := userVotes[0]
	if uv != "user:230000" {
		t.Errorf("Not equal, expected: %v - actual: %v", "user:230000", uv)
	}

	// ---
	articleVotesScore, err := rdb.ZRangeWithScores(ctx, "score:", 0, -1).Result()
	if err != nil {
		t.Errorf("Error found ==> %v\n", err.Error())
	}

	if len(articleVotesScore) == 0 {
		t.Error("Error, articleVotesScore is zero")
	}

	articleZScore := articleVotesScore[0]
	if articleZScore.Score != 432 {
		t.Errorf("Not equal, expected: %v - actual: %v", 432, articleZScore.Score)
	}

	if articleZScore.Member != article {
		t.Errorf("Not equal, expected: %v - actual: %v", article, articleZScore.Member)
	}

	// ---
	articleHashVotes, err := rdb.HGet(ctx, article, "votes").Result()
	if err != nil {
		t.Errorf("Error found ==> %v\n", err.Error())
	}

	if articleHashVotes != "1" {
		t.Errorf("Not equal, expected: %v - actual: %v", "1", articleHashVotes)
	}
}
