package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func main() {
	//chartPath := "samples/haproxy-0.3.25.tgz"
	chartPath := "samples/helm1-0.1.0.tgz"
	//chartPath := "samples/helm1"

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

	vals2 := chartutil.CoalesceTables(map[string]interface{}{}, vals)

	fmt.Println("\n===== Helm Templating ======")

	// Using Method outputs trailing nil
	// e := engine.Engine{Strict: false, LintMode: false}
	// fmt.Println(e.Render(chart, vals))

	m, err := engine.Render(chart, vals2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
		fmt.Println(k)
	}
	fmt.Println(keys)
}
