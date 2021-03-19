package types

import (
	"math"
	"strconv"
	"strings"
)

//Round64  round float64 value https://gist.github.com/DavidVaini/10308388
func Round64(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

//Round32 round float64 value https://gist.github.com/DavidVaini/10308388
func Round32(val float32, roundOn float32, places int) (newVal float32) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * float64(val)
	_, div := math.Modf(digit)
	if div >= float64(roundOn) {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = float32(round / pow)
	return
}

//Ftoa covert float64 value to string with default round
func Ftoa(v float64) string {
	//a := strconv.FormatFloat(v, 'f', 2, 64)
	//fmt.Println(fmt.Printf("Ftoa: %.2f %s %.2f", v, a, Round(v, .5, 2)))

	//return a

	return strconv.FormatFloat(Round64(v, .5, 2), 'f', 2, 64)
}

//Atof covert string value to float with default value
func Atof(v string, d float64) float64 {
	f, err := strconv.ParseFloat(strings.Trim(v, "\""), 64)

	if err != nil {
		return d
	}
	return f
}
