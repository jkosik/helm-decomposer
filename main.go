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

	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n===== Helm Chart ===========")
	fmt.Print(*chart)
	engine := engine.Engine{Strict: false, LintMode: false}

	// manual values
	var extraVals map[string]interface{}
	// extraVals := map[string]interface{}{
	// 	"replicaCount": 3,
	// }

	// Readvalues
	// var myData []byte
	// d, err := chartutil.ReadValues(myData)
	// if err != nil {
	// 	panic(err)
	// }

	// Coalesce values
	// func CoalesceValues(chrt *chart.Chart, vals map[string]interface{}) (Values, error)
	vals, err := chartutil.CoalesceValues(chart, extraVals)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n===== Helm Values - raw ===========")
	fmt.Println(reflect.TypeOf(vals))
	fmt.Println(vals)
	fmt.Println("\n===== Helm Values - in Yaml ===========")
	fmt.Println(vals.YAML())
	fmt.Println("\n====== Helm Templating ==========")
	// func (e Engine) Render(chrt *chart.Chart, values chartutil.Values) (map[string]string, error)
	fmt.Println(engine.Render(chart, vals))

}
