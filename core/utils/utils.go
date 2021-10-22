//Package utils for common method structure to use in all other packages
package utils

import (
	"main/core/status"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MetaData to use Common fields in all other structures
type MetaData struct {
	DocType   string    `json:"doc_type"` // type of the document i.e contract , clause etc
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ResponseID is used to return the response which contains only one ID field
type ResponseID struct {
	ID string `json:"id"`
}

// ResponseMessage is used to return the response which contains only one message field
type ResponseMessage struct {
	Message string `json:"message"`
}

type IDArrayResponse struct {
	ID []string `json:"id"`
}

type ResponseMessageWithIds struct {
	ID []int32 `json:"id"`
	ResponseMessage
}

type ResponseMessageWithId struct {
	ID      string `json:"id,omitempty" metadata:",optional"`
	Message string `json:"message"`
}

// GetByQuery executes the query and returns the byte array result
func GetByQuery(ctx contractapi.TransactionContextInterface, query string, message string) ([]byte, string, error) {
	//fetch the data from the query
	// fmt.Print("----- query --------", query)
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, "", status.ErrInternal.WithError(err)
	}

	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		return nil, "", status.ErrNotFound.WithMessage(message)
	}

	queryResponse, err := resultsIterator.Next()
	if err != nil {
		return nil, "", status.ErrInternal.WithError(err)
	}
	return queryResponse.Value, queryResponse.Key, nil
}
