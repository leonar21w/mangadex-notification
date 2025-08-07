package constants

import "time"

type redisConstants struct{}

var Redis = redisConstants{}

func (r redisConstants) PoolSize() int               { return 5 }
func (r redisConstants) DialTimeout() time.Duration  { return 6 * time.Second }
func (r redisConstants) ReadTimeout() time.Duration  { return 5 * time.Second }
func (r redisConstants) WriteTimeout() time.Duration { return 5 * time.Second }
