package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

type tree []node

type node struct {
	label    string
	children []int // indexes into tree
}

func vis(t tree) {
	if len(t) == 0 {
		fmt.Println("<empty>")
		return
	}
	var f func(int, string)
	f = func(n int, pre string) {
		ch := t[n].children
		if len(ch) == 0 {
			fmt.Println("╴", t[n].label)
			return
		}
		fmt.Println("┐", t[n].label)
		last := len(ch) - 1
		for _, ch := range ch[:last] {
			fmt.Print(pre, "├─")
			f(ch, pre+"│ ")
		}
		fmt.Print(pre, "└─")
		f(ch[last], pre+"  ")
	}
	f(0, "")
}

func main() {

	log.SetFlags(0) // no timestamp
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) != 1 {
		log.Fatalf("supply a chart file or directory")
	}
	chartPath := os.Args[1]

	fmt.Println("\nLoading Helm Chart...")
	loadedChart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.TypeOf(loadedChart))
	// fmt.Print(*loadedChart)

	fmt.Println("\nPopulating Helm Values...")
	// var vals chartutil.Values

	// vals := chartutil.Values{
	// 	"replicaCount": 3,
	// }

	// Signature: func CoalesceValues(chrt *chart.Chart, vals map[string]interface{}) (Values, error)
	// throws nil pointer evaluating interface {}
	// vals, err := chartutil.CoalesceValues(loadedChart, map[string]interface{}{})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(vals.YAML())

	releaseOptions := chartutil.ReleaseOptions{Name: "release1", Namespace: "ns1"}
	// Submitting empty map param {}{}
	vals, err := chartutil.ToRenderValues(loadedChart, map[string]interface{}{},
		releaseOptions, chartutil.DefaultCapabilities)
	// vals, err := chartutil.ToRenderValues(loadedChart, map[string]interface{}{},
	// 	chartutil.ReleaseOptions{}, chartutil.DefaultCapabilities)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// fmt.Println(vals.YAML())

	fmt.Println("\nHelm Templating...")

	// Alternative to engine.Render function. Using Render Method outputs trailing nil.
	// e := engine.Engine{Strict: false, LintMode: false}
	// fmt.Println(e.Render(loadedChart, vals))

	// Templated Chart represented by "m" (map[string]string)
	// where keys are the filenames and values are the file contents
	m, err := engine.Render(loadedChart, vals)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(m)

	fmt.Println("\nChart files found:")
	// Populate keys (filenames)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
		fmt.Println(k)
	}

	fmt.Println("\nSearching images in K8S manifests...")
	// Populate keys (filenames) with "image:" in the file content
	var imageKeys []string
	for _, k := range keys {
		if strings.Contains(m[k], "image:") {
			imageKeys = append(imageKeys, k)
			//fmt.Println(m[imageKeys])

			re := regexp.MustCompile(`image:.+`)
			imageLines := re.FindAllString(m[k], -1)

			if len(imageLines) != 1 {
				fmt.Printf("\nImage found in %s...\n", k)
			}
			// fmt.Println(imageLines)
			for _, i := range imageLines {
				image := strings.TrimPrefix(i, "image:")
				image = strings.TrimSpace(image)
				image = strings.Trim(image, "\"")
				fmt.Println(image)
			}
		}
	}

	fmt.Printf("\nBuilding Tree for the Helm Chart Tree: \"%s\"...\n", loadedChart.Name())

	// Closure must be declared to allow recursions later on
	var depRecursion func(myChart chart.Chart, nodeID int) tree

	// allNodeIDs initialized already to reserve 0 for root.
	// Appending always dummy value "node". Slice keys act as Node IDs. Length represents Node count.
	allNodeIDs := []string{"node"} // 0: node, 1: node,...
	fullTree := tree{{label: loadedChart.Name(), children: []int{}}}
	var currentDepsNodeIDs []int

	depRecursion = func(myChart chart.Chart, nodeID int) tree {
		parent := myChart.Name()
		chartDeps := myChart.Dependencies()

		currentDepsNodeIDs = nil

		fmt.Printf("\n=== Parent chart: %s contains %d dependencies. === \n", parent, len(chartDeps))
		fmt.Println("Tree state:", fullTree)

		// Chart does not have further deps
		if len(chartDeps) == 0 {
			fmt.Println("No dependencies found. Returning...")
			return fullTree
		} else {
			// root Node already declared, len == 1
			shift := len(allNodeIDs)
			for i, dep := range chartDeps {
				// shifting currentDepsNodeIDs overcomes zero-based range indexing
				currentDepsNodeIDs = append(currentDepsNodeIDs, shift+i) // [1,2,3,4], next parent: [5,6,7]...

				// allNodeIDs grow with every new dependencies. Slice keys represent Node IDs. Length represents Node count.
				allNodeIDs = append(allNodeIDs, "node")

				fmt.Printf("Adding \"%s\" Node to the Tree. Current Node count: %d \n", dep.Name(), len(allNodeIDs))
				fullTree = append(fullTree, node{label: dep.Name(), children: []int{}})
			}

			fmt.Printf("New Tree state: %v \n", fullTree)
			fullTree[nodeID] = node{label: parent, children: currentDepsNodeIDs} // NodeID initially passed to the function
			fmt.Printf("Children added to the Node in the Tree: %v \n", fullTree)

			for i, dep := range chartDeps {
				fmt.Printf("Recursive search for: %s\n", dep.Name())
				//fmt.Println(shift + i)
				go depRecursion(*dep, shift+i)
				time.Sleep(time.Second)
			}
		}
		return fullTree
	}

	depRecursion(*loadedChart, 0)

	fmt.Println("\n=== Helm Tree: ===\n")
	vis(fullTree)

}
