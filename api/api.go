package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func sendHttpRequest(payload string) string {
	endpoint := "https://leetcode.com/graphql/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client := http.Client{}

	ctt := bytes.NewBuffer([]byte(payload))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, ctt)
	if err != nil {
		fmt.Println("Error on request setup:", err)
		return ""
	}

	req.Header.Add("accept-language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,es;q=0.6")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://leetcode.com")
	req.Header.Add("referer", "https://leetcode.com")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error after request sent:", err)
		return ""
	}

	defer res.Body.Close()
	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error while parsing response body:", err)
		return ""
	}
	return string(bodyContent)
}

func GetDailyChallenge() (Question, error) {
	query := map[string]interface{}{
		"query": `query
            questionOfToday {
                activeDailyCodingChallengeQuestion {
                    date
                    userStatus
                    link
                    question {
                        companyTags {
                            name
                            slug
                        }
                        solution {
                            body
                        }
                        acRate
                        content
                        categoryTitle
                        hints
                        exampleTestcaseList
                        difficulty
                        freqBar
                        frontendQuestionId: questionFrontendId
                        isFavor
                        paidOnly: isPaidOnly
                        status
                        title
                        titleSlug
                        hasVideoSolution
                        hasSolution
                        topicTags {
                            slug
                        }
                        codeSnippets {
                            lang
                            langSlug
                            code
                        }
                    }
                }
            }`,
		"variables":     map[string]interface{}{},
		"operationName": "questionOfToday",
	}
	var content string
	if payload, err := json.Marshal(query); err != nil {
		fmt.Println("Error while marshalling payload json:", err)
		return Question{}, err
	} else {
		content = sendHttpRequest(string(payload))
	}
	var question DailyQuestion
	if err := json.Unmarshal([]byte(content), &question); err != nil {
		fmt.Println("Error while unmarshalling response body:", err)
		return Question{}, err
	}
	return question.Data.ActiveDailyCodingChallengeQuestion.Question, nil
}

func GetRandomQuestion() (Question, error) {
	query := map[string]interface{}{
		"query": `query randomQuestion($categorySlug: String, $filters: QuestionListFilterInput) {
            randomQuestion(categorySlug: $categorySlug, filters: $filters) {
                companyTags {
                    name
                    slug
                }
                solution {
                    body
                }
                acRate
                content
                categoryTitle
                hints
                exampleTestcaseList
                difficulty
                freqBar
                frontendQuestionId: questionFrontendId
                isFavor
                paidOnly: isPaidOnly
                status
                title
                titleSlug
                hasVideoSolution
                hasSolution
                topicTags {
                    slug
                }
                codeSnippets {
                    lang
                    langSlug
                    code
                }
            }
        }`,
		"variables": map[string]interface{}{
			"categorySlug": "", "filters": map[string]interface{}{},
		},
	}
	var content string
	if payload, err := json.Marshal(query); err != nil {
		fmt.Println("Error while marshalling payload json:", err)
		return Question{}, err
	} else {
		content = sendHttpRequest(string(payload))
	}
	var question RandomQuestionT
	if err := json.Unmarshal([]byte(content), &question); err != nil {
		fmt.Println("Error while unmarshalling response body:", err)
		return Question{}, err
	}
	return question.Data.RandomQuestion, nil
}

func GetQuestionBySlug(slug string) (Question, error) {
	query := map[string]interface{}{
		"query": `query questionContent($titleSlug: String!) {
            question(titleSlug: $titleSlug) {
                companyTags {
                    name
                    slug
                }
                solution {
                    body
                }
                acRate
                content
                categoryTitle
                hints
                exampleTestcaseList
                difficulty
                freqBar
                frontendQuestionId: questionFrontendId
                isFavor
                paidOnly: isPaidOnly
                status
                title
                titleSlug
                hasVideoSolution
                hasSolution
                topicTags {
                    slug
                }
                codeSnippets {
                    lang
                    langSlug
                    code
                }
            }
        }`,
		"variables": map[string]interface{}{
			"titleSlug": slug,
		},
	}
	var content string
	if payload, err := json.Marshal(query); err != nil {
		fmt.Println("Error while marshalling payload json:", err)
		return Question{}, err
	} else {
		content = sendHttpRequest(string(payload))
	}
	var question QuestionBySlug
	if err := json.Unmarshal([]byte(content), &question); err != nil {
		fmt.Println("Error while unmarshalling response body:", err)
		return Question{}, err
	}
	return question.Data.Question, nil
}
