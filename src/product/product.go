package product

import (
	"encoding/json"
	"fmt"
	"main/core/messages"
	"main/core/status"
	"main/core/utils"
	"main/src/schema"
	"main/src/user"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//AddProduct method to add product on blockchain network
//return an id and message if success error if any occurred
func AddProduct(ctx contractapi.TransactionContextInterface, data []byte) (*utils.ResponseMessageWithId, error) {
	request := new(schema.ProductRequest)
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
	//send response  to client
	response := new(utils.ResponseMessageWithId)

	queryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"product_id\":%d}}", utils.DocTypeProduct, request.Product.ProductId)
	product, productId, _ := utils.GetByQuery(ctx, queryString, fmt.Sprintf("user not found"))

	if product == nil {
		// add doctype in the document
		request.Product.DocType = utils.DocTypeProduct
		// change structure to json bytes
		jsonData, err := json.Marshal(request.Product)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// if user not exist , add user to the network
		err = ctx.GetStub().PutState(request.UUID, jsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		response.ID = request.UUID
		response.Message = messages.ProductCreatedSuccess

	} else {
		productFound := new(schema.Product)
		err := json.Unmarshal(product, productFound)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		productFound.Name = request.Product.Name
		productFound.Category = request.Product.Category
		productFound.Description = request.Product.Description
		productFound.UpdatedAt = request.Product.CreatedAt

		productFoundJsonData, err := json.Marshal(productFound)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		err = ctx.GetStub().PutState(productId, productFoundJsonData)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		response.ID = productId
		response.Message = messages.ProductUpdateSuccess
	}

	return response, nil
}

//FetchProduct to get a particular product
//return a product data or error if any occurred
func ViewProduct(ctx contractapi.TransactionContextInterface, data []byte) (*schema.GetProductResponse, error) {
	request := new(schema.GetProductRequest)
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

	response := new(schema.GetProductResponse)

	//check if user have the viewer role
	userHasViewerRole, err := user.UserHasRole(ctx, request.UserID, utils.RoleViewer)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	if !userHasViewerRole {
		return nil, status.ErrNotFound.WithMessage(messages.UserNotAuthorized)
	}

	//query to fetch all the timeline documents corresponding to contract
	productQueryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"product_id\":%d}}", utils.DocTypeProduct, request.ProductId)

	if request.ProductId == 0 {
		productQueryString = fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\"},\"sort\": [{\"created_at\": \"desc\"}]}", utils.DocTypeProduct)
	}

	//fetch data using the query
	resultsIterator, err := ctx.GetStub().GetQueryResult(productQueryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		// iterate over the user data
		var product = new(schema.Product)
		err = json.Unmarshal(queryResponse.Value, product)
		if err != nil {
			return nil, err
		}
		//push timeline data to the array
		response.Product = append(response.Product, product)
	}
	if len(response.Product) == 0 {
		return nil, status.ErrNotFound.WithMessage(messages.ProductNotFound)
	}

	return response, nil
}

// DeleteProductDocType deletes an given asset from the world state.
func DeleteProductDocType(ctx contractapi.TransactionContextInterface, data []byte) (*utils.ResponseMessageWithIds, error) {
	request := new(schema.GetProductRequest)
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

	response := new(utils.ResponseMessageWithIds)

	//check if user have the viewer role
	userHasRemoverRole, err := user.UserHasRole(ctx, request.UserID, utils.RoleRemover)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	if !userHasRemoverRole {
		return nil, status.ErrNotFound.WithMessage(messages.UserNotAuthorized)
	}

	//query to fetch all the timeline documents corresponding to contract
	productQueryString := fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\",\"product_id\":%d}}", utils.DocTypeProduct, request.ProductId)

	if request.ProductId == 0 {
		productQueryString = fmt.Sprintf("{\"selector\":{\"doc_type\":\"%s\"},\"sort\": [{\"created_at\": \"desc\"}]}", utils.DocTypeProduct)
	}

	//fetch data using the query
	resultsIterator, err := ctx.GetStub().GetQueryResult(productQueryString)
	if err != nil {
		return nil, err
	}
	var idSlice []int32
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		// iterate over the user data
		var product = new(schema.Product)
		err = json.Unmarshal(queryResponse.Value, product)
		if err != nil {
			return nil, err
		}

		idSlice = append(idSlice, product.ProductId)

		// update the contract details on the network
		errStub := ctx.GetStub().DelState(queryResponse.Key)
		if errStub != nil {
			return nil, errStub
		}
	}
	response.ID = idSlice
	if len(response.ID) == 0 {
		return nil, status.ErrNotFound.WithMessage(messages.ProductNotFound)
	}
	response.Message = messages.ProductRemoved

	return response, nil
}
