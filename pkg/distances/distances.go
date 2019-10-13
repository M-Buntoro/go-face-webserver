package distances

import (
	"math"
)

func CosineSimilarity(a, b []float32) float32 {
	var dot_product, pow_a, pow_b float32

	for index := 0; index < len(a); index++ {
		dot_product += a[index] * b[index]
		pow_a += a[index] * a[index]
		pow_b += b[index] * b[index]
	}
	pow_a = float32(math.Sqrt(float64(pow_a)))
	pow_b = float32(math.Sqrt(float64(pow_b)))

	sim := dot_product / (pow_a * pow_b)
	if sim > 1.0 {
		sim = 1.0
	} else if sim < 0.0 {
		sim = 0.0
	}

	return sim
}

func EuclidianDistance(a, b []float32) float32 {
	var sum float64 = 0.0

	for index := 0; index < len(a); index++ {
		diff := math.Abs(float64(a[index] - b[index]))
		sum += diff * diff
	}

	sim := float32(math.Sqrt(sum))
	if sim > 1.0 {
		sim = 1.0
	} else if sim < 0.0 {
		sim = 0.0
	}

	return sim
}
