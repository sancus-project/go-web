package context

//go:generate ./params_type.sh :interface{} Int String

// Get slice of interface{} parameters from RouteContext
func (rctx *RoutingContext) GetSlice(key string) ([]interface{}, bool) {
	if v, ok := rctx.Get(key); ok {
		switch w := v.(type) {
		case []interface{}:
			return w, true
		case interface{}:
			return []interface{}{w}, true
		}
	}

	return nil, false
}

// Get int parameter from RouteContext
func (rctx *RoutingContext) GetInt(key string) (int, bool) {
	var zero int

	if v, ok := rctx.Get(key); ok {
		if w, ok := v.(int); ok {
			return w, true
		}
	}

	return zero, false
}

// Get slice of int parameters from RouteContext
func (rctx *RoutingContext) GetIntSlice(key string) ([]int, bool) {
	if v, ok := rctx.Get(key); ok {
		switch w := v.(type) {
		case []int:
			return w, true
		case int:
			return []int{w}, true
		}
	}

	return nil, false
}

// Get string parameter from RouteContext
func (rctx *RoutingContext) GetString(key string) (string, bool) {
	var zero string

	if v, ok := rctx.Get(key); ok {
		if w, ok := v.(string); ok {
			return w, true
		}
	}

	return zero, false
}

// Get slice of string parameters from RouteContext
func (rctx *RoutingContext) GetStringSlice(key string) ([]string, bool) {
	if v, ok := rctx.Get(key); ok {
		switch w := v.(type) {
		case []string:
			return w, true
		case string:
			return []string{w}, true
		}
	}

	return nil, false
}
