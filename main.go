package main

import(
	"encoding/json"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
	"sync"
)

func handle(inputData ContributorResponses){
	var config []validationConfig
	config = append(config, validationConfig{questionCode: "601", derivedQuestionCode: "700"})
	config = append(config, validationConfig{questionCode: "602", derivedQuestionCode: "701"})
	loopOverConfig(config, inputData)
}

func main(){
	lambda.Start(handle)
}

func loopOverConfig(configData []validationConfig, responses ContributorResponses){
	var wg sync.WaitGroup
	wg.Add(len(configData))
	for i:=0; i< len(configData); i++ {
		go loopOverResponse(responses.Responses, configData[i], &wg)
	}
	wg.Wait()
}

func loopOverResponse(responses []responseData, configData validationConfig, waitGroup *sync.WaitGroup){
	for i:=0; i< len(responses); i++{
		bothQAndDqFound(responses[i], configData)
			var data PostData
			output, err := json.Marshal(responses[i])
			if err != nil{
				println("An error occured while marshaling: %s", err)
			}
			data.Data = output
			data.LambdaLocation = os.Getenv("lambda")
			makePost(data)
	}
	waitGroup.Done()
}
