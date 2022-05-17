package context

// Clone makes a copy of the the RouteParams table
func CloneRouteParams(src map[string]interface{}) map[string]interface{} {
	if l := len(src); l > 0 {
		m := make(map[string]interface{})
		for k, v := range src {
			switch vv := v.(type) {
			case []int:
				ww := make([]int, len(vv))
				for i, w := range vv {
					ww[i] = w
				}
				v = ww
			case []string:
				ww := make([]string, len(vv))
				for i, w := range vv {
					ww[i] = w
				}
				v = ww
			case []interface{}:
				ww := make([]interface{}, len(vv))
				for i, w := range vv {
					ww[i] = w
				}
				v = ww
			}

			m[k] = v
		}
		return m
	}
	return nil
}

// Sets value for a RouteParam
func (rctx *RoutingContext) Set(key string, v interface{}) {
	if rctx.RouteParams == nil {
		rctx.RouteParams = make(map[string]interface{}, 1)
	}

	rctx.RouteParams[key] = v
}

// Append value for a RouteParam
func (rctx *RoutingContext) Add(key string, v interface{}) {
	var vv []interface{}

	if rctx.RouteParams == nil {
		rctx.RouteParams = make(map[string]interface{}, 1)
	}

	// extract previous slice if available
	// and convert single item entries to slices
	if w, ok := rctx.RouteParams[key]; ok {
		switch ww := w.(type) {
		case []interface{}:
			vv = ww
		case interface{}, nil:
			vv = append(vv, w)
		}
	}

	// and append our value
	vv = append(vv, v)
	rctx.RouteParams[key] = vv
}

// Get item from RouteParams
func (rctx *RoutingContext) Get(key string) (interface{}, error, bool) {
	if rctx.RouteParams != nil {
		if w, ok := rctx.RouteParams[key]; ok {
			return w, nil, ok
		}
	}

	return nil, nil, false
}
