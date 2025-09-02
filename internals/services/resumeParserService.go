package services

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type ResumeParserService struct{}

func NewResumeParserService() *ResumeParserService {
	return &ResumeParserService{}
}

func (r *ResumeParserService) ParseResume(file multipart.File, header *multipart.FileHeader) (*models.Resume, error) {
	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var parsedText string
	fileExt := strings.ToLower(filepath.Ext(header.Filename))

	switch fileExt {
	case ".pdf":
		parsedText, err = r.parsePDF(content)
	case ".docx":
		parsedText, err = r.parseDOCX(content)
	case ".txt":
		parsedText = string(content)
	default:
		return nil, errors.New("unsupported file type")
	}

	if err != nil {
		return nil, err
	}

	// Extract structured data from text
	resume := r.extractResumeData(parsedText)
	resume.ParsedText = parsedText
	resume.FilePath = header.Filename // In production, save to disk or cloud

	return resume, nil
}

func (r *ResumeParserService) parsePDF(content []byte) (string, error) {
	reader := bytes.NewReader(content)
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return "", err
	}

	var text strings.Builder
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			continue
		}

		ex, err := extractor.New(page)
		if err != nil {
			continue
		}

		pageText, err := ex.ExtractText()
		if err != nil {
			continue
		}

		text.WriteString(pageText)
		text.WriteString("\n")
	}

	return text.String(), nil
}

func (r *ResumeParserService) parseDOCX(content []byte) (string, error) {
	return "", errors.New("DOCX parsing is not fully implemented yet. Please use PDF or TXT files for now")
}

func (r *ResumeParserService) extractResumeData(text string) *models.Resume {
	resume := &models.Resume{}

	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	if email := emailRegex.FindString(text); email != "" {
		resume.Email = email
	}

	phoneRegex := regexp.MustCompile(`(\+?\d{1,3}[-.\s]?)?\(?(\d{3})\)?[-.\s]?(\d{3})[-.\s]?(\d{4})`)
	if phone := phoneRegex.FindString(text); phone != "" {
		resume.Phone = phone
	}

	skillKeywords := []string{"python", "java", "javascript", "go", "c++", "c#", "php", "ruby", "swift", "kotlin",
		"react", "angular", "vue", "node", "django", "flask", "spring", "laravel", "express",
		"mysql", "postgresql", "mongodb", "redis", "docker", "kubernetes", "aws", "azure", "gcp",
		"git", "linux", "windows", "agile", "scrum", "ci/cd", "jenkins", "github", "gitlab"}

	var skills []string
	lowerText := strings.ToLower(text)
	for _, skill := range skillKeywords {
		if strings.Contains(lowerText, skill) {
			skills = append(skills, strings.Title(skill))
		}
	}
	resume.Skills = skills

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) > 3 && len(line) < 50 {
			nameRegex := regexp.MustCompile(`^[a-zA-Z\s.]{3,50}$`)
			if nameRegex.MatchString(line) && !strings.Contains(strings.ToLower(line), "email") &&
			   !strings.Contains(strings.ToLower(line), "phone") && !strings.Contains(line, "@") {
				resume.CandidateName = line
				break
			}
		}
	}

	return resume
}
