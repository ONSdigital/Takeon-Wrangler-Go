package main

func processQuestionCode(inputData responseData, questionCode string) bool{
	println("response qCode: ", inputData.QuestionCode)
	println("config qCode: ", questionCode )
	if inputData.QuestionCode == questionCode{
		return true
	}
	return false
}

func bothQAndDqFound(found isFound) bool{
	if found.qCodeIsFound && found.dqCodeIsFound{
		return true
	}
	return false
}
