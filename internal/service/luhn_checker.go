package service

import "strconv"

type LuhnChecker struct {
}

func (lc LuhnChecker) Check(order string) bool {
	lunh := 0
	for i := 0; i < len(order); i++ {
		k := len(order) - i - 1
		s := order[k : k+1]
		d, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if (i+1)%2 == 0 {
			if d >= 5 {
				lunh += 2*d - 9
			} else {
				lunh += 2 * d
			}
		} else {
			lunh += d
		}
	}
	if lunh%10 == 0 {
		return true
	} else {
		return false
	}
}
