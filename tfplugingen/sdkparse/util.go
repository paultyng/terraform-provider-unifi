package sdkparse

func coallesce(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

func uniq(s ...string) []string {
	n := []string{}
	for _, v := range s {
		if v != "" {
			n = append(n, v)
		}
	}
	return n
}
