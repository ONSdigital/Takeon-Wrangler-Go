package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
)

func handle(inputData ContributorResponses) {
	var config []validationConfig
	config = append(config, validationConfig{questionCode: "601", derivedQuestionCode: "700"})
	config = append(config, validationConfig{questionCode: "602", derivedQuestionCode: "701"})
	loopOverConfig(config, inputData)
}

func main() {
	lambda.Start(handle)
}

func loopOverConfig(configData []validationConfig, responses ContributorResponses) {
	var wg sync.WaitGroup
	fmt.Printf("There are %d elements in config\n", len(configData))
	wg.Add(len(configData))
	for i := 0; i < len(configData); i++ {
		// var found isFound
		// found.dqCodeIsFound = false
		// found.qCodeIsFound = false
		var vetQandResponse VetQandResponse
		go loopOverResponse(responses.Responses, configData[i], vetQandResponse, &wg)
	}
	wg.Wait()
}

func loopOverResponse(responses []responseData, configData validationConfig, vetQandResponse VetQandResponse, waitGroup *sync.WaitGroup) {
	for i := 0; i < len(responses); i++ {
		if processQuestionCode(responses[i], configData.questionCode) {
			vetQandResponse.IsQuestionCodeFound = true
			vetQandResponse.FinalQuestCode = responses[i].QuestionCode
			vetQandResponse.FinalQuestCodeValue = responses[i].Response
		}

		if processQuestionCode(responses[i], configData.derivedQuestionCode) {
			vetQandResponse.IsDerivedQuestFound = true
			vetQandResponse.FinalDerivedQuestCode = responses[i].QuestionCode
			vetQandResponse.FinalDerivedQuestValue = responses[i].Response
		}
		fmt.Printf("responses: %s\n", responses[i])

		// fmt.Printf("found struct: %s\n", found)
		if vetQandResponse.IsQuestionCodeFound && vetQandResponse.IsDerivedQuestFound {
			var data PostData
			var valReq validationRequest
			valReq.PrimaryValue = vetQandResponse.FinalQuestCodeValue
			valReq.ComparisonValue = vetQandResponse.FinalDerivedQuestValue
			fmt.Printf("sending: %+v\n", valReq)
			output, err := json.Marshal(valReq)
			if err != nil {
				println("An error occured while marshaling: %s", err)
			}
			// found.dqCodeIsFound = false
			// found.qCodeIsFound = false
			data.Data = output
			data.LambdaLocation = os.Getenv("lambda")
			makePost(data)
			break

		}
	}
	fmt.Println("A statement just after for loop.")
	waitGroup.Done()
}
