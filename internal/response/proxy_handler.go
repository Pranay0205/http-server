package response

import (
	"http-server/internal/request"
)

var Proxies = map[string]string{
	"/httpbin/": "https://httpbin.org/",
}

func proxyhandler(req request.RequestLine) (*Writer, error) {

	// if strings.HasPrefix()

	return &Writer{}, nil
}
