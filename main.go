package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"

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

	if len(os.Args[1:]) != 1 {
		log.Fatalf("supply a chart file or directory")
	}
	chartPath := os.Args[1]

	fmt.Println("\nLoading Helm Chart...")
	loadedChart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	//fmt.Println(reflect.TypeOf(loadedChart))

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
		log.Println(err)
		fmt.Println("\nWARNING: Helm Chart can not be fully templated. Please check values files on all levels, usage of aliases, etc...")
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
			fmt.Println("No dependencies found. Continuing...")
		} else {
			// root Node already declared, len == 1
			shift := len(allNodeIDs)
			for i, dep := range chartDeps {
				// Composing from scratch slice of child Node IDs for the tested parent.
				// Node ID == Slice KEY IDs for the zer-based Tree which will be submitted to vis().
				// currentDepsNodeIDs's VALUES are +1 to KEYS from the Tree
				currentDepsNodeIDs = append(currentDepsNodeIDs, shift+i) // [1,2,3,4], for the next parent: [5,6,7]...

				// allNodeIDs grows with every new dependencies. Slice keys represent Node IDs (zero-based). Slice length represents Node count.
				allNodeIDs = append(allNodeIDs, "node")

				fmt.Printf("New Node \"%s\" (Node ID: %d) added to the Tree. Current Node count: %d \n", dep.Name(), shift+i, len(allNodeIDs))
				fullTree = append(fullTree, node{label: dep.Name(), children: []int{}})
			}

			fmt.Printf("New Tree state: %v \n", fullTree)
			fullTree[nodeID] = node{label: parent, children: currentDepsNodeIDs} // NodeID initially passed to the function
			fmt.Printf("Childrens in Tree updated for Node \"%s\" (Node ID %d): %v \n", parent, nodeID, fullTree)

			for i, dep := range chartDeps {
				fmt.Printf("Recursive search for: \"%s\", Node ID: %d\n", dep.Name(), shift+i)

				depRecursion(*dep, shift+i)
				//time.Sleep(100 * time.Millisecond)
			}
		}
		return fullTree
	}

	depRecursion(*loadedChart, 0)

	// Wait until only parent program is running
	for runtime.NumGoroutine() > 1 {
		//fmt.Printf("\Runnint Go routines count: %d ", runtime.NumGoroutine())
	}

	fmt.Println("\n=== Helm Tree: ===\n")
	vis(fullTree)

}
