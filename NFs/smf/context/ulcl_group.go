package context

func GetULCLGroupNameFromSUPI(SUPI string) string {
	ulclGroups := smfContext.ULCLGroups
	for name, group := range ulclGroups {
		for _, member := range group {
			if member == SUPI {
				return name
			}
		}
	}
	return ""
}
