package context

func (rctx *RoutingContext) Path() string {
	var s string

	prefix := rctx.RoutePrefix
	path := rctx.RoutePath

	if prefix == "/" {
		s = path
	} else if path == "" {
		s = prefix
	} else {
		s = prefix + path
	}

	return s
}
