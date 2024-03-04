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
	cache        *map[string]valueJSONCache
	cacheChanged bool
	mutex        sync.Mutex
	filename     string
}

type valueJSONCache struct {
	Exp   *int64 `json:"exp"`
	Value []byte `json:"value"`
}

func NewJSONCache(cachePath string) (jsonCache *JSONCache, err error) {
	var (
		cache    = make(map[string]valueJSONCache)
		filename = path.Join(cachePath, "cache.json")
	)

	if _, err = os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			return nil, eris.Wrap(err, "Error on read json cache")
		}
		err = nil
	} else {
		var fileData []byte

		if fileData, err = os.ReadFile(filename); err != nil {
			return nil, eris.Wrap(err, "Error on read json cache")
		}

		if err = json.Unmarshal(fileData, &cache); err != nil {
			return nil, eris.Wrap(err, "Error on read json cache")
		}
	}

	jsonCache = &JSONCache{cache: &cache, filename: filename}

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

	j.cacheChanged = true
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
	if !j.cacheChanged {
		return
	}
	
	jsonStr, err := json.Marshal(j.cache)
	if err != nil {
		return eris.Wrap(err, "Error on save json cache")
	}

	err = os.WriteFile(j.filename, jsonStr, 0666)
	if err != nil {
		return eris.Wrap(err, "Error on save json cache")
	}

	return
}
