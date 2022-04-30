package qlist

func FindQuality(s string, ql QualityList) (float32, bool) {
	for _, qv := range ql {
		switch qv.Value {
		case s, "*":
			// match
			return qv.Quality, true
		}
	}
	// no match
	return 0.0, false
}

func BestEncodingQuality(supported []string, ql QualityList) (string, bool) {
	if len(ql) != 0 {
		// we have Accept-Encoding data
		bestquality := float32(0.0)
		bestencoding := ""

		// pick the best supported match
		for _, encoding := range supported {
			quality, _ := FindQuality(encoding, ql)
			if quality > bestquality {
				bestquality = quality
				bestencoding = encoding
			}
		}

		// but also test for identity
		quality, ok := FindQuality("identity", ql)
		if quality > bestquality {
			// identity is best
			goto identity
		} else if bestencoding == "" {
			// nothing chosen
			if !ok {
				// but identity wasn't forbidden
				goto identity
			}
			// 406
			return "", false
		} else {
			// we have a better option than identity
			return bestencoding, true
		}
	}
identity:
	// nothing chosen and identity not forbidden, pick it
	return "identity", true
}

func BestEncoding(supported []string, header string) (string, bool) {
	// bad header is the same as no header. empty list
	ql, _ := ParseQualityString(header)
	return BestEncodingQuality(supported, ql)
}
