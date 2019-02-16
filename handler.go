package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// slice of url mappings
type urlMap []map[string]string

func parseYAML(yml []byte) (urlMap, error) {
	var u urlMap
	err := yaml.Unmarshal(yml, &u)
	if err != nil {
		return nil, fmt.Errorf("parseYAML: %s", err)
	}
	return u, nil
}

func buildMap(urls urlMap) map[string]string {
	m := make(map[string]string)
	for _, u := range urls {
		path, url := u["path"], u["url"]
		m[path] = url
	}
	return m
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(url))
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	y, err := parseYAML(yml)
	if err != nil {
		return nil, fmt.Errorf("YAMLHandler: %s", err)
	}
	m := buildMap(y)
	return MapHandler(m, fallback), nil
}
