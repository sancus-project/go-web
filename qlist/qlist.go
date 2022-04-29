package qlist

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.sancus.dev/web/errors"
)

const (
	MinimumQuality = 0.
	MaximumQuality = 1.
	Epsilon        = 0.001
)

type QualityValue struct {
	Value   string
	Quality float32
}

func (q QualityValue) String() string {
	if q.Quality+Epsilon > MaximumQuality {
		return q.Value
	} else {
		return fmt.Sprintf("%s;q=%v", q.Value, q.Quality)
	}
}

type QualityList []QualityValue

func (ql QualityList) String() string {
	s := make([]string, len(ql))

	for i, x := range ql {
		s[i] = x.String()
	}

	return strings.Join(s, ", ")
}

func ParseQualityHeader(hdr http.Header, name string) (out QualityList, err error) {
	for k, v := range hdr {
		if strings.EqualFold(name, k) {
			for _, s := range v {
				var q []QualityValue

				q, err = ParseQualityString(s)
				if err != nil {
					return
				}
				out = append(out, q...)
			}
		}
	}
	return
}

func ParseQualityString(qlist string) (out QualityList, err error) {
	for _, s := range strings.Split(qlist, ",") {
		fields := strings.Split(s, ";")

		// remove whitespace
		for i, v := range fields {
			fields[i] = strings.TrimSpace(v)
		}

		if len(fields) > 0 {

			// value
			s := strings.ToLower(fields[0])
			if len(s) == 0 {
				goto invalid
			}

			qv := QualityValue{
				Value:   s,
				Quality: MaximumQuality,
			}

			if len(qv.Value) == 0 {
				goto invalid
			}

			// attributes
			fields = fields[1:]
			for _, s := range fields {
				if strings.HasPrefix(s, "q=") {
					q, err := strconv.ParseFloat(s[2:], 32)
					if err != nil || q < MinimumQuality || q > MaximumQuality {
						goto invalid
					}
					qv.Quality = float32(q)
				}
			}

			out = append(out, qv)
		}
	}
	return
invalid:
	err = errors.ErrInvalidValue("%q", qlist)
	return
}
