package database

import "testing"

func TestGetRedisReturnsNilWhenUninitialized(t *testing.T) {
	redisClient = nil
	if got := GetRedis(); got != nil {
		t.Fatalf("GetRedis() = %v, want nil", got)
	}
}
