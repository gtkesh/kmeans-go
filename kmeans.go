package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
)

type Transaction struct {
	LatLng struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"coordinates"`
}

type Node struct {
	Lat float64
	Lng float64
}

// Takes the list of nodes as observations. Returns k centroids.
// Terminates when we've converged or reach max number of iterations.
func KMeans(nodes []Node, k int, maxIters int) []Node {
	centroids := initCentroids(k, nodes)
	var numIters int
	for {
		// Step 1. Assign nodes to k clusters.
		clusters := make(map[int][]Node)
		for _, node := range nodes {
			i := closest(node, centroids)
			clusters[i] = append(clusters[i], node)
		}

		// Step 2. Update means.
		for k, cluster := range clusters {
			meanNode := mean(cluster)
			if equals(centroids[k], meanNode) {
				// We've converged, so we stop.
				return centroids
			}
			centroids[k] = meanNode
		}
		// If we've reached max num of iterations, we stop.
		if numIters == maxIters {
			return centroids
		}
		numIters++
	}
	return centroids
}

// Compares lat-long pairs of two nodes.
func equals(n1, n2 Node) bool {
	if n1.Lat != n2.Lat || n1.Lng != n2.Lng {
		return false
	}
	return true
}

// Returns the index of the closest node to centroid from a list of nodes.
func closest(centroid Node, nodes []Node) int {
	distances := make([]float64, len(nodes))

	for i, node := range nodes {
		distances[i] = distance(centroid, node)
	}

	current := distances[0]
	var index int
	for i, d := range distances {
		if d < current {
			current = d
			index = i
		}
	}

	return index
}

// Returns the mean (possibly a new) Node from the list of Nodes.
func mean(nodes []Node) Node {
	var latSum, lngSum float64
	for _, node := range nodes {
		latSum += node.Lat
		lngSum += node.Lng
	}

	return Node{
		Lat: latSum / float64(len(nodes)),
		Lng: lngSum / float64(len(nodes)),
	}
}

// Randomly picks and returns initial K centroid nodes.
func initCentroids(k int, nodes []Node) []Node {
	centroids := make([]Node, 0)

	for i := 0; i < k; i++ {
		j := rand.Intn(k)
		centroids = append(centroids, nodes[j])
	}
	return centroids
}

// Returns eucledian distance between two nodes.
func distance(n1, n2 Node) float64 {
	d := math.Hypot((n1.Lat - n2.Lat), (n1.Lng - n2.Lng))
	return d
}

// Parses json data into list of Nodes.
func parse(byt []byte) ([]Node, error) {
	var transactions []Transaction
	err := json.Unmarshal(byt, &transactions)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0)
	for _, t := range transactions {
		nodes = append(nodes, Node{t.LatLng.Lat, t.LatLng.Lng})
	}

	return nodes, nil
}

// Reads file and returns the byte stream.
func read(path string) []byte {
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return byt
}

func main() {
	byt := read("transaction.json")
	nodes, err := parse(byt)
	if err != nil {
		log.Fatal(err)
	}

	k := 3
	maxIterations := 1000000
	centroids := KMeans(nodes, k, maxIterations)
	fmt.Println(centroids)
}
