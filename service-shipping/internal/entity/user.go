package entity

type Role int16

// Role of user.
const (
	RoleAdmin              Role = 1 //00001
	RoleShop               Role = 2 //00010
	RoleAgentNetworkMember Role = 4 //00100
	RoleAccounting         Role = 8 //01000
)

var Roles = []Role{RoleAdmin, RoleShop, RoleAgentNetworkMember, RoleAccounting}

// roles map between role uint16 and role string.
var roles map[Role]string

//User user entity
type User struct {
	Id       string
	Name     string
	Email    string
	Role     Role
	Password string
	CreateDelete
}

func init() {
	roles = make(map[Role]string)
	roles[RoleAdmin] = "Admin"
	roles[RoleShop] = "Shop"
	roles[RoleAgentNetworkMember] = "Agent Network Member"
	roles[RoleAccounting] = "Accounting"
}

// GetRoleName return a role string of user.
func (u User) GetRoleString() string {
	return roles[u.Role]
}

func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u User) IsShop() bool {
	return u.Role == RoleShop
}
