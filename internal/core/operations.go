package core

import (
	"errors"

	"github.com/slashdevops/idp-scim-sync/internal/model"
)

var (
	ErrIdentityProviderGroupsMembersNil = errors.New("identity provider groups members is nil")
	ErrSCIMGroupsMembersNil             = errors.New("scim groups members is nil")
	ErrIdentityProviderGroupsNil        = errors.New("identity provider groups is nil")
	ErrSCIMGroupsNil                    = errors.New("scim groups is nil")
	ErrIdentityProviderUsersNil         = errors.New("identity provider users is nil")
	ErrSCIMUsersNil                     = errors.New("scim users is nil")
)

// groupsOperations returns the differences between the groups in the
// this use the Groups Name as the key.
// SCIM Groups cannot be updated.
// return 4 objest of GroupsMembersResult
// create: groups that exist in "scim" but not in "idp"
// update: groups that exist in "idp" and in "scim" but attributes changed in idp
// equal: groups that exist in both "idp" and "scim" and their attributes are equal
// delete: groups that exist in "scim" but not in "idp"
//
// also extract the id from scim to fill the resutls
func membersOperations(idp, scim *model.GroupsMembersResult) (create *model.GroupsMembersResult, equal *model.GroupsMembersResult, delete *model.GroupsMembersResult, err error) {
	if idp == nil {
		create, equal, delete, err = nil, nil, nil, ErrIdentityProviderGroupsMembersNil
		return
	}
	if scim == nil {
		create, equal, delete, err = nil, nil, nil, ErrSCIMGroupsMembersNil
		return
	}

	idpMemberSet := make(map[string]map[string]model.Member)  // [idp.GroupMembers.Group.Name] -> idp.member.Email -> idp.member
	scimMemberSet := make(map[string]map[string]model.Member) // [scim.Group.Name] -> [scim.member.Email] -> scim.member
	scimGroupsSet := make(map[string]model.Group)             // [scim.Group.Name] -> [scim.Group]

	toCreate := make([]model.GroupMembers, 0)
	toEqual := make([]model.GroupMembers, 0)
	toDelete := make([]model.GroupMembers, 0)

	for _, grpMembers := range idp.Resources {
		idpMemberSet[grpMembers.Group.Name] = make(map[string]model.Member)
		for _, member := range grpMembers.Resources {
			idpMemberSet[grpMembers.Group.Name][member.Email] = member
		}
	}

	for _, grpMembers := range scim.Resources {
		scimGroupsSet[grpMembers.Group.Name] = grpMembers.Group
		scimMemberSet[grpMembers.Group.Name] = make(map[string]model.Member)
		for _, member := range grpMembers.Resources {
			scimMemberSet[grpMembers.Group.Name][member.Email] = member
		}
	}

	for _, grpMembers := range idp.Resources {
		toC := make(map[string][]model.Member)
		toE := make(map[string][]model.Member)

		toC[grpMembers.Group.Name] = make([]model.Member, 0)
		toE[grpMembers.Group.Name] = make([]model.Member, 0)

		// count when both side have members == 0
		noMembers := 0

		// groups equals both sides without members
		if _, ok := scimMemberSet[grpMembers.Group.Name]; ok {
			if len(scimMemberSet[grpMembers.Group.Name]) == 0 && len(idpMemberSet[grpMembers.Group.Name]) == 0 {
				// toE[grpMembers.Group.Name] = make([]model.Member, 0)
				noMembers += 1
			}
		}

		// this case is when the groups is not new in scim
		if grpMembers.Group.SCIMID == "" {
			if _, ok := scimGroupsSet[grpMembers.Group.Name]; ok {
				grpMembers.Group.SCIMID = scimGroupsSet[grpMembers.Group.Name].SCIMID
			}
		}

		for _, member := range grpMembers.Resources {
			if _, ok := scimMemberSet[grpMembers.Group.Name][member.Email]; !ok {
				toC[grpMembers.Group.Name] = append(toC[grpMembers.Group.Name], member)
			} else {
				// check if the groups has the same members before adding to equal
				// TODO: check if the groups has the same members before adding to equal, what happens if some members are different?
				for grpMemberEmail := range scimMemberSet[grpMembers.Group.Name] {
					if grpMemberEmail == member.Email {
						member.SCIMID = scimMemberSet[grpMembers.Group.Name][member.Email].SCIMID
						toE[grpMembers.Group.Name] = append(toE[grpMembers.Group.Name], member)
					}
				}
			}
		}

		if len(toC[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()

			e := model.GroupMembers{
				Items:     len(toC[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toC[grpMembers.Group.Name],
			}
			e.SetHashCode()

			toCreate = append(toCreate, e)
		}

		if noMembers > 0 || len(toE[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()
			ee := model.GroupMembers{
				Items:     len(toE[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toE[grpMembers.Group.Name],
			}
			ee.SetHashCode()

			toEqual = append(toEqual, ee)
		}
	}

	for _, grpMembers := range scim.Resources {
		toD := make(map[string][]model.Member)
		toD[grpMembers.Group.Name] = make([]model.Member, 0)

		for _, member := range grpMembers.Resources {
			if _, ok := idpMemberSet[grpMembers.Group.Name][member.Email]; !ok {
				toD[grpMembers.Group.Name] = append(toD[grpMembers.Group.Name], member)
			}
		}

		if len(toD[grpMembers.Group.Name]) > 0 {
			grpMembers.Group.SetHashCode()

			e := model.GroupMembers{
				Items:     len(toD[grpMembers.Group.Name]),
				Group:     grpMembers.Group,
				Resources: toD[grpMembers.Group.Name],
			}
			e.SetHashCode()

			toDelete = append(toDelete, e)
		}
	}

	create = &model.GroupsMembersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	if create.Items > 0 {
		create.SetHashCode()
	}

	equal = &model.GroupsMembersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	if equal.Items > 0 {
		equal.SetHashCode()
	}

	delete = &model.GroupsMembersResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}
	if delete.Items > 0 {
		delete.SetHashCode()
	}

	return
}

// groupsOperations returns the differences between the groups in the
// this use the Groups Name as the key.
// SCIM Groups cannot be updated.
// return 4 objest of GroupsResult
// create: groups that exist in "scim" but not in "idp"
// update: groups that exist in "idp" and in "scim" but attributes changed in idp
// equal: groups that exist in both "idp" and "scim" and their attributes are equal
// delete: groups that exist in "scim" but not in "idp"
//
// also extract the id from scim to fill the resutls
func groupsOperations(idp, scim *model.GroupsResult) (create *model.GroupsResult, update *model.GroupsResult, equal *model.GroupsResult, delete *model.GroupsResult, err error) {
	if idp == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrIdentityProviderGroupsNil
		return
	}
	if scim == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrSCIMGroupsNil
		return
	}

	idpGroups := make(map[string]struct{})     // [idp.Group.Name ] -> struct{}{}
	scimGroups := make(map[string]model.Group) // [scim.Group.Name] -> scim.Group

	// log.Tracef("idp: %s\n, scim: %s\n", utils.ToJSON(idp), utils.ToJSON(scim))

	toCreate := make([]model.Group, 0)
	toUpdate := make([]model.Group, 0)
	toEqual := make([]model.Group, 0)
	toDelete := make([]model.Group, 0)

	for _, gr := range idp.Resources {
		idpGroups[gr.Name] = struct{}{}
	}

	for _, gr := range scim.Resources {
		scimGroups[gr.Name] = gr
	}

	// loop over idp to see what to create and what to update
	for _, group := range idp.Resources {
		if _, ok := scimGroups[group.Name]; !ok {
			toCreate = append(toCreate, group)
		} else {

			group.SCIMID = scimGroups[group.Name].SCIMID

			if group.IPID != scimGroups[group.Name].IPID {
				toUpdate = append(toUpdate, group)
			} else {
				toEqual = append(toEqual, group)
			}
		}
	}

	// loop over scim to see what to delete
	for _, group := range scim.Resources {
		if _, ok := idpGroups[group.Name]; !ok {
			toDelete = append(toDelete, group)
		}
	}

	create = &model.GroupsResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	if create.Items > 0 {
		create.SetHashCode()
	}

	update = &model.GroupsResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}
	if update.Items > 0 {
		update.SetHashCode()
	}

	equal = &model.GroupsResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	if equal.Items > 0 {
		equal.SetHashCode()
	}

	delete = &model.GroupsResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}
	if delete.Items > 0 {
		delete.SetHashCode()
	}

	return
}

// usersOperations returns the differences between the users in the
// Users Email as the key.
// return 4 objest of UsersResult
// create: users that exist in "scim" but not in "idp"
// update: users that exist in "idp" and in "scim" but attributes changed in idp
// equal: users that exist in both "idp" and "scim" and their attributes are equal
// delete: users that exist in "scim" but not in "idp"
func usersOperations(idp, scim *model.UsersResult) (create *model.UsersResult, update *model.UsersResult, equal *model.UsersResult, delete *model.UsersResult, err error) {
	if idp == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrIdentityProviderUsersNil
		return
	}
	if scim == nil {
		create, update, equal, delete, err = nil, nil, nil, nil, ErrSCIMUsersNil
		return
	}

	idpUsers := make(map[string]struct{})
	scimUsers := make(map[string]model.User)

	toCreate := make([]model.User, 0)
	toUpdate := make([]model.User, 0)
	toEqual := make([]model.User, 0)
	toDelete := make([]model.User, 0)

	for _, usr := range idp.Resources {
		idpUsers[usr.Email] = struct{}{}
	}

	for _, usr := range scim.Resources {
		scimUsers[usr.Email] = usr
	}

	// new users and what equal to them
	for _, usr := range idp.Resources {
		if _, ok := scimUsers[usr.Email]; !ok {
			toCreate = append(toCreate, usr)
		} else {

			usr.SCIMID = scimUsers[usr.Email].SCIMID

			if usr.Name.FamilyName != scimUsers[usr.Email].Name.FamilyName ||
				usr.Name.GivenName != scimUsers[usr.Email].Name.GivenName ||
				usr.Active != scimUsers[usr.Email].Active ||
				usr.IPID != scimUsers[usr.Email].IPID {

				toUpdate = append(toUpdate, usr)
			} else {
				toEqual = append(toEqual, usr)
			}
		}
	}

	for _, usr := range scim.Resources {
		if _, ok := idpUsers[usr.Email]; !ok {
			toDelete = append(toDelete, usr)
		}
	}

	create = &model.UsersResult{
		Items:     len(toCreate),
		Resources: toCreate,
	}
	if create.Items > 0 {
		create.SetHashCode()
	}

	update = &model.UsersResult{
		Items:     len(toUpdate),
		Resources: toUpdate,
	}
	if update.Items > 0 {
		update.SetHashCode()
	}

	equal = &model.UsersResult{
		Items:     len(toEqual),
		Resources: toEqual,
	}
	if equal.Items > 0 {
		equal.SetHashCode()
	}

	delete = &model.UsersResult{
		Items:     len(toDelete),
		Resources: toDelete,
	}
	if delete.Items > 0 {
		delete.SetHashCode()
	}

	return
}

func mergeGroupsResult(grs ...*model.GroupsResult) (merged model.GroupsResult) {
	groups := make([]model.Group, 0)

	for _, gr := range grs {
		groups = append(groups, gr.Resources...)
	}

	merged = model.GroupsResult{
		Items:     len(groups),
		Resources: groups,
	}
	if merged.Items > 0 {
		merged.SetHashCode()
	}

	return
}

func mergeUsersResult(urs ...*model.UsersResult) (merged model.UsersResult) {
	users := make([]model.User, 0)

	for _, u := range urs {
		users = append(users, u.Resources...)
	}

	merged = model.UsersResult{
		Items:     len(users),
		Resources: users,
	}
	if merged.Items > 0 {
		merged.SetHashCode()
	}

	return
}

func mergeGroupsMembersResult(gms ...*model.GroupsMembersResult) (merged model.GroupsMembersResult) {
	groupsMembers := make([]model.GroupMembers, 0)

	for _, gm := range gms {
		groupsMembers = append(groupsMembers, gm.Resources...)
	}

	merged = model.GroupsMembersResult{
		Items:     len(groupsMembers),
		Resources: groupsMembers,
	}
	if merged.Items > 0 {
		merged.SetHashCode()
	}

	return
}
