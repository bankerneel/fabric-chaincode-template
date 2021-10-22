package main

import (
	"main/core/utils"
	"main/src/product"
	"main/src/schema"
	"main/src/user"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*---------- -----------------user routes ---------------*/

//add admin user and admin,viewer,remover roles
func (s *PRODUCTChainCode) AddAdmin(ctx contractapi.TransactionContextInterface, data string) (*utils.ResponseMessageWithId, error) {
	return user.CreateAdmin(ctx, []byte(data))
}

//AddUser route to write new user
func (s *PRODUCTChainCode) AddUser(ctx contractapi.TransactionContextInterface, data string) (*utils.ResponseMessageWithId, error) {
	return user.AddUser(ctx, []byte(data))
}

//FetchUsers route to read Users
func (s *PRODUCTChainCode) FetchUsers(ctx contractapi.TransactionContextInterface) (*schema.UsersResponse, error) {
	return user.FetchUsers(ctx)
}

//Add viwer or remover role to user
func (s *PRODUCTChainCode) AddRoleToUser(ctx contractapi.TransactionContextInterface, data string) (*utils.ResponseMessageWithId, error) {
	return user.AddRoleToUser(ctx, []byte(data))
}

//AddProduct to write new product
func (s *PRODUCTChainCode) AddProduct(ctx contractapi.TransactionContextInterface, data string) (*utils.ResponseMessageWithId, error) {
	return product.AddProduct(ctx, []byte(data))
}

//ViewProduct to view single or all product if user have viewer role
func (s *PRODUCTChainCode) ViewProduct(ctx contractapi.TransactionContextInterface, data string) (*schema.GetProductResponse, error) {
	return product.ViewProduct(ctx, []byte(data))
}

//DeleteProductDocType route to delete a single product or all product if user have remover role
func (s *PRODUCTChainCode) DeleteProductDocType(ctx contractapi.TransactionContextInterface, data string) (*utils.ResponseMessageWithIds, error) {
	return product.DeleteProductDocType(ctx, []byte(data))
}
