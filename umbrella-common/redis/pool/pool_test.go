package pool

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	radix "github.com/fzzy/radix/redis"
	"github.com/stretchr/testify/assert"

	"umbrella-go/umbrella-common/redis"
)

func TestPool(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	var lastRequest redis.Request
	cw := func(next redis.Cmder, ctx context.Context, cmd string, args ...interface{}) *radix.Reply {
		lastRequest.Cmd = cmd
		lastRequest.Args = args
		return next(ctx, cmd, args...)
	}
	pool, err := NewPool("tcp", s.Addr(), 10, redis.WrapCmder(cw))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Empty()

	c, err := pool.Get()
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.Cmd(context.Background(), "ping").Str()
	assert.Equal(t, "PONG", r)
	pool.CarefullyPut(c, &err)

	assert.Equal(t, "ping", lastRequest.Cmd)
	assert.Len(t, lastRequest.Args, 0)
}

func TestLazyPool(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	lazyPool := NewCustomLazyPool("tcp", s.Addr(), 2, radix.Dial)
	const N = 10
	pools := make(chan *Pool, N)
	for i := 0; i < N; i++ {
		go func() {
			p, err := lazyPool.GetPool()
			if err != nil {
				t.Fatal(err)
			}
			pools <- p
		}()
	}

	pool := <-pools
	for i := 1; i < N; i++ {
		assert.Equal(t, pool, <-pools)
	}

	c, err := pool.Get()
	if err != nil {
		t.Fatal(err)
	}

	r, err := c.Cmd(context.Background(), "ping").Str()
	assert.Equal(t, "PONG", r)
	pool.CarefullyPut(c, &err)
}
