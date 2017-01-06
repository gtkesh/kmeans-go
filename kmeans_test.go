package main

import (
	"fmt"
	"testing"
)

func TestEuclideanDistance(t *testing.T) {
	tests := []struct {
		n1               Node
		n2               Node
		expectedDistance float64
	}{
		{
			Node{Lat: 1.1, Lng: 2.2},
			Node{Lat: 3.3, Lng: 4.4},
			3.111269837220809,
		},
		{
			Node{Lat: 40.9027901845628, Lng: -73.0782184632683},
			Node{Lat: 40.85770758108904, Lng: -73.08187007904053},
			0.04523024910079703,
		},
	}

	for _, tt := range tests {
		dist := distance(tt.n1, tt.n2)
		if dist != tt.expectedDistance {
			t.Error(
				"expected", tt.expectedDistance,
				"got", dist,
			)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		in  []byte
		out []Node
	}{
		{
			[]byte(`[{"coordinates": {"lat": 40.9027901845628, "lng": -73.0782184632683}},
                {"coordinates": {"lat": 40.85770758108904, "lng": -73.08187007904053}}]`),
			[]Node{
				Node{Lat: 40.9027901845628, Lng: -73.0782184632683},
				Node{Lat: 40.85770758108904, Lng: -73.08187007904053},
			},
		},
	}

	for _, tt := range tests {
		nodes, err := parse(tt.in)
		if err != nil {
			fmt.Println(err)
			t.Error("expected no error")
		}
		for i := 0; i < len(tt.out); i++ {
			if !equals(nodes[i], tt.out[i]) {
				t.Error(
					"expected", tt.out,
					"got", nodes,
				)
			}
		}
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		nodes            []Node
		expectedMeanNode Node
	}{
		{[]Node{Node{0.0, 1.1}}, Node{0.0, 1.1}},
		{[]Node{Node{2.0, 3.0}, Node{4.0, 5.0}}, Node{3.0, 4.0}},
	}
	for _, tt := range tests {
		node := mean(tt.nodes)
		if !equals(node, tt.expectedMeanNode) {
			t.Error(
				"expected", tt.expectedMeanNode,
				"got", node,
			)
		}
	}
}

func TestClosest(t *testing.T) {
	tests := []struct {
		nodes         []Node
		centroid      Node
		expectedIndex int
	}{
		{[]Node{Node{0.0, 1.1}}, Node{0.0, 1.1}, 0},
		{[]Node{Node{0.0, 1.1}, Node{22.2, 10.123}}, Node{0.0, 1.1}, 0},
		{[]Node{Node{0.0, 123.123}, Node{22.2, 10.123}, Node{5555.2, 10.123}}, Node{0.0, 1.1}, 1},
	}

	for _, tt := range tests {
		index := closest(tt.centroid, tt.nodes)
		if index != tt.expectedIndex {
			t.Error(
				"expected", tt.expectedIndex,
				"got", index,
			)
		}
	}
}
