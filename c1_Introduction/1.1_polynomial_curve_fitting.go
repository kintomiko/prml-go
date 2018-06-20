package c1_Introduction

import (
	"errors"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

var twoPi = math.Pi * 2

type Seeder func(float64) float64

type Noiser func() float64

type Sampler func() (float64, error)

var sin = func(x float64) float64 {
	return math.Sin(twoPi * x)
}

var UniformSampler = func() (float64, error) {
	return rand.Float64(), nil
}

var BoxMullerGaussianRandomNumberGenerator = func(expected float64, variance float64) func() float64 {
	flag := 0
	x := 0.0
	y := 0.0
	return func() float64 {
		if flag == 0 {
			u, _ := UniformSampler()
			v, _ := UniformSampler()
			x = math.Sqrt((-2)*math.Log(u)) * math.Cos(twoPi*v)
			y = math.Sqrt((-2)*math.Log(u)) * math.Sin(twoPi*v)
			flag = 1
			return x*variance + expected
		} else {
			flag = 0
			return y*variance + expected
		}
	}
}

type Generator func() ([]float64, error)

var GeneratorImpl = func(sampler Sampler, seeder Seeder, noiser Noiser) func() ([]float64, error) {
	return func() ([]float64, error) {
		input, err := sampler()
		if err == nil {
			return []float64{input, seeder(input) + noiser()}, nil
		}
		return []float64{-1, -1}, errors.New("reached end of input")
	}
}

func fitting(gen Generator, sampleCount int, order int) []float64 {
	//init
	M := mat.NewDense(order, order, nil)
	C := mat.NewVecDense(order, nil)

	//calculate m and c in m * w = c
	for s := 0; s < sampleCount; s++ {
		trainDate, _ := gen()
		for i := 0; i < order; i++ {
			for j := 0; j < order; j++ {
				M.Set(i, j, M.At(i, j)+math.Pow(trainDate[0], float64(i+j)))
			}
			C.SetVec(i, C.AtVec(i)+trainDate[1]*math.Pow(trainDate[0], float64(i)))
		}
	}

	W := mat.NewVecDense(order, nil)
	detM := mat.Det(M)

	for i := 0; i < order; i++ {
		I := mat.DenseCopyOf(M)
		I.SetCol(i, C.RawVector().Data)
		detI := mat.Det(I)
		W.SetVec(i, detI/detM)
	}

	return W.RawVector().Data
}
