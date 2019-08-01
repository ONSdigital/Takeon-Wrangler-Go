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

type validationRequest struct {
	PrimaryValue    string `json:"primaryValue"`
	ComparisonValue string `json:"comparisonValue"`
	MetaData        string `json:"metaData"`
}

type validationResult struct {
	ValueFormula string `json:"valueFormula"`
	Triggered    bool   `json:"triggered"`
	MetaData     string `json:"metaData"`
}

type VetQResponse struct {
	QuestionCode string `json:"questionCode"`
	ValueFormula string `json:"valueFormula"`
	Triggered    bool   `json:"triggered"`
}

type VetQDqResponse struct {
	FinalQuestCode         string
	FinalDerivedQuestCode  string
	FinalQuestCodeValue    string
	FinalDerivedQuestValue string
	IsQuestionCodeFound    bool
	IsDerivedQuestFound    bool
	ValidationResult       validationResult
}

type PersistInpData struct {
	Reference         string         `json:"reference"`
	Period            string         `json:"period"`
	Survey            string         `json:"survey"`
	Instance          string         `json:"instance"`
	ValidationName    string         `json:"validationName"`
	ValidationResults []VetQResponse `json:"validationResults"`
}

type PersistPostReq struct {
	Type  string         `json:"type"`
	Input PersistInpData `json:"input"`
}

type isFound struct {
	qCodeIsFound  bool
	dqCodeIsFound bool
}
