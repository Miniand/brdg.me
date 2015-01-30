package helper

import (
	"fmt"
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
