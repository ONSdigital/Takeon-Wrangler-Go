package main

import (
	"reflect"
	"testing"
)

func TestUnpack(t *testing.T) {
	test := `{"reference": "123456789", "period": "202013", "survey": "001", "responses": [{"questionCode": "123", "response": "456"}]}`
	x := unpackJSON(test)
	testStrut := ContributorResponses{}

	testStrut.Reference = "123456789"
	testStrut.Period = "202013"
	testStrut.Survey = "001"

	if testStrut.Reference != x.Reference {
		t.Errorf("References Don't match! Expected: %s but got: %s", x.Reference, testStrut.Reference)
	}
}

func TestCheck(t *testing.T) {
	var testData responseData
	testData.QuestionCode = "123"
	testData.Response = "001"
	qCode := "123"
	x := processQuestionCode(testData, qCode)
	if x != true {
		t.Errorf("Function did not return true. Question Data qCode: %s\nTest qCode: %s", testData.QuestionCode, qCode)
	}
	qCode = "541"
	x = processQuestionCode(testData, qCode)
	if x != false {
		t.Errorf("Function did not return true. Question Data qCode: %s\nTest qCode: %s", testData.QuestionCode, qCode)
	}
}

func TestLoopOverResponse(t *testing.T) {
	cvalqRequest := make(chan validationQRequest)
	var config = validationConfig{questionCode: "601", derivedQuestionCode: "700"}
	// config = append(config, validationConfig{questionCode: "602", derivedQuestionCode: "701"})
	var responses = []responseData{responseData{QuestionCode: "601", Response: "146"},
		responseData{QuestionCode: "602", Response: "150"},
		responseData{QuestionCode: "700", Response: "148"},
		responseData{QuestionCode: "603", Response: "602"},
		responseData{QuestionCode: "701", Response: "603"},
		responseData{QuestionCode: "605", Response: "604"},
		responseData{QuestionCode: "606", Response: "605"},
		responseData{QuestionCode: "607", Response: "606"},
		responseData{QuestionCode: "608", Response: "607"},
		responseData{QuestionCode: "609", Response: "608"},
	}
	go loopOverResponse(responses, config, cvalqRequest)
	got := <-cvalqRequest
	var want = validationQRequest{Qcode: "601", PrimaryValue: "146", ComparisonValue: "148", MetaData: ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("loopOverResponse(responses, config, cvalqRequest) = %v, want %v", got, want)
	}

}

func TestLoopOverConfig(t *testing.T) {
	var config []validationConfig
	config = append(config, validationConfig{questionCode: "601", derivedQuestionCode: "700"})
	config = append(config, validationConfig{questionCode: "602", derivedQuestionCode: "701"})

	var inputData = ContributorResponses{Reference: "4990012",
		Period: "201211",
		Survey: "066",
		BpmID:  "instanceid",
		Responses: []responseData{responseData{QuestionCode: "601", Response: "146"},
			responseData{QuestionCode: "602", Response: "150"},
			responseData{QuestionCode: "700", Response: "148"},
			responseData{QuestionCode: "603", Response: "602"},
			responseData{QuestionCode: "701", Response: "603"},
			responseData{QuestionCode: "605", Response: "604"},
			responseData{QuestionCode: "606", Response: "605"},
			responseData{QuestionCode: "607", Response: "606"},
			responseData{QuestionCode: "608", Response: "607"},
			responseData{QuestionCode: "609", Response: "608"},
		},
	}

	cvalReqSlice := make(chan []validationQRequest)
	go loopOverConfig(config, inputData, cvalReqSlice)
	got := <-cvalReqSlice
	var want = []validationQRequest{validationQRequest{Qcode: "601", PrimaryValue: "146", ComparisonValue: "148", MetaData: ""},
		validationQRequest{Qcode: "602", PrimaryValue: "150", ComparisonValue: "603", MetaData: ""},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("loopOverConfig(config, inputData, cvalReqSlice) = %v, want %v", got, want)
	}

}

// func TestBothFound(t *testing.T) {
// 	var respData responseData
// 	var dqData responseData

// 	var config validationConfig
// 	respData.QuestionCode = "123"
// 	dqData.QuestionCode = "123"
// 	config.questionCode = "123"
// 	if bothQAndDqFound(respData, config) != true {
// 		t.Errorf("Function did not return true")
// 	}
// 	respData.QuestionCode = "111"
// 	if bothQAndDqFound(respData, config) != false {
// 		t.Errorf("Function didn't return false")
// 	}
// }
