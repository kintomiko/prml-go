package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

var two_pi = math.Pi * 2

type Seeder func(float64) float64

type Noiser func() float64

type Sampler func() (float64, error)

var sin = func(x float64) float64 {
	return math.Sin(two_pi * x)
}

var uniformSampler = func() (float64, error) {
	return rand.Float64(), nil
}

var gaussianNoiser = func(expected float64, variance float64) func() float64 {
	flag := 0
	x := 0.0
	y := 0.0
	return func() float64 {
		if flag == 0 {
			u, _ := uniformSampler()
			v, _ := uniformSampler()
			x = math.Sqrt((-2)*math.Log(u)) * math.Cos(two_pi*v)
			y = math.Sqrt((-2)*math.Log(u)) * math.Sin(two_pi*v)
			flag = 1
			return x*variance + expected
		} else {
			flag = 0
			return y*variance + expected
		}
	}
}

var generator = func(sampler Sampler, seeder Seeder, noiser Noiser) func() ([]float64, error) {
	return func() ([]float64, error) {
		input, err := sampler()
		if err == nil {
			return []float64{input, seeder(input) + noiser()}, nil
		}
		return []float64{-1, -1}, errors.New("reached end of input")
	}
}

func main() {
	gen := generator(uniformSampler, sin, gaussianNoiser(0, 0.2))

	for i := 0; i < 1000; i++ {
		pair, _ := gen()
		fmt.Printf("generated x: %v, t: %v\n", pair[0], pair[1])
	}
}
