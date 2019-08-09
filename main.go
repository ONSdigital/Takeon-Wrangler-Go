package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var region = os.Getenv("AWS_REGION")

func main() {
	lambda.Start(handle)
}

func handle(inputData ContributorResponses) {
	var config []validationConfig
	//configs of questionCode and  derivedQuestionCode are assigned here
	config = append(config, validationConfig{questionCode: "601", derivedQuestionCode: "700"})
	config = append(config, validationConfig{questionCode: "602", derivedQuestionCode: "701"})

	//Define the channel to get the responses that matches the config's questionCode and  derivedQuestionCode
	cvalReqSlice := make(chan []validationQRequest)
	go loopOverConfig(config, inputData, cvalReqSlice)
	valReqSlice := <-cvalReqSlice
	fmt.Printf("valReqSlice %v\n", valReqSlice)

	//send the responses to value comparision lambda to get the validation results of the responses
	cvalResultSlice := make(chan []VetQResponse)
	go callValueCompLambda(valReqSlice, cvalResultSlice)

	inputMetadata := InputMetaData{inputData.Reference, inputData.Period, inputData.Survey, inputData.BpmID}
	var wg sync.WaitGroup
	wg.Add(1)
	//save the resultant responses including validation results into S3
	go callPersistLambda(cvalResultSlice, inputMetadata, &wg)
	wg.Wait()
}

//function to persist data into s3 bucket
func callPersistLambda(cvalResultSlice chan []VetQResponse, inputMetadata InputMetaData, waitGroup *sync.WaitGroup) {
	var persistInpData PersistInpData
	persistInpData.Period = inputMetadata.Period
	persistInpData.Reference = inputMetadata.Reference
	persistInpData.Survey = inputMetadata.Survey
	persistInpData.Instance = inputMetadata.BpmID
	persistInpData.ValidationName = "QvDQ"
	ValidationResults := <-cvalResultSlice
	persistInpData.ValidationResults = ValidationResults

	DataToPersist, err := json.Marshal(persistInpData)
	if err != nil {
		fmt.Printf("An error occured while marshaling persistPostReq: %s", err)
	}

	fmt.Printf("Data going to save into s3 bucket.\n")
	//fmt.Printf("Data going to save into s3 bucket: %v\n", string(DataToPersist))

	reader := strings.NewReader(string(DataToPersist))

	fmt.Printf("Region: %q\n", region)
	config := &aws.Config{
		Region: aws.String(region),
	}

	sess := session.New(config)

	uploader := s3manager.NewUploader(sess)

	bucket := os.Getenv("S3_BUCKET")
	filename := strings.Join([]string{persistInpData.Survey, persistInpData.Period, persistInpData.Reference, persistInpData.Instance, persistInpData.ValidationName}, "")
	fmt.Printf("Bucket filename: %q\n", filename)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   reader,
	})

	if err != nil {
		fmt.Printf("Unable to upload %q to %q, %v", filename, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to s3 bucket %q\n", filename, bucket)
	waitGroup.Done()
}

//This function calls value comparision lambda
func callValueCompLambda(valReqSlice []validationQRequest, cvalResultSlice chan []VetQResponse) {
	var valdResult validationResult
	valAllResult := []VetQResponse{}
	var vetqResponse VetQResponse
	for _, valReq := range valReqSlice {
		valReqItem := validationRequest{valReq.PrimaryValue, valReq.ComparisonValue, valReq.MetaData}
		ValueCompData, err := json.Marshal(valReqItem)
		if err != nil {
			println("An error occured while marshaling valReq: %s", err)
		}
		ValueCompLambda := os.Getenv("VALUE_COMPARISION_LAMBDA")
		ValueCompResult := invokeLambda(ValueCompLambda, ValueCompData)
		json.Unmarshal(ValueCompResult, &valdResult)
		vetqResponse = VetQResponse{QuestionCode: valReq.Qcode,
			ValueFormula: valdResult.ValueFormula,
			Triggered:    valdResult.Triggered}
		valAllResult = append(valAllResult, vetqResponse)

	}

	cvalResultSlice <- valAllResult

}

//This function will loop through all the configs and will evaluate the input responses.
func loopOverConfig(configData []validationConfig, responses ContributorResponses, cvalReqSlice chan []validationQRequest) {
	valAllRequest := []validationQRequest{}
	fmt.Printf("Wrangler Input: %+v\n", responses)
	fmt.Printf("Config: %+v\n", configData)
	for i := 0; i < len(configData); i++ {
		cvalqRequest := make(chan validationQRequest)
		go loopOverResponse(responses.Responses, configData[i], cvalqRequest)
		valRequest := <-cvalqRequest
		if valRequest.Qcode != "" {
			valAllRequest = append(valAllRequest, valRequest)
		}
	}
	cvalReqSlice <- valAllRequest
}

//This function will resturn the responses that matches the questionCode and derivedQuestionCode of a particular config.
func loopOverResponse(responses []responseData, configData validationConfig, cvalReq chan validationQRequest) {
	var valReq validationQRequest
	var vetQDqResponse VetQDqResponse
	for i := 0; i < len(responses); i++ {
		if processQuestionCode(responses[i], configData.questionCode) {
			vetQDqResponse.IsQuestionCodeFound = true
			vetQDqResponse.FinalQuestCode = responses[i].QuestionCode
			vetQDqResponse.FinalQuestCodeValue = responses[i].Response
		}

		if processQuestionCode(responses[i], configData.derivedQuestionCode) {
			vetQDqResponse.IsDerivedQuestFound = true
			vetQDqResponse.FinalDerivedQuestCode = responses[i].QuestionCode
			vetQDqResponse.FinalDerivedQuestValue = responses[i].Response
		}
		// Check if both QuestionCode and DerivedQuestionCode are found
		if vetQDqResponse.IsQuestionCodeFound && vetQDqResponse.IsDerivedQuestFound {
			valReq.Qcode = vetQDqResponse.FinalQuestCode
			valReq.PrimaryValue = vetQDqResponse.FinalQuestCodeValue
			valReq.ComparisonValue = vetQDqResponse.FinalDerivedQuestValue
			cvalReq <- validationQRequest{vetQDqResponse.FinalQuestCode, vetQDqResponse.FinalQuestCodeValue, vetQDqResponse.FinalDerivedQuestValue, ""}
			break

		}
	}
	if vetQDqResponse.FinalQuestCodeValue == "" || vetQDqResponse.FinalDerivedQuestValue == "" {
		cvalReq <- validationQRequest{}
	}
}

// func sendBpmResponse(bpmInstance string, validationName string, topicArn string) {

// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("eu-west-2"),
// 	})

// 	if err != nil {
// 		fmt.Println("NewSession error:", err)
// 		return
// 	}

// 	client := sns.New(sess)
// 	input := &sns.PublishInput{
// 		Message:  aws.String(bpmInstance + validationName),
// 		TopicArn: aws.String(topicArn),
// 	}

// 	result, err := client.Publish(input)

// 	if err != nil {
// 		fmt.Println("Publish error:", err)
// 		return
// 	}

// 	fmt.Println("sendBpmResponse:")
// 	fmt.Println(result)

// }
