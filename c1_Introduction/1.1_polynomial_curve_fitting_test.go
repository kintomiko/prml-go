package c1_Introduction

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const float64EqualityThreshold = 1e-2
const SAMPLING_NUMBER = 100000

func TestGenerator(t *testing.T) {
	gen := GeneratorImpl(UniformSampler, sin, BoxMullerGaussianRandomNumberGenerator(0, 0.2))

	for i := 0; i < SAMPLING_NUMBER; i++ {
		_, err := gen()
		if err != nil {
			t.Error("failed to generate new random variable")
		}
	}
}

func TestGaussianRandomGenerator(t *testing.T) {
	ONE := 1.0
	MINUS_TWO := -2.0

	cases := []struct {
		expect   float64
		variance float64
		lb       *float64
		hb       *float64
		expected float64
	}{
		{0, 1, &ONE, nil, 0.1587},
		{0, 1, &MINUS_TWO, &ONE, 0.8186},
		{0, 1, nil, &MINUS_TWO, 0.0228},
		{2.4, 1.8, &ONE, nil, 0.7816},
		{2.4, 1.8, &MINUS_TWO, &ONE, 0.2111},
		{2.4, 1.8, nil, &MINUS_TWO, 0.0073},
		{-1.3, 2.8, &ONE, nil, 0.2057},
		{-1.3, 2.8, &MINUS_TWO, &ONE, 0.393},
		{-1.3, 2.8, nil, &MINUS_TWO, 0.4013},
	}

	for _, testcase := range cases {

		gaussian := BoxMullerGaussianRandomNumberGenerator(testcase.expect, testcase.variance)
		count := 0
		for i := 0; i < SAMPLING_NUMBER; i++ {
			randomVariable := gaussian()
			if inCase(randomVariable, testcase.lb, testcase.hb) {
				count++
			}
		}
		possibility := float64(count) / SAMPLING_NUMBER

		if !almostEqual(possibility, testcase.expected) {
			t.Errorf("Possibility incorrect, case: %v, got: %f",
				testcase,
				possibility,
			)
		}
	}
}

func TestFitting(t *testing.T) {
	w := fitting(
		GeneratorImpl(UniformSampler, sin, BoxMullerGaussianRandomNumberGenerator(0, 2)),
		1000,
		8,
	)

	fmt.Printf("W: %v", w)
	graphPointNum := 100
	input := make([][]float64, graphPointNum)
	for i := 0; i < graphPointNum; i++ {
		x := float64(i) / float64(graphPointNum)
		input[i] = []float64{x, y(x, w)}
	}

	expected := make([][]float64, graphPointNum)
	for i := 0; i < graphPointNum; i++ {
		x := float64(i) / float64(graphPointNum)
		expected[i] = []float64{x, sin(x)}
	}
	plotLine(input, expected)
}

func y(x float64, w []float64) float64 {
	result := 0.0
	for i := 0; i < len(w); i++ {
		result += w[i] * math.Pow(x, float64(i))
	}
	return result
}

func plotLine(actual, expected [][]float64) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Polynomial Curve Fitting"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	actualPts := make(plotter.XYs, len(actual))
	for i, data := range actual {
		actualPts[i].X = data[0]
		actualPts[i].Y = data[1]
	}

	expectedPts := make(plotter.XYs, len(expected))
	for i, data := range expected {
		expectedPts[i].X = data[0]
		expectedPts[i].Y = data[1]
	}

	err = plotutil.AddLinePoints(p,
		"actual", actualPts,
		"expected", expectedPts,
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func inCase(rv float64, lb, hb *float64) bool {
	if lb != nil && rv < *lb {
		return false
	}
	if hb != nil && rv >= *hb {
		return false
	}
	return true
}
