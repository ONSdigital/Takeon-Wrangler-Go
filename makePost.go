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

	input := &lambda.InvokeInput{
		FunctionName: aws.String(data.LambdaLocation),
		Payload: data.Data,	
	}
	
	result, err := svc.Invoke(input)

	if err != nil{
		fmt.Errorf("There was an error during post: %s\n", err.Error)
		panic(err)
	}
	println("Response: %s\n", result.StatusCode)
}