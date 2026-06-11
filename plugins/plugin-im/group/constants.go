package group

// Group types
const (
	GroupTypeMixed        = "mixed"
	GroupTypeConsumerOnly = "consumer_only"
)

// Member roles
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// Member statuses
const (
	MemberActive = "active"
	MemberLeft   = "left"
	MemberKicked = "kicked"
)

// Group statuses
const (
	GroupNormal    = "normal"
	GroupDissolved = "dissolved"
)
