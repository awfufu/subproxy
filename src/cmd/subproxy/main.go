package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Name   string `yaml:"name"`
	Suburl string `yaml:"suburl"`
	Proxy  string `yaml:"proxy"`
}

type Config struct {
	Listen string  `yaml:"listen"`
	Routes []Route `yaml:"routes"`
}

func main() {
	configPath := flag.String("f", "config.yaml", "path to config file")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Fatalf("unmarshal config: %v", err)
	}

	mux := http.NewServeMux()

	for _, r := range conf.Routes {
		target, err := url.Parse(r.Suburl)
		if err != nil {
			log.Printf("invalid suburl for %s: %v", r.Name, err)
			continue
		}

		transport := http.DefaultTransport.(*http.Transport).Clone()
		if r.Proxy != "" {
			proxyURL, _ := url.Parse(r.Proxy)
			transport.Proxy = http.ProxyURL(proxyURL)
			log.Printf("add route /%s with proxy", r.Name)
		} else {
			log.Printf("add route /%s", r.Name)
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
			req.URL.Path = target.Path
			req.URL.RawQuery = target.RawQuery
		}

		routePath := "/" + r.Name
		mux.HandleFunc(routePath, func(w http.ResponseWriter, req *http.Request) {
			proxy.ServeHTTP(w, req)
		})
	}

	fmt.Printf("listening on %s\n", conf.Listen)
	if err := http.ListenAndServe(conf.Listen, mux); err != nil {
		log.Fatal(err)
	}
}
