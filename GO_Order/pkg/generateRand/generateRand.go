package generaterand

import "math/rand/v2"

func RandStr(length int) string {
	randStr := ""
	i := 0
	for length > i {
		randomer := rand.IntN(123)
		if (randomer > 47 && randomer < 58) || (randomer > 64 && randomer < 91) || (randomer > 96) {
			randStr += string(byte(randomer))
			i++
		}
	}
	return randStr
}
func RandNumberStr(length int) string {
	randStr := ""
	i := 0
	for length > i {
		randomer := rand.IntN(123)
		if randomer > 47 && randomer < 58 {
			randStr += string(byte(randomer))
			i++
		}
	}
	return randStr
}
