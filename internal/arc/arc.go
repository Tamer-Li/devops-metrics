package arc

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GZIPHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func RZIPBody(w http.ResponseWriter, r *http.Request) (io.Reader, error) {
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	return reader, nil
}

func Compress(data []byte) (*bytes.Buffer, error) {
	var compressedData bytes.Buffer
	gzWriter := gzip.NewWriter(&compressedData)

	_, err := gzWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("Failed data to gzip: %w", err)
	}

	err = gzWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed close gzip: %w", err)
	}

	return &compressedData, nil
}
