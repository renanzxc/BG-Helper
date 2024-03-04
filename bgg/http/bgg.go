package http

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/renanzxc/BG-Helper/utils/cache"
	"github.com/rotisserie/eris"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

type HTTP struct {
	client  *http.Client
	baseURL *url.URL
	cache   cache.Cache
}

func NewBGGClient(cache cache.Cache) (*HTTP, error) {
	baseURL, err := url.Parse("https://boardgamegeek.com/xmlapi2/")
	if err != nil {
		return nil, eris.Wrap(err, "Error on parse base URL")
	}

	return &HTTP{
		client:  new(http.Client),
		baseURL: baseURL,
		cache:   cache,
	}, nil
}

func (h *HTTP) GetUserCollection(ctx context.Context, username string, useCache bool) (userCollection *CollectionXML, err error) {
	var (
		req  *http.Request
		res  *http.Response
		body []byte
	)

	if req, err = h.newRequest(http.MethodGet, "/collection", map[string][]string{
		"username": {username},
	}, nil); err != nil {
		return userCollection, eris.Wrap(err, "")
	}

	if res, err = h.do(ctx, req, useCache); err != nil {
		return userCollection, eris.Wrap(err, "")
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return userCollection, eris.Wrap(err, "")
	}

	userCollection = new(CollectionXML)
	if err = xml.Unmarshal(body, userCollection); err != nil {
		return userCollection, eris.Wrap(err, "")
	}

	return
}

func (h *HTTP) generateURL(urlIn string, params ...map[string][]string) string {
	var (
		baseURL  = *h.baseURL
		finalURL = &baseURL
	)

	finalURL = finalURL.JoinPath(urlIn)
	if len(params) > 0 {
		queryParams := url.Values(params[0])
		finalURL.RawQuery = queryParams.Encode()
	}

	return finalURL.String()
}

func (h *HTTP) newRequest(method, url string, queryParams map[string][]string, body io.Reader) (req *http.Request, err error) {
	return http.NewRequest(method, h.generateURL(url, queryParams), body)
}

func (h *HTTP) do(ctx context.Context, req *http.Request, useCache bool) (res *http.Response, err error) {
	var (
		body        []byte
		bodyStr     string
		reqCacheKey = getReqCacheKey(req)
	)

	if useCache {
		if body, err = h.cache.Get(ctx, reqCacheKey); err != nil && !errors.Is(err, cache.ErrNoCache) {
			return nil, eris.Wrap(err, "Erro on cache")
		} else if err == nil {
			zap.L().Debug(reqCacheKey + " has cache!")

			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Request:    req,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		}
	}

	for {
		zap.L().Debug("Exec req " + reqCacheKey)
		if res, err = h.client.Do(req); err != nil {
			return
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		bodyStr = string(body)
		if strings.Contains(bodyStr, "Please try again later for access.") {
			time.Sleep(time.Second * 3)
			continue
		}

		if strings.Contains(bodyStr, "Rate limit exceeded.") {
			zap.L().Debug("rate limit")
			time.Sleep(time.Second * 3)
			continue
		}

		res.Body = io.NopCloser(bytes.NewReader(body))

		break
	}

	if err = h.cache.Set(ctx, reqCacheKey, body, true); err != nil {
		return nil, eris.Wrap(err, "Erro on cache")
	}
	zap.L().Debug(reqCacheKey + " cache created!")

	return
}

func getReqCacheKey(req *http.Request) string {
	return fmt.Sprintf("%s-%s", req.Method, req.URL.String())
}
