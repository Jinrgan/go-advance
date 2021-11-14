package captcha

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Verifier struct {
	Redis *redis.Client
}

// TODO: refactor to repo
func (v *Verifier) Verify(mobile, code string) error {
	val, err := v.Redis.Get(mobile).Result()
	if err != nil {
		return err
	}

	if val != code {
		return fmt.Errorf("wrong code")
	}

	return nil
}
