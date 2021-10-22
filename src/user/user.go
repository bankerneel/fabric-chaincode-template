package user

import (
	"encoding/json"
	"fmt"
	"main/core/messages"
	"main/core/status"
	"main/core/utils"
	"main/src/schema"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//AddUser method to add user on blockchain network
func CreateAdmin(ctx contractapi.TransactionContextInterface, data []byte) (*utils.ResponseMessageWithId, error) {
	request := new(schema.CreateAdminRequest)
	// parse json bytes to struct
	err := json.Unmarshal(data, request)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	//send response userid id to client
	response := new(utils.ResponseMessageWithId)

	queryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"user_id\":%d}}", utils.DocTypeUser, utils.AdminId)
	admin, adminId, _ := utils.GetByQuery(ctx, queryString, fmt.Sprintf("admin not found"))

	if admin == nil {
		//create admin role
		adminRole := new(schema.Role)
		adminRole.DocType = utils.DocTypeRole
		adminRole.Name = utils.RoleAdmin
		adminRole.CreatedAt = request.CreatedAt
		//add into the database
		jsonData, err := json.Marshal(adminRole)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// if user not exist , add user to the network
		err = ctx.GetStub().PutState(request.AdminRoleUUID, jsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		//create viewer role
		viewerRole := new(schema.Role)
		viewerRole.DocType = utils.DocTypeRole
		viewerRole.Name = utils.RoleViewer
		viewerRole.CreatedAt = request.CreatedAt
		//add into the database
		jsonData, err = json.Marshal(viewerRole)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// if user not exist , add user to the network
		err = ctx.GetStub().PutState(request.ViewerRoleUUID, jsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		//create remover role
		removerRole := new(schema.Role)
		removerRole.DocType = utils.DocTypeRole
		removerRole.Name = utils.RoleRemover
		removerRole.CreatedAt = request.CreatedAt
		//add into the database
		jsonData, err = json.Marshal(removerRole)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// if user not exist , add user to the network
		err = ctx.GetStub().PutState(request.RemoverRoleUUID, jsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		//create admin
		admin := new(schema.User)
		admin.DocType = utils.DocTypeUser
		admin.UserID = utils.AdminId
		admin.FirstName = utils.RoleAdmin
		admin.LastName = utils.RoleAdmin
		var adminRoles []string
		adminRoles = append(adminRoles, request.AdminRoleUUID, request.RemoverRoleUUID, request.ViewerRoleUUID)
		admin.Roles = adminRoles
		admin.CreatedAt = request.CreatedAt
		//add into the database
		jsonData, err = json.Marshal(admin)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// if user not exist , add user to the network
		err = ctx.GetStub().PutState(request.AdminUUID, jsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		response.ID = request.AdminUUID
		response.Message = messages.AdminCreatedSuccess

	} else {
		response.ID = adminId
		response.Message = messages.AdminExist
	}

	return response, nil
}

//AddUser method to add user on blockchain network
func AddUser(ctx contractapi.TransactionContextInterface, data []byte) (*utils.ResponseMessageWithId, error) {
	request := new(schema.UserRequest)
	// parse json bytes to struct
	err := json.Unmarshal(data, request)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	// Validate the input data
	err = request.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}
	//send response userid id to client
	response := new(utils.ResponseMessageWithId)
	if request.User.UserID == 1 {
		response.Message = messages.UpdateAdminError
		return response, nil
	}

	queryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"user_id\":%d}}", utils.DocTypeUser, request.User.UserID)
	user, userId, _ := utils.GetByQuery(ctx, queryString, fmt.Sprintf("user not found"))

	if user == nil {
		// add doctype in the document
		request.User.DocType = utils.DocTypeUser
		// change structure to json bytes
		jsonData, err := json.Marshal(request.User)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// if user not exist , add user to the network
		err = ctx.GetStub().PutState(request.UUID, jsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		response.ID = request.UUID
		response.Message = messages.UserCreatedSuccess

	} else {
		userFound := new(schema.User)
		err := json.Unmarshal(user, userFound)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		userFound.FirstName = request.User.FirstName
		userFound.LastName = request.User.LastName
		userFound.UpdatedAt = request.User.CreatedAt

		userFoundJsonData, err := json.Marshal(userFound)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		err = ctx.GetStub().PutState(userId, userFoundJsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		response.ID = userId
		response.Message = messages.UserUpdateSuccess
	}

	return response, nil
}

//FetchUsers to get all user
//return array of user data or error if any occured
func FetchUsers(ctx contractapi.TransactionContextInterface) (*schema.UsersResponse, error) {
	//query to fetch all the timeline documents corresponding to contract
	queryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\"},\"sort\": [{\"created_at\": \"desc\"}]}", utils.DocTypeUser)

	//fetch data using the query
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	//Array of user user data
	var userResponse = new(schema.UsersResponse)

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		// iterate over the user data
		var user = new(schema.User)
		err = json.Unmarshal(queryResponse.Value, user)
		if err != nil {
			return nil, err
		}
		fmt.Println(user)
		if user.Roles == nil {
			var emptyArray []string
			user.Roles = emptyArray
		}
		//push timeline data to the array
		userResponse.Data = append(userResponse.Data, user)
	}
	if len(userResponse.Data) == 0 {
		return nil, status.ErrNotFound.WithMessage(messages.UserNotFound)
	}

	return userResponse, nil
}

//AddRoleToUser method to assign role to the user
func AddRoleToUser(ctx contractapi.TransactionContextInterface, data []byte) (*utils.ResponseMessageWithId, error) {
	request := new(schema.AddRoleToUserRequest)
	// parse json bytes to struct
	err := json.Unmarshal(data, request)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	// Validate the input data
	err = request.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	//send response userid id to client
	response := new(utils.ResponseMessageWithId)

	var RoleName string

	if request.Role == 1 {
		RoleName = utils.RoleViewer
	} else if request.Role == 2 {
		RoleName = utils.RoleRemover
	} else {
		response.Message = messages.WrongRole
		return response, nil
	}

	adminIsAdmin, err := UserHasRole(ctx, request.UserID, utils.RoleAdmin)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	if adminIsAdmin {

		userQueryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"user_id\":%d}}", utils.DocTypeUser, request.UpdateUserID)
		user, userId, err := utils.GetByQuery(ctx, userQueryString, fmt.Sprintf("updated user not found"))
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		roleQueryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"name\":\"%s\"}}", utils.DocTypeRole, RoleName)
		role, roleId, err := utils.GetByQuery(ctx, roleQueryString, fmt.Sprintf("role not found"))
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		if user != nil {
			response.ID = userId
			if role != nil {

				roleAdded := true

				userFound := new(schema.User)
				err := json.Unmarshal(user, userFound)
				if err != nil {
					return nil, status.ErrInternal.WithError(err)
				}
				roles := userFound.Roles
				for i := range roles {
					if roles[i] == roleId {
						roleAdded = false
					}
				}
				if roleAdded {
					roles = append(roles, roleId)
					userFound.Roles = roles
					userFound.UpdatedAt = request.CreatedAt
					jsonData, err := json.Marshal(userFound)
					if err != nil {
						return nil, status.ErrInternal.WithError(err)
					}
					err = ctx.GetStub().PutState(userId, jsonData)
					if err != nil {
						return nil, status.ErrInternal.WithError(err)
					}
					response.Message = messages.RoleAssignedSuccess
				} else {
					response.Message = messages.RoleAlreadyAssigned
				}
			} else {
				response.Message = messages.RoleNotExist
			}
		} else {
			response.Message = messages.UserNotExist
		}
	} else {
		response.Message = messages.UserNotAdmin
	}

	return response, nil
}

func UserHasRole(ctx contractapi.TransactionContextInterface, userId int32, roleName string) (bool, error) {

	userHasRole := false

	userQueryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"user_id\":%d}}", utils.DocTypeUser, userId)
	user, _, err := utils.GetByQuery(ctx, userQueryString, fmt.Sprintf("user not found"))
	if err != nil {
		return false, status.ErrInternal.WithError(err)
	}

	viewerRoleQueryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"name\":\"%s\"}}", utils.DocTypeRole, roleName)
	_, viewerRoleId, err := utils.GetByQuery(ctx, viewerRoleQueryString, fmt.Sprintf("viewer role not found"))
	if err != nil {
		return false, status.ErrInternal.WithError(err)
	}

	UserFound := new(schema.User)
	err = json.Unmarshal(user, UserFound)
	if err != nil {
		return false, status.ErrInternal.WithError(err)
	}

	userRoles := UserFound.Roles
	for i := range userRoles {
		if userRoles[i] == viewerRoleId {
			userHasRole = true
		}
	}

	return userHasRole, nil
}
