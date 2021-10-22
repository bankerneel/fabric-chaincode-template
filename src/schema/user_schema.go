package schema

import "main/core/utils"

//User structure
type User struct {
	UserID    int32    `json:"user_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles,omitempty" metadata:",optional"`
	utils.MetaData
}

//User request structure
type UserRequest struct {
	UUID string `json:"uuid"`
	User
}

type CreateAdminRequest struct {
	AdminUUID       string `json:"admin_uuid"`
	AdminRoleUUID   string `json:"admin_role_uuid"`
	ViewerRoleUUID  string `json:"viewer_role_uuid"`
	RemoverRoleUUID string `json:"remover_role_uuid"`
	utils.MetaData
}

//PaginationTermResponse structure
type UsersResponse struct {
	Data []*User `json:"user"`
}

//Role structure
type Role struct {
	Name string `json:"name"`
	utils.MetaData
}

//Role request structure
type RoleRequest struct {
	UUID string `json:"uuid"`
	Role
}

//AddRoleToUser request structure
type AddRoleToUserRequest struct {
	UserID       int32 `json:"user_id"`
	UpdateUserID int32 `json:"update_user_id"`
	Role         int   `json:"role"`
	utils.MetaData
}
