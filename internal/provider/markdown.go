package provider

import "strconv"

func markdownValueListInt(values []int) string {
	switch {
	case len(values) == 0:
		return ""
	case len(values) == 1:
		return "`" + strconv.Itoa(values[0]) + "`"
	default:
		s := ""
		for i := 0; i < len(values)-1; i++ {
			s += "`" + strconv.Itoa(values[i]) + "`, "
		}
		s += " and `" + strconv.Itoa(values[len(values)-1]) + "`"
		return s
	}
}
