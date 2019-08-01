package main

import (
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

func TestBothFound(t *testing.T) {
	var respData responseData
	var dqData responseData
	var config validationConfig
	respData.QuestionCode = "123"
	dqData.QuestionCode = "123"
	config.questionCode = "123"
	if bothQAndDqFound(respData, config) != true {
		t.Errorf("Function did not return true")
	}
	respData.QuestionCode = "111"
	if bothQAndDqFound(respData, config) != false {
		t.Errorf("Function didn't return false")
	}
}
