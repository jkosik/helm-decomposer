package main

import (
	"fmt"
	"reflect"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func main() {
	//chartPath := "samples/hyperbridge-1.1.300.tgz"
	//chartPath := "samples/hyperbridge-data-1.0.160.tgz"
	chartPath := "samples/nginx-12.0.0.tgz"
	//chartPath := "samples/haproxy-0.3.25.tgz"

	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n===== Helm Chart ===========")
	fmt.Print(chart)

	var vals chartutil.Values
	// To override Chart values
	// vals := chartutil.Values{
	// 	"replicaCount": 3,
	// }

	// Coalesce values
	// func CoalesceValues(chrt *chart.Chart, vals map[string]interface{}) (Values, error)
	coalescedVals, err := chartutil.CoalesceValues(chart, vals)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n===== Helm Values - raw ===========")
	fmt.Println(reflect.TypeOf(coalescedVals.AsMap())) // map[string]interface {}
	fmt.Println(reflect.TypeOf(coalescedVals))         // chartutil.Values
	fmt.Println(coalescedVals)
	fmt.Println("===== Helm Values - Yaml ===========")
	fmt.Println(coalescedVals.YAML())

	fmt.Println("\n====== Helm Templating ==========")
	//e := engine.Engine{Strict: false, LintMode: false}
	var e engine.Engine
	//fmt.Println(e.Render(chart, vals))
	fmt.Println(e.Render(chart, coalescedVals))
}
