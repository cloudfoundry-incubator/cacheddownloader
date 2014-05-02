package cacheddownloader

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileCache interface {
	Fetch(url *url.URL, cacheKey string) (io.ReadCloser, error)
}

type CachedFile struct {
	size   int64
	access time.Time
}

type Cache struct {
	cachedPath     string
	uncachedPath   string
	maxSizeInBytes int64
	downloader     Downloader
	lock           *sync.Mutex

	cachedFiles map[string]CachedFile
}

func New(cachedPath string, uncachedPath string, maxSizeInBytes int64, downloader Downloader) *Cache {
	os.RemoveAll(cachedPath)
	os.MkdirAll(cachedPath, 0770)
	return &Cache{
		cachedPath:     cachedPath,
		uncachedPath:   uncachedPath,
		maxSizeInBytes: maxSizeInBytes,
		downloader:     downloader,
		lock:           &sync.Mutex{},
		cachedFiles:    map[string]CachedFile{},
	}
}

func (c *Cache) Fetch(url *url.URL, cacheKey string) (io.ReadCloser, error) {
	if cacheKey == "" {
		return c.fetchUncachedFile(url)
	} else {
		cacheKey = fmt.Sprintf("%x", md5.Sum([]byte(cacheKey)))
		return c.fetchCachedFile(url, cacheKey)
	}
}

func (c *Cache) fetchUncachedFile(url *url.URL) (io.ReadCloser, error) {
	destinationFile, err := ioutil.TempFile(c.uncachedPath, "uncached")
	if err != nil {
		return nil, err
	}

	_, _, err = c.downloader.Download(url, destinationFile, time.Time{})
	if err != nil {
		os.Remove(destinationFile.Name())
		return nil, err
	}
	os.Remove(destinationFile.Name()) //OK, 'cause that's how unix works
	destinationFile.Seek(0, 0)

	return destinationFile, nil
}

func (c *Cache) fetchCachedFile(url *url.URL, cacheKey string) (io.ReadCloser, error) {
	c.recordAccessForCacheKey(cacheKey)

	path := c.pathForCacheKey(cacheKey)

	f, err := os.Open(path)
	fileExists := err == nil

	modTime := time.Time{}
	if fileExists {
		info, err := os.Stat(path)
		if err != nil {
			f.Close()
			return nil, err
		}
		modTime = info.ModTime()
	}

	//download the file to a temporary location
	tempFile, err := ioutil.TempFile(c.uncachedPath, "temporary")
	if err != nil {
		if fileExists {
			f.Close()
		}
		return nil, err
	}
	defer os.Remove(tempFile.Name()) //OK, even if we return tempFile 'cause that's how UNIX works.

	didDownload, size, err := c.downloader.Download(url, tempFile, modTime)
	if err != nil {
		if fileExists {
			f.Close()
		}
		return nil, err
	}

	if !didDownload {
		return f, nil
	}

	//make room for the file and move it in (if possible)
	c.moveFileIntoCache(cacheKey, tempFile.Name(), size)

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func (c *Cache) pathForCacheKey(cacheKey string) string {
	return filepath.Join(c.cachedPath, cacheKey)
}

func (c *Cache) moveFileIntoCache(cacheKey string, sourcePath string, size int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if size > c.maxSizeInBytes {
		//file does not fit in cache...
		return
	}

	usedSpace := int64(0)
	for ck, f := range c.cachedFiles {
		if ck != cacheKey {
			usedSpace += f.size
		}
	}

	for c.maxSizeInBytes < usedSpace+size {
		oldestAccessTime, oldestCacheKey := time.Now(), ""
		for ck, f := range c.cachedFiles {
			if ck != cacheKey {
				if f.access.Before(oldestAccessTime) {
					oldestCacheKey = ck
					oldestAccessTime = f.access
				}
			}
		}
		os.Remove(c.pathForCacheKey(oldestCacheKey))
		usedSpace -= c.cachedFiles[oldestCacheKey].size
		delete(c.cachedFiles, oldestCacheKey)
	}

	f := c.cachedFiles[cacheKey]
	f.size = size
	c.cachedFiles[cacheKey] = f
	os.Rename(sourcePath, c.pathForCacheKey(cacheKey))
}

func (c *Cache) recordAccessForCacheKey(cacheKey string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	f := c.cachedFiles[cacheKey]
	f.access = time.Now()
	c.cachedFiles[cacheKey] = f
}