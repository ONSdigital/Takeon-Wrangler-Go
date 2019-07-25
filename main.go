package main

import(
	"encoding/json"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
	"sync"
	"fmt"
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
	fmt.Printf("There are %d elements in config\n", len(configData))
	wg.Add(len(configData))
	for i:=0; i< len(configData); i++ {
		var found isFound
		found.dqCodeIsFound = false
		found.qCodeIsFound = false
		go loopOverResponse(responses.Responses, configData[i], found, &wg)
	}
	wg.Wait()
}

func loopOverResponse(responses []responseData, configData validationConfig, found isFound, waitGroup *sync.WaitGroup){
	for i:=0; i< len(responses); i++{
		if processQuestionCode(responses[i], configData.questionCode){
			found.qCodeIsFound = true
		}
		
		if processQuestionCode(responses[i], configData.derivedQuestionCode){
			found.dqCodeIsFound = true
		}
		fmt.Printf("responses: %s\n", responses[i])
		fmt.Printf("found struct: %s\n", found)
		if bothQAndDqFound(found){
			var data PostData
			fmt.Printf("sending: %+v\n", responses[i])
			output, err := json.Marshal(responses[i])
			if err != nil{
				println("An error occured while marshaling: %s", err)
			}
			found.dqCodeIsFound = false
			found.qCodeIsFound = false
			data.Data = output
			data.LambdaLocation = os.Getenv("lambda")
			makePost(data)
		}
	}
	waitGroup.Done()
}
