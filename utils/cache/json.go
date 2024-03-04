package cache

import (
	"context"
	"encoding/json"
	"github.com/rotisserie/eris"
	"os"
	"path"
	"sync"
	"time"
)

type JSONCache struct {
	cache *map[string]valueJSONCache
	mutex sync.Mutex
	file  *os.File
}

type valueJSONCache struct {
	Exp   *int64 `json:"exp"`
	Value []byte `json:"value"`
}

func NewJSONCache(cachePath string) (jsonCache *JSONCache, err error) {
	var cache = make(map[string]valueJSONCache)

	file, err := os.Create(path.Join(cachePath, "cache.json"))
	if err != nil {
		return nil, eris.Wrap(err, "Error on generate json cache")
	}

	jsonCache = &JSONCache{cache: &cache, file: file}

	return
}

func (j *JSONCache) Set(ctx context.Context, key string, value []byte, shortTime bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var (
		exp *int64
	)

	if shortTime {
		expI := time.Now().Add(time.Hour * 24 * 7).Unix()
		exp = &expI
	}
	(*j.cache)[key] = valueJSONCache{Value: value, Exp: exp}

	return nil
}

func (j *JSONCache) Get(ctx context.Context, key string) (bytes []byte, err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	jsonValue, ok := (*j.cache)[key]
	if !ok {
		return []byte{}, ErrNoCache
	}
	if jsonValue.Exp != nil && time.Since(time.Unix(*jsonValue.Exp, 0)) > 0 {
		delete(*j.cache, key)
		return []byte{}, ErrNoCache
	}
	return jsonValue.Value, nil
}

func (j *JSONCache) Down() (err error) {
	defer j.file.Close()

	encoder := json.NewEncoder(j.file)
	if err = encoder.Encode(j.cache); err != nil {
		return eris.Wrap(err, "Error on save json cache")
	}

	return
}
