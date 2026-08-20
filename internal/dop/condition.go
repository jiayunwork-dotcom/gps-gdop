package dop

import "math"

// InfinityNorm returns the maximum absolute row sum of a matrix, the
// operator norm used for the condition number estimate.
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

// ConditionNumber estimates the conditioning of the normal matrix as the
// product of the infinity norms of the matrix and its inverse. A value
// near unity means the constellation geometry is robust; a large value
// flags a nearly singular geometry even when the inversion still succeeds.
func ConditionNumber(n, inv Mat4) float64 {
	return InfinityNorm(n) * InfinityNorm(inv)
}

// LogCondition scales the condition number for display.
func LogCondition(cond float64) float64 {
	return math.Log10(cond)
}
