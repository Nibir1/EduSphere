# EduSphere: AI-Powered Academic Assistant

**EduSphere** is a full-stack **Generative AI academic assistant** designed to analyze student transcripts, generate personalized course and scholarship recommendations, and create downloadable academic reports — all locally and securely.  

Built with a **production-grade Golang Fiber backend**, **React (Vite) frontend**, and **Ollama AI inference engine**, EduSphere demonstrates how to integrate **local LLM-powered reasoning** into modern full-stack systems.

---

## 🚀 Key Features

### 🎓 AI-Powered Academic Intelligence
- Upload transcripts (PDF or text) and automatically extract clean, structured text.  
- Generate **personalized course recommendations** using embedded reasoning and prompt engineering.  
- Summarize academic transcripts into concise, meaningful profiles.  
- Perform **AI-driven scholarship searches** by integrating Brave Search with model inference.  

### 💬 ChatGPT-Style AI Chat
- Fully local chat interface for **real-time conversation with the AI model**.  
- Streaming responses (token-by-token rendering).  
- Memoryless session for privacy — all chats are temporary.  

### 🧾 PDF Report Generation
- Auto-generates professional PDF summaries that include transcript analysis, recommendations, and scholarships.  
- Clickable links in the PDF open directly in the browser.  
- Files are stored locally under `/summaries` and can later be moved to cloud storage (e.g., AWS S3).  

### 🧠 Local LLM Inference
- Uses **Ollama** running locally with `gemma3:4b-it-qat`.  
- Works fully offline for privacy-preserving inference.  
- Easily replaceable with OpenAI API or custom hosted models.

---

## 🏗️ Architecture Overview

```plaintext
EduSphere/
├── server/            # Golang Fiber backend (AI logic, DB, Ollama integration)
├── client/            # React (Vite) frontend for chat and UI
└── README.md          # This overview file
```

**System Flow:**
1. User uploads transcript → Backend extracts and stores text.  
2. AI summarizes transcript → Generates strengths and recommendations.  
3. Brave search runs → Fetches scholarships and filters through LLM reasoning.  
4. Combined summary → PDF generated with courses + scholarships.  
5. Users can chat directly with the AI model in a ChatGPT-style interface.

---

## 🧰 Technologies Used

| Layer | Stack |
|-------|--------|
| **Frontend** | React (Vite), Axios, TailwindCSS, Lucide Icons |
| **Backend** | Golang Fiber, PostgreSQL (via sqlc), Ollama API |
| **AI/Inference** | Local inference with `gemma3:4b-it-qat` |
| **PDF Engine** | gofpdf |
| **Auth** | Paseto Token-based Authentication |

---

## ⚙️ Setup Instructions

### Prerequisites
- [Go 1.22+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [PostgreSQL 15+](https://www.postgresql.org/)
- [Ollama](https://ollama.ai) (for local model inference)

### Clone the Repository
```bash
git clone https://github.com/Nibir1/EduSphere.git
cd EduSphere
```

### Backend Setup
```bash
cd server
cp .env.example .env
go mod tidy
make migrateup
go run main.go
```

### Frontend Setup
```bash
cd client
npm install
npm run dev
```

Then visit **http://localhost:5173**

---

## 💡 Future Enhancements
- 🌐 AWS S3 support for cloud file persistence.  
- 🧬 Support for multiple model backends (OpenAI, Anthropic, Groq).  
- 📊 Analytics dashboard for performance tracking.  
- 🔐 Persistent chat memory using vector storage (e.g., pgvector).

---

## 🧭 Why This Project Matters
EduSphere showcases a **modern AI engineering stack** that combines **generative intelligence, production backend design, and privacy-first local inference**. It’s ideal as a **portfolio project** for demonstrating full-stack AI systems that are both technically deep and user-facing.
