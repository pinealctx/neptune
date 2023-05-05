package dl

import (
	"fmt"
	"github.com/pinealctx/neptune/tex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	defaultByteBuff = 64 * 1024 // 64k
)

type HTTPObj struct {
	// http response
	Resp *http.Response
	// object name
	Name string
}

func (h *HTTPObj) Close() error {
	return h.Resp.Body.Close()
}

// DownloadObj : download a url object, return http.Response and object name.
// Warning: Call "defer HTTPObj.Resp.Body.Close()" at the beginning.
// uri -- url string
// for instance
//
//	uri -> https://example.com/1.jpg
//	    HTTPObj.Resp -> http.Response
//	    HTTPObj.Name -> 1.jpg
//	    err -> error return when failed
func DownloadObj(uri string) (*HTTPObj, error) {
	var fileURL, err = url.Parse(uri)
	if err != nil {
		return nil, err
	}
	var uPath = fileURL.Path
	var sgs = strings.Split(uPath, "/")
	var l = len(sgs)
	if l <= 1 {
		return nil, fmt.Errorf("invalid.url:%s", uri)
	}
	var hc = http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	rsp, err := hc.Get(uri) // nolint:bodyclose
	if err != nil {
		return nil, err
	}

	return &HTTPObj{
		Resp: rsp,
		Name: sgs[l-1],
	}, nil
}

// Download2Buffer : download a url object to bytes buffer
// uri -- url string
// for instance
//
//	uri -> https://example.com/1.jpg
//	    name -> 1.jpg
//	    data -> binary data
//	    err -> error return when failed
func Download2Buffer(uri string) (string, []byte, error) {
	var rsp, err = DownloadObj(uri)
	if err != nil {
		return "", nil, err
	}
	defer func() {
		_ = rsp.Resp.Body.Close()
	}()
	var stream *tex.Buffer
	if rsp.Resp.ContentLength > 0 {
		stream = tex.NewSizedBuffer(int(rsp.Resp.ContentLength))
	} else {
		stream = tex.NewSizedBuffer(defaultByteBuff)
	}
	_, err = stream.ReadFrom(rsp.Resp.Body)
	if err != nil {
		return "", nil, err
	}
	return rsp.Name, stream.Bytes(), nil
}

// Download2Path : download a url object then save it to file
// uri -- url string
// dir -- current file dir
// for instance
//
//	uri -> https://example.com/1.jpg
//	dir -> /home/pics
//	would download file then save as /home/pics/1.jpg
func Download2Path(uri string, dir string) error {
	var rsp, err = DownloadObj(uri)
	if err != nil {
		return err
	}
	defer func() {
		_ = rsp.Resp.Body.Close()
	}()
	var fileName = path.Join(dir, rsp.Name)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
		if err != nil {
			_ = os.Remove(fileName)
		}
	}()

	_, err = io.Copy(file, rsp.Resp.Body)
	return err
}
