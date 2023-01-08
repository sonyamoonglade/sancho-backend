package tests

import (
	"bytes"
	"encoding/json"
	"io"
)

const baseURL = "http://localhost:8000"

func readBody(rc io.ReadCloser) []byte {
	b, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	rc.Close()
	return b
}

func newBody(b interface{}) io.Reader {
	bodyBytes, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bodyBytes)
}

func checkIsDescending(arr []int32) bool {
	for i := 1; i < len(arr); i++ {
		a, b := arr[i-1], arr[i]
		if a < b {
			return false
		}
	}
	return true
}

func buildURL(path string) string {
	return baseURL + path
}

func StringPtr(s string) *string { return &s }

func IntPtr[N int | int8 | int16 | int32 | int64](n N) *N { return &n }
