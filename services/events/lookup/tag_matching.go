package lookup

func isMatchingTags(eventTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	tags := make(map[string]struct{}, len(eventTags))
	for _, tag := range eventTags {
		tags[tag] = struct{}{}
	}

	for _, filterTag := range filterTags {
		if _, ok := tags[filterTag]; !ok {
			return false
		}
	}

	return true
}
