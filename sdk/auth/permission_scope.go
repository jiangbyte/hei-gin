package auth

const (
	permissionPathDirect   = "P0"
	permissionPathUserRole = "P1"
)

var permissionPathPriority = map[string]int{
	permissionPathDirect:   0,
	permissionPathUserRole: 1,
}

var dataScopePriority = map[string]int{
	"SELF":             0,
	"CUSTOM_GROUP":     1,
	"CUSTOM_ORG":       2,
	"GROUP_AND_BELOW":  3,
	"GROUP":            4,
	"ORG_AND_BELOW":    5,
	"ORG":              6,
	"ALL":              7,
}

func mostRestrictiveScope(scopes ...string) string {
	if len(scopes) == 0 {
		return ""
	}

	result := scopes[0]
	minPriority, ok := dataScopePriority[result]
	if !ok {
		minPriority = len(dataScopePriority) + 1
	}

	for _, scope := range scopes[1:] {
		priority, exists := dataScopePriority[scope]
		if !exists {
			continue
		}
		if priority < minPriority {
			minPriority = priority
			result = scope
		}
	}

	return result
}
