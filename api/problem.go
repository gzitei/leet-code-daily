package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type symbol string

const (
	easy     symbol = "\U0001F7E2" //🟢
	medium          = "\U0001F534" //🟡
	hard            = "\U0001F534" //🔴
	bookmark        = "\U0001F516" //🔖
	pin             = "\U0001F4CC" //📌
	calendar        = "\U0001F4C5" //📆
	chart           = "\U0001F4CA" //📊
	globe           = "\U0001F310" //🌐
	lamp            = "\U0001F4A1" //💡

)

var dir, rootDir string

var status map[string]string = map[string]string{
	"Easy":   string(easy),
	"Medium": string(medium),
	"Hard":   string(hard),
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func WriteFile(fname, content string) {
	fp := filepath.Join(rootDir, dir, fname)
	if FileExists(fp) {
		return
	}
	os.WriteFile(fp, []byte(content), 0o666)
}

func CreateFolder() (string, error) {
	dirPath := filepath.Join(rootDir, dir)
	if FileExists(dirPath) {
		return dirPath, nil
	}
	if err := os.Mkdir(dirPath, 0o777); err != nil {
		return "", err
	}
	return dirPath, nil
}

func createMarkdown(question Question) {
	converter := md.NewConverter("", true, nil)
	tagList := []string{}
	for _, v := range question.TopicTags {
		tagList = append(tagList, "#"+v.Slug)
	}
	link := fmt.Sprintf("https://leetcode.com/problems/%s", question.TitleSlug)
	title := fmt.Sprintf("%s. %s", question.FrontendQuestionID, question.Title)
	if markdown, err := converter.ConvertString(question.Content); err == nil {
		markdown = strings.ReplaceAll(markdown, "(../", "("+link+"/")
		tags := "```yaml\n" + bookmark + "tags: " + strings.Join(tagList, " ") + "\n```"
		s := fmt.Sprintf("# %s\n##### %s %s | %s %s | %s %s | %s %.2f%% | %s [Leet-Code #%s](%s)\n---\n%s\n---\n",
			title,
			pin,
			question.CategoryTitle,
			calendar,
			time.Now().Format("2006-01-02"),
			status[question.Difficulty],
			question.Difficulty,
			chart,
			question.AcRate,
			globe,
			question.FrontendQuestionID,
			link,
			tags,
		)
		s += markdown
		s += "\n\n"
		if len(question.SimilarQuestionList) > 0 {
			s += "> ### Questões Similares:"
		}
		for _, similar := range question.SimilarQuestionList {
			similarTopics := []string{}
			for _, topic := range similar.TopicTags {
				similarTopics = append(similarTopics, "#"+topic.Slug)
			}
			s += fmt.Sprintf("\n> %s [%s. %s](https://leetcode.com/problems/%s) %s",
				status[similar.Difficulty],
				similar.FrontendQuestionID,
				similar.Title,
				similar.TitleSlug,
				strings.Join(similarTopics, " "),
			)
		}
		s += "\n---\n"
		if len(question.Hints) > 0 {
			s += "> ### Dicas:"
		}
		for _, h := range question.Hints {
			s += "\n>" + lamp + h
		}
		WriteFile("README.md", s)
	} else {
		fmt.Println(err)
	}
	if question.HasSolution {
		if markdown, err := converter.ConvertString(question.Solution.Body); err == nil {
			markdown = strings.ReplaceAll(markdown, "(../", "("+link+"/")
			WriteFile(".solution.md", markdown)
		}
	}
}

func CreateCodingFile(q Question) {
	test := strings.Join(q.ExampleTestcaseList, ", ")
	snippets := q.CodeSnippets
	for _, snippet := range snippets {
		if snippet.LangSlug == "golang" {
			WriteFile("main.go", fmt.Sprintf("package main\n\n%s\n\nfunc main() {\n\t/*\n\t\t%s\n\t*/\n}", snippet.Code, test))
		}
	}
}

func SetUpEnv(d string, question Question) (string, error) {
	rootDir = d
	dir = question.FrontendQuestionID
	created, err := CreateFolder()
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return "", err
	}
	createMarkdown(question)
	CreateCodingFile(question)
	return created, nil
}
