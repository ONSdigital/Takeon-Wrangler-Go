package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

type PostData struct{
	Data []byte
	LambdaLocation string
}

func makePost(data PostData){
	svc := lambda.New(session.New())

	fmt.Printf(string(data.Data))

	fmt.Printf(string(data.LambdaLocation)+"\n")

	input := &lambda.InvokeInput{
		FunctionName: aws.String(data.LambdaLocation),
		Payload: data.Data,
		InvocationType: aws.String("RequestResponse"),
	}
	fmt.Printf("sending from makePost: %+v\n", string(data.Data))
	result, err := svc.Invoke(input)
	if err != nil{
		fmt.Errorf("There was an error during post: %s\n", err.Error)
		panic(err)
	}
	fmt.Printf("Response: %s\n", result.Payload)
}
