package main

type ContributorResponses struct{
	Reference string `json:"reference"`
	Period string `json:"period"`
	Survey string `json:"survey"`
	BpmID string `json:"instance"`

	Responses []responseData `json:"responses"`
}

type responseData struct{
	QuestionCode string `json:"questionCode"`
	Response string `json:"response"`
}

type validationConfig struct{
	questionCode string
	derivedQuestionCode string
}

type isFound struct{
	qCodeIsFound bool
	dqCodeIsFound bool
}

type sendStruct struct{
	Reference string `json:"reference"`
	Period string `json:"period"`
	Survey string `json:"survey"`
	BpmID string `json:"instance"`


}