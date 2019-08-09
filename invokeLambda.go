package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func invokeLambda(LambdaName string, Data []byte) []byte {

	svc := lambda.New(session.New())

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(LambdaName),
		Payload:        Data,
		InvocationType: aws.String("RequestResponse"),
	}
	fmt.Printf("sending to "+LambdaName+" from invokeLambda: %+v\n", string(Data))
	result, err := svc.Invoke(input)
	if err != nil {
		fmt.Println("error in invokeLambda: ")
		fmt.Println(err.Error())
		// fmt.Errorf("There was an error during invokeLambda: %s", err)
		panic(err)
	}
	fmt.Printf("Response from "+LambdaName+" : %s\n", result.Payload)
	return result.Payload
}
