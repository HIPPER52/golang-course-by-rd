package roles

type Role string

const (
	Admin    Role = "admin"
	Operator Role = "operator"
	Client   Role = "client"
)
