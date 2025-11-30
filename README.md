# 🧠 EduSphere — AI-Powered Academic Assistant

> **Ambitious Full-Stack Generative AI System** built with **Golang Fiber**, **PostgreSQL**, and **React (Vite)** — integrating local LLM inference, dynamic reasoning, and production-grade software design.

EduSphere isn’t just another AI demo — it’s a **complete Generative AI platform** that transforms student transcripts into **personalized insights, course recommendations, and scholarship matches** — all powered by **on-device LLM inference** through Ollama.

It represents a **real-world AI Systems Engineering project**, blending backend scalability, AI reasoning, and modern UI design into a cohesive and professional-grade product.

---

## 🚀 Why EduSphere Is an Ambitious Project

- 🧩 **End-to-End System Design** — Secure multi-user authentication, inference orchestration, and persistent data handling.  
- 🧠 **AI Reasoning Layer** — Summarization, academic profiling, and context-based course & scholarship discovery.  
- ⚙️ **Production-Grade Backend** — Golang Fiber + PostgreSQL + structured routes + concurrency-safe architecture.  
- 💬 **Streaming Chat Interface** — Real-time chat UX built in React, mirroring ChatGPT’s conversational flow.  
- 📄 **Dynamic PDF Reports** — Auto-generated academic summaries with integrated AI reasoning.  
- 🔒 **Local & Private Inference** — Runs LLMs directly on-device with **Ollama**, ensuring privacy and independence from cloud APIs.

---

## 🧩 System Architecture

| Layer | Stack | Description |
|-------|--------|-------------|
| **Frontend** | React (Vite), TailwindCSS, Lucide Icons | Real-time chat UI, scholarship discovery, and summary dashboards |
| **Backend** | Golang (Fiber), PostgreSQL | Token auth, modular routes, and production-grade data persistence |
| **AI Engine** | Ollama + Local LLMs (Gemma / Llama / Mistral) | Summarization, reasoning, and conversational inference |
| **Storage** | PostgreSQL + Filesystem | Structured persistence for recommendations and generated reports |
| **DevOps** | Docker, Makefile | Local development setup, easy build and run workflow |

---

## 🧠 AI Engineering Highlights

- Local inference using **Ollama** (no cloud dependency)  
- Custom **prompt orchestration** for multi-step academic reasoning  
- **RAG-ready architecture** for future integration with vector databases  
- **Hybrid reasoning** combining transcript data and Brave search results  
- ChatGPT-style **streaming LLM chat** with markdown rendering  
- Professional **PDF generation pipeline** with summaries, recommendations, and scholarships

---

## 🧰 System Workflow

```plaintext
User Uploads Transcript
        ↓
AI Summarizes Academic Profile
        ↓
Course Recommendations (LLM Reasoning)
        ↓
Scholarship Matching (Web Search + AI Filtering)
        ↓
PDF Report Generation
        ↓
Optional Chat with AI (Real-Time Streaming)
```
---
### ✨ List of Functionalities this project can do

- Analyze uploaded student transcripts and extract structured insights.  
- Generate AI-powered academic summaries highlighting key strengths and subjects.  
- Recommend personalized academic courses based on transcript content, inferred interests, and existing courses from the database.  
- Perform real-time web searches (via Brave Search API) for scholarships relevant to a student’s profile.  
- Use AI filtering and ranking to match the most suitable scholarships based on relevance and fit.  
- Integrate scholarship details (title, description, match score, and URL) directly into the user dashboard.  
- Dynamically include discovered scholarships into summary reports or PDF exports.  
- Handle duplicate filtering and sanitization to ensure unique and clean scholarship results.  
- Generate academic summaries using local LLM inference (via Ollama) — no external API needed.  
- Run natural language reasoning pipelines locally for complete data privacy.  
- Support prompt-based orchestration for summarization, recommendation, and contextual reasoning.  
- Provide a real-time chat interface for natural, conversational interaction with EduSphere AI.  
- Generate dynamic academic reports (PDF) including summary, recommendations, and scholarships.  
- Include clickable scholarship links within the PDF for user convenience.  
- Use structured report layouts with user information, creation date, and clean professional typography.  
- Save generated PDFs locally on the backend, linked securely to individual user accounts.  
- Allow users to download saved reports directly from the web interface.  
- Support deletion of old reports with full backend file cleanup.  
- Automatically prevent duplicate or incomplete report generation.  
- Use **Golang Fiber** for a fast, production-grade backend API.  
- Handle authentication and authorization via secure **JWT tokens**.  
- Maintain per-user data isolation — transcripts, summaries, and PDFs are always user-specific.  
- Store structured data using **PostgreSQL**, with relationships between users, recommendations, and summaries.  
- Manage concurrent API calls (e.g., AI inference and PDF generation) safely and efficiently.  
- Support long-running AI inference operations with extended HTTP timeouts.  
- Log detailed backend operations for full transparency and debugging.  
- Built with **React (Vite)** — fast, modular, and optimized for developer experience.  
- Fully responsive UI for both desktop and mobile devices.  
- Feature a dynamic dashboard displaying document statistics, course counts, and scholarship matches.  
- Provide progress indicators for long-running AI tasks (loading, generating, saving, etc.).  
- Handle smooth state management for simultaneous actions (e.g., generating summaries while fetching scholarships).  
- Implement secure token-based authentication with automatic session expiration handling.  
- Present a clean, professional, and accessible UI using **TailwindCSS** and **Lucide icons**.  

---

## ⚙️ Setup & Run

### Prerequisites

- **Golang** ≥ 1.22  
- **Node.js** ≥ 18  
- **PostgreSQL**  
- **Ollama** installed locally (`https://ollama.ai`)  

### Backend Setup

```bash
cd server
go run main.go
```

### Frontend Setup

```bash
cd client
npm install
npm run dev
```

---

## 🧩 Key Features Summary

| Feature | Description |
|----------|--------------|
| AI Transcript Summarization | LLM-based academic insight extraction |
| Course Recommendations | Personalized academic paths based on transcript content |
| Scholarship Discovery | Brave API + AI filtering for relevant global scholarships |
| Dynamic PDF Reports | Summaries, recommendations, and scholarships in one file |
| Real-Time Chat | ChatGPT-style chat with streaming responses |
| Privacy First | Fully local inference using Ollama (no data leaves your system) |

---

## 💡 Why It Matters

EduSphere demonstrates **end-to-end Generative AI Systems Engineering** —  
combining **AI reasoning, backend scalability, and human-centered interaction** into a seamless platform.

It’s designed to showcase the kind of **architecture and applied AI thinking** that modern companies expect from **AI Engineers and Full-Stack Developers** building production-grade GenAI tools.

---

## 🧠 Built With

- **Golang (Fiber Framework)** — backend & API design  
- **PostgreSQL** — structured relational data storage  
- **React + Vite + TailwindCSS** — frontend experience  
- **Ollama (Local LLM Inference)** — private, on-device AI reasoning  
- **Docker + Makefile** — streamlined dev & deployment environment  

---

## 🏆 Project Scope

EduSphere reflects:  
- Real-world **LLM orchestration** and **AI safety practices**  
- Production-grade **backend design principles**  
- Deep understanding of **human-AI interaction systems**  
- Full-stack integration of **AI, data, and UX**  

---

> **EduSphere** — A showcase of applied AI engineering, full-stack system design, and the power of local intelligence.
