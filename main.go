package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func main() {

	log.SetFlags(0) // no timestamp
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) != 1 {
		log.Fatalf("supply a chart file or directory")
	}
	chartPath := os.Args[1]

	fmt.Println("\n===== Load Helm Chart =====")
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.TypeOf(chart))
	// fmt.Print(*chart)

	fmt.Println("\n===== Populating Values =====")
	// var vals chartutil.Values

	// vals := chartutil.Values{
	// 	"replicaCount": 3,
	// }

	// Signature: func CoalesceValues(chrt *chart.Chart, vals map[string]interface{}) (Values, error)
	// throws nil pointer evaluating interface {}
	// vals, err := chartutil.CoalesceValues(chart, map[string]interface{}{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(vals.YAML())

	releaseOptions := chartutil.ReleaseOptions{Name: "release1", Namespace: "ns1"}
	// Submitting empty map param {}{}
	vals, err := chartutil.ToRenderValues(chart, map[string]interface{}{},
		releaseOptions, chartutil.DefaultCapabilities)
	// vals, err := chartutil.ToRenderValues(chart, map[string]interface{}{},
	// 	chartutil.ReleaseOptions{}, chartutil.DefaultCapabilities)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// fmt.Println(vals.YAML())

	fmt.Println("\n===== Helm Templating ======")

	// Using Method outputs trailing nil
	// e := engine.Engine{Strict: false, LintMode: false}
	// fmt.Println(e.Render(chart, vals))

	// m becomes map[string]string
	m, err := engine.Render(chart, vals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)

	fmt.Println("\n===== Chart files found =====")
	// Identify chart's filenames (keys of m)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
		fmt.Println(k)
	}

	fmt.Println("\n===== Searching images in K8S manifests =====\n")
	// Idenitfy file contents (values of m)
	for _, k := range keys {
		if strings.Contains(m[k], "image:") {
			fmt.Printf("=== Image found in %s ===\n", k)
			fmt.Println(m[k])
		}
	}

}
