package egothor

// DiffApply returns string
func DiffApply(orig string, diff string) string {
	// Check for word
	if orig == "" || diff == "" {
		return orig
	}

	str := []rune(orig)
	pos := len(str) - 1
	// pos in string is out of bound
	if pos < 0 {
		return orig
	}
	runeDiff := []rune(diff)

	for i := 0; i < len(runeDiff)/2; i++ {
		cmd := runeDiff[2*i]
		param := runeDiff[2*i+1]
		parNum := param - 'a' + 1
		switch cmd {
		case '-':
			pos = pos - int(parNum) + 1
			break
		case 'R':

			if pos >= 0 {
				str[pos] = rune(param)
				break
			} else {
				// pos in string is out of bound
				return orig
			}
		case 'D':
			o := int(pos + 1)
			pos -= int(parNum) - 1
			if pos >= 0 && o >= 0 {
				str = append(str[:pos], str[o:]...)
				break
			} else  {
				// pos in string is out of bound
				return orig
			}
		case 'I':
			pos++
			if pos < 0 {
				// pos in string is out of bound
				return orig
			}
			str = append(str, 0)
			copy(str[pos+1:], str[pos:])
			str[pos] = rune(param)
			break
		}
		pos--
	}
	return string(str)
}
