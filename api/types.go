package api

type Question struct {
	AcRate        float64 `json:"acRate"`
	CategoryTitle string  `json:"categoryTitle"`
	CodeSnippets  []struct {
		Code     string `json:"code"`
		Lang     string `json:"lang"`
		LangSlug string `json:"langSlug"`
	} `json:"codeSnippets"`
	CompanyTags         interface{}   `json:"companyTags"`
	Content             string        `json:"content"`
	Difficulty          string        `json:"difficulty"`
	ExampleTestcaseList []string      `json:"exampleTestcaseList"`
	FreqBar             interface{}   `json:"freqBar"`
	FrontendQuestionID  string        `json:"frontendQuestionId"`
	HasSolution         bool          `json:"hasSolution"`
	HasVideoSolution    bool          `json:"hasVideoSolution"`
	Hints               []string      `json:"hints"`
	IsFavor             bool          `json:"isFavor"`
	NextChallenges      []interface{} `json:"nextChallenges"`
	PaidOnly            bool          `json:"paidOnly"`
	SimilarQuestionList []struct {
		Difficulty         string `json:"difficulty"`
		FrontendQuestionID string `json:"frontendQuestionId"`
		MetaData           string `json:"metaData"`
		PaidOnly           bool   `json:"paidOnly"`
		Title              string `json:"title"`
		TitleSlug          string `json:"titleSlug"`
		TopicTags          []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"topicTags"`
	} `json:"similarQuestionList"`
	Solution struct {
		Body string `json:"body"`
	} `json:"solution"`
	Status    interface{} `json:"status"`
	Title     string      `json:"title"`
	TitleSlug string      `json:"titleSlug"`
	TopicTags []struct {
		Slug string `json:"slug"`
	} `json:"topicTags"`
}
type QuestionBySlug struct {
	Data struct {
		Question Question `json:"question"`
	} `json:"data"`
}

type RandomQuestionT struct {
	Data struct {
		RandomQuestion Question `json:"randomQuestion"`
	} `json:"data"`
}

type DailyQuestion struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion struct {
			Question   Question
			Date       string `json:"date"`
			Link       string `json:"link"`
			UserStatus string `json:"userStatus"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}
