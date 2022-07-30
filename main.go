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

	// Prepare vals
	var vals map[string]interface{}
	// To override values
	// vals := map[string]interface{}{
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
	coalescedVals, err := chartutil.CoalesceValues(chart, vals)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n===== Helm Values - raw ===========")
	fmt.Println(reflect.TypeOf(coalescedVals))
	fmt.Println(coalescedVals)
	fmt.Println("\n===== Helm Values - in Yaml ===========")
	fmt.Println(coalescedVals.YAML())
	fmt.Println("\n====== Helm Templating ==========")
	// func (e Engine) Render(chrt *chart.Chart, values chartutil.Values) (map[string]string, error)
	fmt.Println(engine.Render(chart, coalescedVals))

}
