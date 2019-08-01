package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func handle(inputData ContributorResponses) {
	var config []validationConfig
	config = append(config, validationConfig{questionCode: "601", derivedQuestionCode: "700"})
	config = append(config, validationConfig{questionCode: "602", derivedQuestionCode: "701"})
	loopOverConfig(config, inputData)
	// topic := os.Getenv("TOPIC_ARN")
	// sendBpmResponse("B123-P456-M789", "QvsDQ", topic)
}

func main() {
	lambda.Start(handle)
}

func loopOverConfig(configData []validationConfig, responses ContributorResponses) {
	var wg sync.WaitGroup
	vetqResponses := []VetQResponse{}
	fmt.Printf("Wrangler Input: %+v\n", responses)
	fmt.Printf("Config: %+v\n", configData)
	wg.Add(len(configData))
	for i := 0; i < len(configData); i++ {
		// var vetQDqResponse VetQDqResponse
		// go loopOverResponse(responses.Responses, configData[i], VetQDqResponse, &wg)
		var vetqResponse VetQResponse
		vetqResponse = loopOverResponse(responses.Responses, configData[i], &wg)
		vetqResponses = append(vetqResponses, vetqResponse)
	}
	wg.Wait()
	var persistInpData PersistInpData
	var persistPostReq PersistPostReq
	persistInpData.Period = responses.Period
	persistInpData.Reference = responses.Reference
	persistInpData.Survey = responses.Survey
	persistInpData.Instance = responses.BpmID
	persistInpData.ValidationName = "QvDQ"
	persistInpData.ValidationResults = vetqResponses

	persistPostReq.Type = "POST"
	persistPostReq.Input = persistInpData

	DataToPersist, err := json.Marshal(persistPostReq)
	if err != nil {
		println("An error occured while marshaling persistPostReq: %s", err)
	}

	PersistLambdaName := os.Getenv("PERSIST_LAMBDA")

	PersistResult := invokeLambda(PersistLambdaName, DataToPersist)
	var PersistRetn string
	json.Unmarshal(PersistResult, &PersistRetn)

}

func loopOverResponse(responses []responseData, configData validationConfig, waitGroup *sync.WaitGroup) VetQResponse {
	var valdResult validationResult
	var vetqResponse VetQResponse
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
			// var data PostData
			var valReq validationRequest
			valReq.PrimaryValue = vetQDqResponse.FinalQuestCodeValue
			valReq.ComparisonValue = vetQDqResponse.FinalDerivedQuestValue
			ValueCompData, err := json.Marshal(valReq)
			if err != nil {
				println("An error occured while marshaling: %s", err)
			}
			ValueCompLambda := os.Getenv("VALUE_COMPARISION_LAMBDA")
			ValueCompResult := invokeLambda(ValueCompLambda, ValueCompData)
			json.Unmarshal(ValueCompResult, &valdResult)
			vetqResponse = VetQResponse{QuestionCode: vetQDqResponse.FinalQuestCode,
				ValueFormula: valdResult.ValueFormula,
				Triggered:    valdResult.Triggered}
			break

		}

	}
	waitGroup.Done()
	return vetqResponse
}

func sendBpmResponse(bpmInstance string, validationName string, topicArn string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})

	if err != nil {
		fmt.Println("NewSession error:", err)
		return
	}

	client := sns.New(sess)
	input := &sns.PublishInput{
		Message:  aws.String(bpmInstance + validationName),
		TopicArn: aws.String(topicArn),
	}

	result, err := client.Publish(input)

	if err != nil {
		fmt.Println("Publish error:", err)
		return
	}

	fmt.Println("sendBpmResponse:")
	fmt.Println(result)

}
