package main

import (
	"log"
	"main/core/messages"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PRODUCTChainCode implementation
type PRODUCTChainCode struct {
	contractapi.Contract
}

func main() {
	PRODUCTChainCode, err := contractapi.NewChaincode(&PRODUCTChainCode{})
	if err != nil {
		log.Panicf(messages.ChaincodeCreateError, err.Error())
	}

	if err := PRODUCTChainCode.Start(); err != nil {
		log.Panicf(messages.ChaincodeStartError, err.Error())
	}
}
