package cache

import (
	"context"
	"github.com/rotisserie/eris"
)

var ErrNoCache = eris.New("No cache")

type Cache interface {
	Get(ctx context.Context, key string) (bytes []byte, err error)
	Set(ctx context.Context, key string, value []byte, shortTime bool) error
	Down() error
}
