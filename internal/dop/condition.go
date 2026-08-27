package dop

import "math"

func InfinityNorm(m Mat4) float64 {
	max := 0.0
	for i := 0; i < 4; i++ {
		sum := 0.0
		for j := 0; j < 4; j++ {
			sum += math.Abs(m.Data[i][j])
		}
		if sum > max {
			max = sum
		}
	}
	return max
}

func ConditionNumber(n, inv Mat4) float64 {
	return InfinityNorm(n) * InfinityNorm(inv)
}

func LogCondition(cond float64) float64 {
	return math.Log10(cond)
}
