// Images transformations API
//
// The purpose of this API is to provide a set of
// endpoints that will transform and optimise images.
// Then it becomes easy to use the API with <img> and <picture> tags in web development.
//
// Version: 2
// Schemes: https
// Host: pixboost.com
// BasePath: /api/2/
// Security:
// - api_key:
// SecurityDefinitions:
//   api_key:
//     type: apiKey
//     name: auth
//     in: query
// swagger:meta
package main

import (
	"flag"
	"github.com/Pixboost/transformimgs/v6/img"
	"github.com/Pixboost/transformimgs/v6/img/loader"
	"github.com/Pixboost/transformimgs/v6/img/processor"
	"github.com/dooman87/kolibri/health"
	"net/http"
	"os"
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
	flag.IntVar(&procNum, "proc", runtime.NumCPU(), "Number of images processors to run. Defaults to number of CPUs")
	flag.Parse()

	p, err := processor.NewImageMagick(im, imIdent)

	if err != nil {
		img.Log.Errorf("Can't create image magic processor: %+v", err)
		os.Exit(1)
	}

	img.CacheTTL = cache
	srv, err := img.NewService(&loader.Http{}, p, procNum)
	if err != nil {
		img.Log.Errorf("Can't create image service: %+v", err)
		os.Exit(2)
	}

	router := srv.GetRouter()
	router.HandleFunc("/health", health.Health)

	img.Log.Printf("Running the applicaiton on port 8080...\n")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		img.Log.Errorf("Error while stopping application: %+v", err)
		os.Exit(3)
	}
	os.Exit(0)
}
