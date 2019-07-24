package main

func processQuestionCode(inputData responseData, questionCode string) bool{
	if inputData.QuestionCode == questionCode{
		return true
	}
	return false
}

func bothQAndDqFound(response responseData, config validationConfig) bool{
	question := processQuestionCode(response, config.questionCode)
	derivedQuestion := processQuestionCode(response, config.derivedQuestionCode)
	if question && derivedQuestion{
		return true
	}
	return false
}
