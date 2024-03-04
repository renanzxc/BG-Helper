package cache

import (
	"context"
	"github.com/rotisserie/eris"
)

var ErrNoCache = eris.New("Não possui cache")

type Cache interface {
	Get(ctx context.Context, key string) (bytes []byte, err error)
	Set(ctx context.Context, key string, value []byte, shortTime bool)
	Down() error
}
