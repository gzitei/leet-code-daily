package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type Challenge struct {
	Data struct {
		ActiveDailyCodingChallengeQuestion struct {
			Date     string `json:"date"`
			Link     string `json:"link"`
			Question struct {
				AcRate        float64 `json:"acRate"`
				CategoryTitle string  `json:"categoryTitle"`
				CodeSnippets  []struct {
					Code     string `json:"code"`
					Lang     string `json:"lang"`
					LangSlug string `json:"langSlug"`
				} `json:"codeSnippets"`
				CompanyTags        interface{}   `json:"companyTags"`
				Content            string        `json:"content"`
				DataSchemas        []interface{} `json:"dataSchemas"`
				Difficulty         string        `json:"difficulty"`
				FreqBar            interface{}   `json:"freqBar"`
				FrontendQuestionID string        `json:"frontendQuestionId"`
				HasSolution        bool          `json:"hasSolution"`
				HasVideoSolution   bool          `json:"hasVideoSolution"`
				Hints              []string      `json:"hints"`
				IsFavor            bool          `json:"isFavor"`
				PaidOnly           bool          `json:"paidOnly"`
				Solution           struct {
					Body string `json:"body"`
				} `json:"solution"`
				Status    interface{} `json:"status"`
				Title     string      `json:"title"`
				TitleSlug string      `json:"titleSlug"`
				TopicTags []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Slug string `json:"slug"`
				} `json:"topicTags"`
			} `json:"question"`
			UserStatus string `json:"userStatus"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

func getDailyCodingChallenge() (Challenge, error) {
	url := "https://leetcode.com/graphql/"
	method := "POST"
	payload := strings.NewReader(`{
    "query": "\n    query questionOfToday {\n  activeDailyCodingChallengeQuestion {\n    date\n    userStatus\n    link\n    question {\n      codeSnippets {\n      lang\n      langSlug\n      code\n    }\n    dataSchemas\n      companyTags {name, slug\n}\n      solution {body}\n    acRate\n      content\n      categoryTitle\n      difficulty\n      freqBar\n      frontendQuestionId: questionFrontendId\n      isFavor\n      hints\n      paidOnly: isPaidOnly\n      status\n      title\n      titleSlug\n      hasVideoSolution\n      hasSolution\n      topicTags {\n        name\n        id\n        slug\n      }\n    }\n  }\n}\n    ",
    "variables": {},
    "operationName": "questionOfToday"
}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return Challenge{}, err
	}
	req.Header.Add("accept-language", "pt-BR,pt;q=0.9")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("origin", "https://leetcode.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Challenge{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return Challenge{}, err
	}
	var content Challenge
	err = json.Unmarshal(body, &content)
	if err != nil {
		fmt.Println(err)
	}
	return content, err
}

func WriteFile(fname, content string) {
	fp := filepath.Join(rootDir, dir, fname)
	os.WriteFile(fp, []byte(content), 0o666)
}

var status map[string]string = map[string]string{
	"Easy":   "🟢",
	"Medium": "🟡",
	"Hard":   "🔴",
}

func CreateFolder() string {
	dirPath := filepath.Join(rootDir, dir)
	if err := os.Mkdir(dirPath, 0o777); err != nil {
		fmt.Println(err)
	}
	return dirPath
}

func goToFolder(str string) {
	cmd := exec.Command("/usr/bin/tmux", "new-session", "-d", "-s", dir, "bash", "-c", "cd "+str+" && nvim -c 'Alpha' "+str)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	cmd = exec.Command("gnome-terminal", "--", "bash", "-c", "tmux attach -t "+dir)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

var rootDir, dir string

func main() {
	rootDir = os.Args[1]
	challenge, err := getDailyCodingChallenge()
	if err != nil {
		fmt.Println(err)
		return
	}
	question := challenge.Data.ActiveDailyCodingChallengeQuestion.Question
	converter := md.NewConverter("", true, nil)
	tagList := []string{}
	for _, v := range question.TopicTags {
		tagList = append(tagList, "#"+v.Slug)
	}
	dir = question.FrontendQuestionID
	created := CreateFolder()
	link := fmt.Sprintf("https://leetcode.com/problems/%s", question.TitleSlug)
	title := fmt.Sprintf("%s. %s", question.FrontendQuestionID, question.Title)
	if markdown, err := converter.ConvertString(question.Content); err == nil {
		markdown = strings.ReplaceAll(markdown, "(../", "("+link+"/")
		tags := "```yaml\n🔖tags: " + strings.Join(tagList, " ") + "\n```"
		s := fmt.Sprintf("# %s\n##### 📌 %s | 📆 %s | [%s %s]() | 📊 %.2f%% | 🌐 [Leet-Code #%s](%s)\n%s\n---\n",
			title,
			question.CategoryTitle,
			challenge.Data.ActiveDailyCodingChallengeQuestion.Date,
			status[question.Difficulty],
			question.Difficulty,
			question.AcRate,
			question.FrontendQuestionID,
			link,
			tags,
		)
		s += markdown
		s += "\n"
		for _, h := range question.Hints {
			s += "\n>💡" + h
		}
		WriteFile("challenge.md", s)
	} else {
		fmt.Println(err)
	}
	if markdown, err := converter.ConvertString(question.Solution.Body); err == nil {
		markdown = strings.ReplaceAll(markdown, "(../", "("+link+"/")
		WriteFile(".solution.md", markdown)
	}
	for _, snp := range question.CodeSnippets {
		switch snp.Lang {
		case "Python3":
			WriteFile("solution.py", snp.Code)
		case "TypeScript":
			WriteFile("solution.ts", snp.Code)
		case "Go":
			WriteFile("main.go", "package main\n\n"+snp.Code+"\n\nfunc main() {\n}")
		}
	}
	goToFolder(created)
}
