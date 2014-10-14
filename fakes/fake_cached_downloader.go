// This file was generated by counterfeiter
package fakes

import (
	"io"
	"net/url"
	"sync"

	"github.com/pivotal-golang/cacheddownloader"
)

type FakeCachedDownloader struct {
	FetchStub        func(urlToFetch *url.URL, cacheKey string, transformer cacheddownloader.CacheTransformer) (io.ReadCloser, error)
	fetchMutex       sync.RWMutex
	fetchArgsForCall []struct {
		urlToFetch  *url.URL
		cacheKey    string
		transformer cacheddownloader.CacheTransformer
	}
	fetchReturns struct {
		result1 io.ReadCloser
		result2 error
	}
}

func (fake *FakeCachedDownloader) Fetch(urlToFetch *url.URL, cacheKey string, transformer cacheddownloader.CacheTransformer) (io.ReadCloser, error) {
	fake.fetchMutex.Lock()
	fake.fetchArgsForCall = append(fake.fetchArgsForCall, struct {
		urlToFetch  *url.URL
		cacheKey    string
		transformer cacheddownloader.CacheTransformer
	}{urlToFetch, cacheKey, transformer})
	fake.fetchMutex.Unlock()
	if fake.FetchStub != nil {
		return fake.FetchStub(urlToFetch, cacheKey, transformer)
	} else {
		return fake.fetchReturns.result1, fake.fetchReturns.result2
	}
}

func (fake *FakeCachedDownloader) FetchCallCount() int {
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	return len(fake.fetchArgsForCall)
}

func (fake *FakeCachedDownloader) FetchArgsForCall(i int) (*url.URL, string, cacheddownloader.CacheTransformer) {
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	return fake.fetchArgsForCall[i].urlToFetch, fake.fetchArgsForCall[i].cacheKey, fake.fetchArgsForCall[i].transformer
}

func (fake *FakeCachedDownloader) FetchReturns(result1 io.ReadCloser, result2 error) {
	fake.FetchStub = nil
	fake.fetchReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

var _ cacheddownloader.CachedDownloader = new(FakeCachedDownloader)