package helper

import (
	"fmt"
	"math"
	"strings"
)

func MatchStringInStringMap(input string, strs map[int]string) (int, error) {
	keyMap := map[int]int{}
	strArr := make([]string, len(strs))
	i := 0
	for key, str := range strs {
		keyMap[i] = key
		strArr[i] = str
		i += 1
	}
	key, err := MatchStringInStrings(input, strArr)
	return keyMap[key], err
}

func MatchStringInStrings(input string, strs []string) (int, error) {
	in := strings.ToLower(input)
	lowerStrs := make([]string, len(strs))
	for i, s := range strs {
		lowerStrs[i] = strings.ToLower(strings.TrimSpace(s))
	}
	// Check for exact matches
	for i, s := range lowerStrs {
		if in == s {
			return i, nil
		}
	}
	// Check for unique partial matches
	skipped := map[int]bool{}
	for i, b := range []byte(in) {
		found := 0
		foundStr := 0
		for s, str := range lowerStrs {
			if skipped[s] || i >= len(str) || b != str[i] {
				skipped[s] = true
				continue
			}
			found += 1
			foundStr = s
		}
		switch found {
		case 0:
			break
		case 1:
			return foundStr, nil
		}
	}
	return 0, fmt.Errorf("could not match '%s' uniquely to anything in (%s)",
		input, strings.Join(strs, ", "))
}

func StringInStrings(input string, strs []string) (int, error) {
	for i, s := range strs {
		if input == s {
			return i, nil
		}
	}
	return 0, fmt.Errorf("could not find '%s' in (%s)",
		input, strings.Join(strs, ", "))
}

var scaleStrs = []string{
	"thousand",
	"million",
	"billion",
	"trillion",
	"quadrillion",
	"quintillion",
	"sextillion",
	"septillion",
	"octillion",
	"nonillion",
	"decillion",
	"undecillion",
	"duodecillion",
	"tredecillion",
	"quattuordecillion",
	"quindecillion",
	"sexdecillion",
	"septendecillion",
	"octodecillion",
	"novemdecillion",
	"vigintillion",
}

func numberScaleStr(n int) (quant string, pre, rem int) {
	if n < 1000 {
		panic("must be above 1000")
	}
	digits := int(math.Log10(float64(n)))
	scale := digits/3 - 1
	scaleDigits := (scale + 1) * 3
	if l := len(scaleStrs); scale >= l {
		scale = l - 1
	}
	quant = scaleStrs[scale]
	pow := int(math.Pow10(scaleDigits))
	pre = n / pow
	rem = n - pre*pow
	return
}

func numberSubScaleStr(n int) string {
	if n <= 0 || n > 1000 {
		panic("must be between 1 and 999")
	}
	parts := []string{}
	if n >= 100 {
		h := n / 100
		n -= h * 100
		parts = append(parts, numberSubScaleStr(h), "hundred")
		if n > 0 {
			parts = append(parts, "and")
		}
	}
	if n >= 20 {
		t := n / 10
		n -= t * 10
		tStr := ""
		switch t {
		case 2:
			tStr = "twenty"
		case 3:
			tStr = "thirty"
		case 4:
			tStr = "fourty"
		case 5:
			tStr = "fifty"
		case 6:
			tStr = "sixty"
		case 7:
			tStr = "seventy"
		case 8:
			tStr = "eighty"
		case 9:
			tStr = "ninety"
		}
		parts = append(parts, tStr)
	}
	if n > 0 {
		var s string
		switch n {
		case 1:
			s = "one"
		case 2:
			s = "two"
		case 3:
			s = "three"
		case 4:
			s = "four"
		case 5:
			s = "five"
		case 6:
			s = "six"
		case 7:
			s = "seven"
		case 8:
			s = "eight"
		case 9:
			s = "nine"
		case 10:
			s = "ten"
		case 11:
			s = "eleven"
		case 12:
			s = "twelve"
		case 13:
			s = "thirteen"
		case 14:
			s = "fourteen"
		case 15:
			s = "fifteen"
		case 16:
			s = "sixteen"
		case 17:
			s = "seventeen"
		case 18:
			s = "eighteen"
		case 19:
			s = "nineteen"
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, " ")
}

func NumberStr(n int) string {
	if n == 0 {
		return "zero"
	} else if n < 0 {
		return "negative " + NumberStr(-n)
	}
	parts := []string{}
	for n > 1000 {
		quant, pre, rem := numberScaleStr(n)
		parts = append(parts, NumberStr(pre), quant)
		n = rem
	}
	if n > 0 {
		parts = append(parts, numberSubScaleStr(n))
	}
	return strings.Join(parts, " ")
}
