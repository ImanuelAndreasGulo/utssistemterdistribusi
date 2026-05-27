package counter

// MergeState merges two G-Counter states by taking the maximum value per node.
// It returns the merged state and a flag indicating whether any change occurred.
func MergeState(current, incoming map[string]int) (map[string]int, bool) {
	merged := make(map[string]int, len(current)+len(incoming))
	changed := false

	for k, v := range current {
		merged[k] = v
	}

	for k, v := range incoming {
		if currentVal, ok := merged[k]; !ok || v > currentVal {
			merged[k] = v
			if !ok || v != currentVal {
				changed = true
			}
		}
	}

	if len(merged) != len(current) {
		changed = true
	}

	return merged, changed
}
