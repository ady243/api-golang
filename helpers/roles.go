package helpers

import "github.com/ady243/teamup/internal/models"

type Permission struct {
	CanEdit   bool
	CanDelete bool
	CanView   bool
}

var RolePermissions = map[models.Role]Permission{
	models.Player: {
		CanEdit:   true,
		CanDelete: true,
		CanView:   true,
	},
	models.Referee: {
		CanEdit:   true,
		CanDelete: false,
		CanView:   true,
	},
	models.Organizer: {
		CanEdit:   true,
		CanDelete: true,
		CanView:   true,
	},
}

func GetPermissions(role models.Role) Permission {
	return RolePermissions[role]
}
