package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func main() {

	inputChart := flag.String("chart", "sample-helm-charts/nginx", "Helm Chart to process. Submit tar.gz or folder name.")
	outputFile := flag.Bool("o", false, "Write output to helm-decomposer-output.md. (default \"false\")")
	processImages := flag.Bool("i", false, "Inspect images used in the Helm Chart. (default \"false\")")

	flag.Parse()

	fmt.Println(*inputChart)
	fmt.Println(*outputFile)
	fmt.Println(*processImages)

	fmt.Printf("\nLoading Helm Chart \"%s\"...\n", *inputChart)
	loadedChart, err := loader.Load(*inputChart)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nPopulating Helm Values...")

	// colVals, err := chartutil.CoalesceValues(loadedChart, map[string]interface{}{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	releaseOptions := chartutil.ReleaseOptions{Name: "release1", Namespace: "ns1"}
	// Submitting empty map param {}{}
	vals, err := chartutil.ToRenderValues(loadedChart, map[string]interface{}{},
		releaseOptions, chartutil.DefaultCapabilities)

	// templatedVals, _ := vals.YAML()
	// fmt.Println("Templated Values: \n", templatedVals)

	fmt.Println("\nHelm Templating...")

	// Templated Chart represented by "m" (map[string]string)
	// where keys are the filenames and values are the file contents
	m, err := engine.Render(loadedChart, vals)
	if err != nil {
		log.Println(err)
		fmt.Println("\nWARNING: Helm Chart can not be fully templated. Please check values files on all levels, usage of aliases, etc...")
	}
	fmt.Println("Templated manifests: \n", m)

	// Detect images in the Chart structure
	detectImages(m)

	// Build visual tree of Chart dependencies
	fmt.Printf("\nBuilding Tree for the Helm Chart Tree: \"%s\"...\n", loadedChart.Name())

	// Closure must be declared to allow recursions later on
	var depRecursion func(myChart chart.Chart, nodeID int) tree

	// allNodeIDs initialized already to reserve 0 for root. Needed by vis() in tree.go
	// Slice keys act as Node IDs. Values are always "dummy". Length represents Node count.
	allNodeIDs := []string{"node"} // 0: node, 1: node,...
	fullTree := tree{{label: loadedChart.Name(), children: []int{}}}
	var currentDepsNodeIDs []int

	depRecursion = func(myChart chart.Chart, nodeID int) tree {
		parent := myChart.Name()
		chartDeps := myChart.Dependencies()

		currentDepsNodeIDs = nil

		// fmt.Printf("\n=== Parent chart: %s contains %d dependencies. === \n", parent, len(chartDeps))
		// fmt.Println("Tree state:", fullTree)

		// Chart does not have further deps
		if len(chartDeps) == 0 {
			// fmt.Println("No dependencies found. Continuing...")
		} else {
			// root Node already declared, i.e. len == 1. Child Node IDs are shifted.
			shift := len(allNodeIDs)
			for i, dep := range chartDeps {
				// Composing from scratch slice of child Node IDs for the tested parent.
				// Node ID == allNodeIDs slice KEY IDs.
				// currentDepsNodeIDs's VALUES are shifted +1 to KEYS from the allNodeIDs
				currentDepsNodeIDs = append(currentDepsNodeIDs, shift+i) // [1,2,3,4], for the next parent: [5,6,7]...

				// allNodeIDs grows with every new dependencies. Slice keys represent Node IDs (zero-based). Slice length represents Node count.
				allNodeIDs = append(allNodeIDs, "node")

				// fmt.Printf("New Node \"%s\" (Node ID: %d) added to the Tree. Current Node count: %d \n", dep.Name(), shift+i, len(allNodeIDs))
				fullTree = append(fullTree, node{label: dep.Name(), children: []int{}})
			}

			// fmt.Printf("New Tree state: %v \n", fullTree)
			fullTree[nodeID] = node{label: parent, children: currentDepsNodeIDs} // NodeID initially passed to the function
			// fmt.Printf("Childrens in Tree updated for Node \"%s\" (Node ID %d): %v \n", parent, nodeID, fullTree)

			for i, dep := range chartDeps {
				// fmt.Printf("Recursive search for: \"%s\", Node ID: %d\n", dep.Name(), shift+i)
				depRecursion(*dep, shift+i)
				//time.Sleep(100 * time.Millisecond)
			}
		}
		return fullTree
	}

	depRecursion(*loadedChart, 0)

	fmt.Println("\n=== Helm Tree: ===\n")

	// If output file needed
	if *outputFile {
		f, err := os.Create("helm-decomposer-output.md")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Capturing this at stdout
		vis(fullTree)

		w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		f.Write(out)

		// Print captured
		fmt.Printf("%s", out)

	} else {
		vis(fullTree)
	}

}
