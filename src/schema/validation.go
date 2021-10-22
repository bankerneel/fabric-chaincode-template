package schema

import (
	"main/core/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (data UserRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UUID),
		validation.Field(&data.User),
	)
}

func (data User) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UserID, validation.Required.Error(utils.UserIdRequired), validation.NotNil.Error(utils.UserIdRequired)),
		validation.Field(&data.FirstName, validation.Required.Error(utils.FirstNameRequired), validation.NotNil.Error(utils.FirstNameRequired)),
		validation.Field(&data.LastName, validation.Required.Error(utils.LastNameRequired), validation.NotNil.Error(utils.LastNameRequired)),
	)
}

func (data RoleRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UUID),
		validation.Field(&data.Role),
	)
}

func (data Role) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Name, validation.Required.Error(utils.UserIdRequired), validation.NotNil.Error(utils.UserIdRequired)),
	)
}

func (data AddRoleToUserRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UserID, validation.Required.Error(utils.UserIdRequired), validation.NotNil.Error(utils.UserIdRequired)),
		validation.Field(&data.UpdateUserID, validation.Required.Error(utils.UserIdRequired), validation.NotNil.Error(utils.UserIdRequired)),
		validation.Field(&data.Role, validation.Required.Error(utils.RoleNameRequired), validation.NotNil.Error(utils.RoleNameRequired)),
	)
}

func (data ProductRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UUID),
		validation.Field(&data.Product),
	)
}

func (data Product) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.ProductId, validation.Required.Error(utils.ProductIdRequired), validation.NotNil.Error(utils.ProductIdRequired)),
		validation.Field(&data.Name, validation.Required.Error(utils.NameRequired), validation.NotNil.Error(utils.NameRequired)),
		validation.Field(&data.Category, validation.Required.Error(utils.CategoryRequired), validation.NotNil.Error(utils.CategoryRequired)),
		validation.Field(&data.Description, validation.Required.Error(utils.DescriptionRequired), validation.NotNil.Error(utils.DescriptionRequired)),
	)
}

func (data GetProductRequest) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UserID, validation.Required.Error(utils.UserIdRequired), validation.NotNil.Error(utils.UserIdRequired)),
		// validation.Field(&data.ProductID, validation.Required.Error(utils.ProductIdRequired), validation.NotNil.Error(utils.ProductIdRequired)),
	)
}
