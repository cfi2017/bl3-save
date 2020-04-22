package assets

import (
	"io"
	"net/http"
	"os"
	"time"
)

var (
	httpClient = http.Client{
		Timeout: time.Second * 30,
	}
	publisher = "https://bl3.swiss.dev"
)

func downloadAsset(path, url string) error {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if fileExists(path) {
		t := modTime(path)
		request.Header.Add("If-Modified-Since", t.Format(http.TimeFormat))
	}

	r, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if r.StatusCode == http.StatusNotModified {
		// don't write in this case
		return nil
	}
	defer r.Body.Close()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return err
	}
	return nil
}

func modTime(path string) time.Time {
	s, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return s.ModTime()
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
