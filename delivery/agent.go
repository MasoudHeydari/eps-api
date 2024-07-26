package delivery

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/MasoudHeydari/eps-api/model"
	"github.com/google/uuid"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Agent struct {
	phoneCollector                  *colly.Collector
	keywordCollector                *colly.Collector
	findPhoneRegexp, findEmailRegex *regexp.Regexp
	client                          *http.Client
}

func NewAgent() *Agent {
	phoneCollector := colly.NewCollector()
	phoneCollector.AllowURLRevisit = true
	phoneCollector.SetRequestTimeout(20 * time.Second)

	keywordCollector := colly.NewCollector()
	keywordCollector.AllowURLRevisit = true
	keywordCollector.SetRequestTimeout(20 * time.Second)
	return &Agent{
		phoneCollector:   phoneCollector,
		keywordCollector: keywordCollector,
		findPhoneRegexp:  regexp.MustCompile(`(?:[-+() ]*\d){10,13}`),
		findEmailRegex:   regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
		client:           &http.Client{Timeout: 20 * time.Second},
	}
}

func (a *Agent) GetRateLimiter() *rate.Limiter {
	// TODO: Complete this function
	return nil
}

func (a *Agent) CreateJob(query model.Query) (uuid.UUID, error) {
	logrus.Tracef("Start Google SERP API, query: %+v", query)
	requestBody := []map[string]interface{}{
		{
			"keyword":       query.Text,
			"location_code": query.Location,
			"language_code": query.LangCode,
			"depth":         query.Depth,
		},
	}
	var (
		apiResponse model.APIResponse
		jobID       uuid.UUID
	)

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return jobID, fmt.Errorf("CreateJob.json.Marshal: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.dataforseo.com/v3/serp/google/organic/task_post",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return jobID, fmt.Errorf("CreateJob.http.NewRequest: %w", err)
	}

	// Set the headers
	token := base64.StdEncoding.EncodeToString([]byte("m.heydari4883@gmail.com:22599da38215faea"))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return jobID, fmt.Errorf("CreateJob.a.client.Do: %w", err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return jobID, fmt.Errorf("CreateJob.json.NewDecoder.Decode: %w", err)
	}
	jobID, err = getJobID(apiResponse)
	if err != nil {
		return jobID, fmt.Errorf("CreateJob: %w", err)
	}
	logrus.Println("*********************************************************")
	logrus.Println("*********************************************************")
	logrus.Infof("agent.CreateJob: job created successfully for query '%+v' and job id is: %s\n", query, jobID)
	logrus.Infof("agent.CreateJob: Total API Cost for job id %s is $%f", jobID, apiResponse.Cost)
	logrus.Println("*********************************************************")
	logrus.Println("*********************************************************")
	return jobID, nil
}

func (a *Agent) PollJob(jobID uuid.UUID) ([]model.Item, bool, error) {
	u := fmt.Sprintf(
		"https://api.dataforseo.com/v3/serp/google/organic/task_get/regular/%s",
		jobID,
	)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, false, fmt.Errorf("PollJob.http.NewRequest: %w", err)
	}

	// Set the headers
	token := base64.StdEncoding.EncodeToString([]byte("m.heydari4883@gmail.com:22599da38215faea"))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("PollJob.a.client.Do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("PollJob.client.Do.StatusCode: unintended staus code %d", resp.StatusCode)
	}

	var apiResponse model.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, false, fmt.Errorf("PollJob.json.NewDecoder.Decode: %w", err)
	}

	if len(apiResponse.Tasks) == 0 {
		return nil, false, fmt.Errorf("PollJob: len(apiResponse.Tasks) == 0,  no item found for job id %s", jobID)
	}
	if apiResponse.Tasks[0].StatusCode == 40102 {
		// means this job id has not any result
		return nil, true, fmt.Errorf("PollJob: there is no result for job id %s - status code is 40102", jobID)
	}
	if len(apiResponse.Tasks[0].Result) == 0 {
		return nil, false, fmt.Errorf("PollJob: len(apiResponse.Tasks[0].Result) == 0, no item found for job id %s", jobID)
	}
	if len(apiResponse.Tasks[0].Result[0].Items) == 0 {
		return nil, false, fmt.Errorf("PollJob: len(apiResponse.Tasks[0].Result[0].Items) == 0, no item found for job id %s", jobID)
	}
	logrus.Println("*********************************************************")
	logrus.Info("PollJob: Total result found: ", len(apiResponse.Tasks[0].Result[0].Items))
	logrus.Println("*********************************************************")
	return apiResponse.Tasks[0].Result[0].Items, false, nil
}

