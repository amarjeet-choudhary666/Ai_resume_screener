package services

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	client *genai.Client
}

type AIMatchResult struct {
	Score           float64  `json:"score"`
	Reasoning       string   `json:"reasoning"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	ExperienceMatch float64  `json:"experience_match"`
	EducationMatch  float64  `json:"education_match"`
}

type AISkillExtraction struct {
	Skills         []string `json:"skills"`
	ExperienceYears int      `json:"experience_years"`
	EducationLevel  string   `json:"education_level"`
}

func NewAIService(apiKey string) (*AIService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &AIService{client: client}, nil
}

func (a *AIService) Close() {
	a.client.Close()
}

func (a *AIService) EnhanceMatching(resumeText, jobDescription string) (*AIMatchResult, error) {
	ctx := context.Background()
	model := a.client.GenerativeModel("gemini-1.5-flash")

	prompt := fmt.Sprintf(`Analyze the following resume and job description for comprehensive matching.
Return a JSON response with the following structure:
{
  "score": <number 0-100>,
  "reasoning": "<detailed explanation>",
  "matched_skills": ["skill1", "skill2"],
  "missing_skills": ["skill3", "skill4"],
  "experience_match": <number 0-100>,
  "education_match": <number 0-100>
}

Resume:
%s

Job Description:
%s

Focus on:
- Technical skills matching
- Experience level compatibility
- Education requirements
- Cultural fit indicators
- Overall suitability`, resumeText, jobDescription)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Try to parse as JSON
	var result AIMatchResult
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		// If JSON parsing fails, extract score and reasoning manually
		return a.parseManualResponse(responseText)
	}

	return &result, nil
}

func (a *AIService) parseManualResponse(response string) (*AIMatchResult, error) {
	result := &AIMatchResult{}

	// Extract score using regex
	scoreRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)(?:/100|%)`)
	if match := scoreRegex.FindStringSubmatch(response); len(match) > 1 {
		if score, err := strconv.ParseFloat(match[1], 64); err == nil {
			result.Score = score
		}
	}

	// Extract reasoning (everything after score)
	parts := strings.Split(response, "\n")
	result.Reasoning = strings.Join(parts[1:], "\n")

	return result, nil
}

func (a *AIService) ExtractSkillsFromText(text string) (*AISkillExtraction, error) {
	ctx := context.Background()
	model := a.client.GenerativeModel("gemini-1.5-flash")

	prompt := fmt.Sprintf(`Extract technical skills, experience level, and education from the following text.
Return a JSON response with the following structure:
{
  "skills": ["skill1", "skill2", "skill3"],
  "experience_years": <number>,
  "education_level": "<bachelor|master|phd|associate|high_school>"
}

Text:
%s

Focus on:
- Programming languages and frameworks
- Tools and technologies
- Soft skills
- Years of experience
- Highest education level`, text)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	var result AISkillExtraction
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		// Fallback to manual extraction
		return a.manualSkillExtraction(text)
	}

	return &result, nil
}

func (a *AIService) manualSkillExtraction(text string) (*AISkillExtraction, error) {
	result := &AISkillExtraction{}

	// Simple skill extraction based on common keywords
	skillKeywords := []string{
		"python", "java", "javascript", "go", "c++", "c#", "php", "ruby", "swift", "kotlin",
		"react", "angular", "vue", "node", "django", "flask", "spring", "laravel", "express",
		"mysql", "postgresql", "mongodb", "redis", "docker", "kubernetes", "aws", "azure", "gcp",
		"git", "linux", "windows", "agile", "scrum", "ci/cd", "jenkins", "github", "gitlab",
		"html", "css", "typescript", "graphql", "rest", "api", "microservices",
	}

	lowerText := strings.ToLower(text)
	for _, skill := range skillKeywords {
		if strings.Contains(lowerText, skill) {
			result.Skills = append(result.Skills, strings.Title(skill))
		}
	}

	// Extract experience years
	expRegex := regexp.MustCompile(`(\d+)\s*(?:year|yr)s?\s*(?:of\s*)?experience`)
	if match := expRegex.FindStringSubmatch(lowerText); len(match) > 1 {
		if years, err := strconv.Atoi(match[1]); err == nil {
			result.ExperienceYears = years
		}
	}

	// Extract education level
	educationLevels := map[string]string{
		"phd": "phd", "doctorate": "phd", "doctoral": "phd",
		"master": "master", "masters": "master", "msc": "master", "ms": "master",
		"bachelor": "bachelor", "bachelors": "bachelor", "bsc": "bachelor", "bs": "bachelor",
		"associate": "associate", "diploma": "associate",
	}

	for _, level := range educationLevels {
		if strings.Contains(lowerText, level) {
			result.EducationLevel = level
			break
		}
	}

	return result, nil
}

func (a *AIService) GenerateJobSummary(jobDescription string) (string, error) {
	ctx := context.Background()
	model := a.client.GenerativeModel("gemini-1.5-flash")

	prompt := fmt.Sprintf(`Create a concise summary of the following job description, highlighting:
- Key responsibilities
- Required skills and qualifications
- Experience level needed
- Company/industry context

Job Description:
%s

Keep the summary under 200 words.`, jobDescription)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
