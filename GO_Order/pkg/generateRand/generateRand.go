package generaterand

import "math/rand/v2"

func RandStr(lenght int) string {
	randStr := ""
	i := 0
	for lenght > i {
		randomer := rand.IntN(123)
		if (randomer > 47 && randomer < 58) || (randomer > 64 && randomer < 91) || (randomer > 96) {
			randStr += string(byte(randomer))
			i++
		}
	}
	return randStr
}