func (a *Agent) extractKeywords(path string) ([]string, error) {
	var (
		maxLen               = 6
		maxLenForEachKeyword = 72
		h1h2h3QuerySelector  = "h1, h2, h3"
	)
	keyWordsMap := make(map[string]struct{})
	keyWords := make([]string, 0)
	defer func() {
		a.keywordCollector.OnHTMLDetach(h1h2h3QuerySelector)
	}()
	a.keywordCollector.OnHTML(h1h2h3QuerySelector, func(e *colly.HTMLElement) {
		keyWordsMap[e.Text] = struct{}{}
	})
	err := a.keywordCollector.Visit(path)
	if err != nil {
		return keyWords, fmt.Errorf("extractKeywords: %w", err)
	}
	if len(keyWordsMap) == 0 {
		return keyWords, fmt.Errorf("no keyword found from %q", path)
	}
	for k := range keyWordsMap {
		k = strings.TrimSpace(k)
		k = strings.ReplaceAll(k, "\t", "")
		k = strings.ReplaceAll(k, "\n", "")
		if len(k) > maxLenForEachKeyword {
			k = k[:maxLenForEachKeyword]
		}
		keyWords = append(keyWords, strings.TrimSpace(k))
	}
	if len(keyWords) > maxLen {
		keyWords = keyWords[:maxLen]
	}
	return keyWords, nil
}

func (a *Agent) extractEmails(path string) ([]string, error) {
	emailsMap := make(map[string]struct{})
	emails := make([]string, 0)
	resp, err := a.client.Get(path)
	if err != nil {
		return emails, fmt.Errorf("extractEmails.Get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return emails, fmt.Errorf("extractEmails.ReadAll: %w", err)
	}
	emailsArray := a.findEmailRegex.FindAllString(string(body), -1)
	for _, email := range emailsArray {
		emailsMap[email] = struct{}{}
	}
	if len(emailsMap) == 0 {
		return nil, fmt.Errorf("no email found from %q", path)
	}
	for k := range emailsMap {
		emails = append(emails, k)
	}
	return emails, nil
}

func (a *Agent) extractPhoneNumbers(path string) ([]string, error) {
	matches := make([]string, 0)
	divQuerySelector := "div"
	defer func() {
		a.phoneCollector.OnHTMLDetach(divQuerySelector)
	}()
	a.phoneCollector.OnHTML(divQuerySelector, func(e *colly.HTMLElement) {
		newMatches := a.findPhoneRegexp.FindAllString(e.Text, -1)
		matches = append(
			matches,
			newMatches...,
		)
	})
	err := a.phoneCollector.Visit(path)
	if err != nil {
		return nil, fmt.Errorf("extractPhoneNumbers: %w", err)
	}
	return matches, nil
}

func (a *Agent) extractPhoneNumbersFromAllPossibleURLs(p string) ([]string, error) {
	phoneNums := make([]string, 0)
	paths := make(map[string]struct{}, 0)
	u, err := url.Parse(p)
	if err != nil {
		return phoneNums, fmt.Errorf("extractPhoneNumbersFromAllPossibleURLs.Parse: %w", err)
	}
	paths[p] = struct{}{}
	paths[u.Scheme+"://"+u.Host+"/about-us"] = struct{}{}
	paths[u.Scheme+"://"+u.Host+"/about"] = struct{}{}
	paths[u.Scheme+"://"+u.Host+"/contact-us"] = struct{}{}
	paths[u.Scheme+"://"+u.Host+"/contact"] = struct{}{}

	for pp := range paths {
		phones, err := a.extractPhoneNumbers(pp)
		if err != nil {
			logrus.Infof("extractPhoneNumbersFromAllPossibleURLs: %v", err)
		}
		phoneNums = append(phoneNums, phones...)
	}
	if len(phoneNums) == 0 {
		return phoneNums, fmt.Errorf("no phone number found from %q", p)
	}
	phonesMap := make(map[string]struct{})
	for _, phone := range phoneNums {
		phone = strings.TrimSpace(phone)
		phonesMap[phone] = struct{}{}
	}
	results := make([]string, 0, len(phonesMap))
	for k := range phonesMap {
		results = append(results, k)
	}
	return results, nil
}

func getJobID(response model.APIResponse) (uuid.UUID, error) {
	if response.TasksCount == 0 {
		return uuid.UUID{}, fmt.Errorf("getJobID: tasks count is zero")
	}
	return uuid.Parse(response.Tasks[0].ID)
}
