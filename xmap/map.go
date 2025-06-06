package xmap

func ToMap[FK comparable, FV any, TK comparable, TV any](from map[FK]FV, cb func(FK, FV) (TK, TV)) map[TK]TV {
	if from == nil {
		return nil
	}
	to := make(map[TK]TV, len(from))
	for k, v := range from {
		tk, tv := cb(k, v)
		to[tk] = tv
	}
	return to
}
