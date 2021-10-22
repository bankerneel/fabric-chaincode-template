// Package utils common constants that are being used in other packages
package utils

// Constant to share in complete project
//Doctype for document types like clause , contract
const (
	ISDeleted      int    = 0
	DocTypeUser    string = "user"
	DocTypeRole    string = "role"
	DocTypeProduct string = "product"
	AdminId        int32  = 1
	RoleAdmin      string = "admin"
	RoleViewer     string = "viewer"
	RoleRemover    string = "remover"
	StatusConflict int    = 400
	ConflictError  string = "Conflict Error"
)
