package indicators

type Operator func(x, y float64) bool
type ArrayOperator func(x, y []float64) bool

func CrossOver(value1, value2 []float64) bool {
	return value1[1] < value2[1] && value1[0] > value2[0]
}

func CrossUnder(value1, value2 []float64) bool {
	return value1[1] > value2[1] && value1[0] < value2[0]
}

func Higher(value1, value2 float64) bool {
	return value1 > value2
}

func Lower(value1, value2 float64) bool {
	return value1 < value2
}

func Equal(value1, value2 float64) bool {
	return value1 < value2
}
