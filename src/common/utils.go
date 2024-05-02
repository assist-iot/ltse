package common

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func GetEnvironment(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = def
	}
	return val
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func ReverseProxy(urlRemote string, context *gin.Context, proxyPath string) error {
	if proxyPath == "" {
		proxyPath = "/"
	}

	remote, err := url.Parse(urlRemote)

	if err != nil {
		log.Println(err)
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	//Define the director func
	//This is a good place to log, for example
	proxy.Director = func(req *http.Request) {
		req.Header = context.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = proxyPath
	}

	proxy.ServeHTTP(context.Writer, context.Request)

	return nil
}

func RemoveDuplicates(items []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range items {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
