package main

type ContributorResponses struct {
	Reference string `json:"reference"`
	Period    string `json:"period"`
	Survey    string `json:"survey"`
	BpmID     string `json:"instance"`

	Responses []responseData `json:"responses"`
}

type responseData struct {
	QuestionCode string `json:"questionCode"`
	Response     string `json:"response"`
}

type validationConfig struct {
	questionCode        string
	derivedQuestionCode string
}

type isFound struct {
	qCodeIsFound  bool
	dqCodeIsFound bool
}

type sendStruct struct {
	Reference string `json:"reference"`
	Period    string `json:"period"`
	Survey    string `json:"survey"`
	BpmID     string `json:"instance"`
}

type validationRequest struct {
	PrimaryValue    string `json:"primaryValue"`
	ComparisonValue string `json:"comparisonValue"`
	MetaData        string `json:"metaData"`
}

type validationResult struct {
	QuestionCode string
	ValueFormula string
	Triggered    string
}

type VetQandResponse struct {
	FinalQuestCode         string
	FinalDerivedQuestCode  string
	FinalQuestCodeValue    string
	FinalDerivedQuestValue string
	IsQuestionCodeFound    bool
	IsDerivedQuestFound    bool
	ValidationResult       validationResult
}
