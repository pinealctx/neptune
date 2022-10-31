package dl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// Download : download a url object then save it to file
// uri -- url string
// dir -- current file dir
// for instance
//     uri -> https://example.com/1.jpg
//     dir -> /home/pics
//     would download file then save as /home/pics/1.jpg
func Download(uri string, dir string) error {
	var fileURL, err = url.Parse(uri)
	if err != nil {
		return err
	}
	var uPath = fileURL.Path
	var sgs = strings.Split(uPath, "/")
	var l = len(sgs)
	if l <= 1 {
		return fmt.Errorf("invalid.url:%s", uri)
	}
	var fileName = path.Join(dir, sgs[l-1])

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	var cli = http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	rsp, err := cli.Get(uri)
	if err != nil {
		_ = file.Close()
		_ = os.Remove(fileName)
		return err
	}
	defer func() {
		_ = rsp.Body.Close()
		_ = file.Close()
	}()
	_, err = io.Copy(file, rsp.Body)
	return err
}
