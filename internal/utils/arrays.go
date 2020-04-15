package utils

func InArrayS(needle string, stack []string) bool {
	return InArraySF(needle, stack, func(need string, stackVal string) bool {
		return need == stackVal
	})
}

func InArraySF(needle string, stack []string, handler func(need string, stackVal string) bool) bool {
	for _, v := range stack {
		if handler(needle, v) {
			return true
		}
	}
	return false
}

func ArrayDiffS(src []string, newer []string) (added []string, removed []string) {
	for _, oldVal := range src {
		if !InArrayS(oldVal, newer) {
			removed = append(removed, oldVal)
		}
	}
	for _, newVal := range newer {
		if !InArrayS(newVal, src) {
			added = append(added, newVal)
		}
	}
	return added, removed
}
