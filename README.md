# AI Resume Screener Backend

A comprehensive AI-powered resume screening system built with Go, Gin, and Google's Gemini AI. This backend provides intelligent candidate-job matching, automated resume parsing, and advanced scoring algorithms to streamline the recruitment process.

## 🚀 Features

### Core Functionality

#### 📄 Resume Parsing Service
- **Multi-format Support**: Accepts PDF, DOCX, and TXT resume files
- **Intelligent Extraction**: Automatically extracts candidate information including:
  - Personal details (name, email, phone)
  - Education history
  - Work experience and roles
  - Skills and certifications
- **AI-Enhanced Skill Detection**: Uses Gemini AI to identify additional skills and competencies
- **File Validation**: Size limits (10MB) and format validation

#### 💼 Job Description Management
- **RESTful API**: Create and manage job postings
- **Flexible Requirements**: Define required and nice-to-have skills
- **Experience Levels**: Support for entry, mid, and senior level positions
- **Education Requirements**: Specify minimum education levels

#### 🤖 AI-Powered Matching Engine
- **Semantic Matching**: Uses Gemini AI for intelligent resume-job comparison
- **Multi-dimensional Scoring**: Combines traditional and AI-based scoring:
  - Required skills matching (40%)
  - Nice-to-have skills (20%)
  - Experience compatibility (30%)
  - Education requirements (10%)
- **Detailed Reasoning**: AI provides explanations for match scores
- **Skill Gap Analysis**: Identifies missing skills and competencies

#### 📊 Candidate Shortlisting
- **Automated Ranking**: Sort candidates by match score
- **Filtering Options**: Filter by minimum score, experience level, etc.
- **Batch Processing**: Process multiple resumes against job requirements
- **Score Persistence**: Store and retrieve historical match results

### Advanced Features

#### 🔐 Security & Middleware
- **JWT Authentication**: Secure API access with token-based authentication
- **CORS Support**: Cross-origin resource sharing for web applications
- **Rate Limiting**: Prevent abuse with configurable request limits (100/minute)
- **Request Logging**: Comprehensive logging for debugging and monitoring
- **Error Handling**: Robust error responses and panic recovery

#### 🗄️ Database Integration
- **PostgreSQL Support**: Full relational database integration
- **Redis Caching**: High-performance caching for frequently accessed data
- **Auto-migration**: Automatic database schema updates
- **GORM ORM**: Type-safe database operations

#### 📈 Performance & Scalability
- **Concurrent Processing**: Handle multiple requests simultaneously
- **Efficient Parsing**: Optimized file processing and text extraction
- **Database Indexing**: Optimized queries for large datasets
- **Memory Management**: Efficient resource utilization

## 🛠️ Technology Stack

- **Backend Framework**: Gin (Go)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **AI Integration**: Google Gemini AI
- **File Processing**: unidoc/pdf for PDF parsing
- **Authentication**: JWT tokens
- **Configuration**: Viper for environment management

## 📋 API Endpoints

### User Management
```
POST /user/register - Register a new user
```

### Resume Management
```
POST /resume/upload - Upload and parse resume file
```

### Job Management
```
POST /job/create - Create a new job description
GET  /job/list - List all job descriptions
POST /job/match/:jobId - Match candidates to a specific job
GET  /job/top/:jobId - Get top candidates for a job
```

### System
```
GET / - Health check and system information
```

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or higher
- PostgreSQL database
- Redis server
- Google Gemini API key

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/your-username/ai-resume-screener.git
cd ai-resume-screener/backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
Create a `.env` file in the backend directory:
```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/resume_screener?sslmode=disable
REDIS_URL=redis://localhost:6379
AI_API_KEY=your_gemini_api_key_here
JWT_SECRET=your_jwt_secret_here
```

4. **Run database migrations**
```bash
go run cmd/api/main.go
```
The application will automatically migrate the database schema on startup.

5. **Start the server**
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## 📖 Usage Examples

### Upload a Resume
```bash
curl -X POST -F "resume=@resume.pdf" http://localhost:8080/resume/upload
```

### Create a Job Description
```bash
curl -X POST http://localhost:8080/job/create \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Software Engineer",
    "description": "We are looking for an experienced software engineer...",
    "required_skills": ["Go", "PostgreSQL", "Docker"],
    "nice_to_have_skills": ["Kubernetes", "AWS"],
    "experience_level": "senior",
    "min_experience": 5,
    "education_required": "bachelor"
  }'
```

### Match Candidates
```bash
curl -X POST http://localhost:8080/job/match/job-uuid-here
```

## 🔧 Configuration

### Environment Variables
- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection URL
- `AI_API_KEY`: Google Gemini API key
- `JWT_SECRET`: Secret key for JWT token generation

### File Upload Configuration
- **Max file size**: 10MB
- **Supported formats**: PDF, DOCX, TXT
- **Upload directory**: `./uploads/`

## 🤖 AI Integration

The system integrates with Google's Gemini AI for:
- **Resume Analysis**: Extract skills and competencies from unstructured text
- **Semantic Matching**: Understand context and meaning in job descriptions
- **Intelligent Scoring**: Provide nuanced match scores beyond keyword matching
- **Reasoning**: Explain why candidates are a good or poor fit

### AI Features
- **Skill Extraction**: Identify technical and soft skills
- **Experience Analysis**: Understand career progression and expertise
- **Education Matching**: Compare education levels and relevance
- **Cultural Fit**: Assess alignment with company values

## 📊 Scoring Algorithm

The matching score is calculated using a weighted algorithm:

```
Total Score = (Required Skills × 0.4) + (Nice-to-Have Skills × 0.2) + (Experience × 0.3) + (Education × 0.1)
```

When AI is enabled, the final score combines traditional matching (60%) with AI analysis (40%).

## 🔒 Security Features

- **Input Validation**: All API inputs are validated and sanitized
- **Rate Limiting**: Prevents abuse with configurable limits
- **CORS Protection**: Configurable cross-origin policies
- **Error Handling**: Secure error responses without information leakage
- **File Upload Security**: Strict file type and size validation

## 📈 Performance

- **Concurrent Requests**: Handles multiple simultaneous uploads and processing
- **Efficient Parsing**: Optimized text extraction from various file formats
- **Database Optimization**: Indexed queries for fast candidate retrieval
- **Caching**: Redis-based caching for improved response times

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

## 📝 Development

### Project Structure
```
backend/
├── cmd/api/           # Application entry point
├── internals/
│   ├── api/
│   │   ├── controller/    # HTTP request handlers
│   │   ├── middlewares/   # Custom middleware
│   │   └── routes/        # Route definitions
│   ├── config/        # Configuration management
│   ├── database/      # Database connection and setup
│   ├── models/        # Data models
│   └── services/      # Business logic services
├── uploads/           # File upload directory
└── go.mod             # Go module dependencies
```

### Adding New Features
1. Define data models in `internals/models/`
2. Implement business logic in `internals/services/`
3. Create API handlers in `internals/api/controller/`
4. Define routes in `internals/api/routes/`
5. Add middleware if needed in `internals/api/middlewares/`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the GitHub repository
- Check the documentation for common solutions
- Review the API endpoints and examples

## 🚀 Future Enhancements

- [ ] Frontend web application
- [ ] Advanced analytics and reporting
- [ ] Integration with ATS systems
- [ ] Real-time notifications
- [ ] Bulk resume processing
