package roleset

import (
	"stash.corp.synacor.com/hack/werewolf/role"
)

func Fiver() *Roleset {
	return &Roleset{
		Name:        "Fast Fiver",
		Description: "Four villagers. One wolf. Two days to find them.",
		Roles: []*role.Role{
			role.Werewolf(),
			role.Cultist(),
			role.Hunter(),
			role.Seer(),
			role.Villager(),
		},
	}
}
