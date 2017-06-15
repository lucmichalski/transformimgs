// Images transformations API
//
// The purpose of this API is to provide a set of
// endpoints that will transform and optimise images.
// Then it becomes easy to use the API with <img> and <picture> tags in web development.
//
// Version: 1
// Schemes: http
// Host: pixboost.com
// BasePath: /api/1/
// Security:
// - api_key:
// SecurityDefinitions:
// - api_key:
//   type: apiKey
//   name: auth
//   in: query
// swagger:meta
package main

import (
	"flag"
	"github.com/dooman87/kolibri/health"
	"github.com/dooman87/transformimgs/img"
	"log"
	"net/http"
	"runtime"
)

func main() {
	var (
		im      string
		imIdent string
		cache   int
		procNum int
	)
	flag.StringVar(&im, "imConvert", "", "Imagemagick convert command")
	flag.StringVar(&imIdent, "imIdentify", "", "Imagemagick identify command")
	flag.IntVar(&cache, "cache", 86400,
		"Number of seconds to cache image after transformation (0 to disable cache). Default value is 86400 (one day)")
	flag.IntVar(&procNum, "proc", runtime.NumCPU()*2, "Number of images processors to run. Defaults to number of CPU * 2")
	flag.Parse()

	p, err := img.NewProcessor(im, imIdent)
	if err != nil {
		log.Fatalf("Can't create image magic processor: %+v", err)
	}

	img.CacheTTL = cache
	srv, err := img.NewService(&img.ImgUrlReader{}, p, cache, procNum)
	if err != nil {
		log.Fatalf("Can't create image service: %+v", err)
	}

	http.HandleFunc("/health", health.Health)
	http.Handle("/", srv.GetRouter())

	log.Println("Running the applicaiton on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
