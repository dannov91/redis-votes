package actions

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ONE_WEEK_IN_SECONDS = 7 * 86400
	VOTE_SCORE          = 432
	ctx                 = context.Background()
)

func ArticleVote(conn *redis.Client, user string, article string) error {
	today := time.Now().Unix()
	cutoff := int(today) - ONE_WEEK_IN_SECONDS

	articleScore, err := conn.ZScore(ctx, "time:", article).Result()
	if err != nil {
		return err
	}

	if articleScore < float64(cutoff) {
		return nil
	}

	articleID := strings.Split(article, ":")[1]

	res, err := conn.SAdd(ctx, "voted:"+articleID, user).Result()
	if err != nil {
		return err
	}

	if res > 0 { // if set value does not exist...
		if _, err := conn.ZIncrBy(ctx, "score:", float64(VOTE_SCORE), article).Result(); err != nil {
			return err
		}

		if _, err := conn.HIncrBy(ctx, article, "votes", 1).Result(); err != nil {
			return err
		}
	}

	return nil
}
