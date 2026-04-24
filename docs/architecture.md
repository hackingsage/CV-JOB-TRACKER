# Architecture

## Overview
CareerFlow AI is a portfolio-grade full stack application for managing job applications and measuring resume-to-role fit.

## Stack
- Frontend: React + Vite + Recharts
- API: Go + Gin + PostgreSQL + JWT auth
- Analysis service: Python + FastAPI
- Infra: Docker Compose

## Flow
1. User registers or logs in through the Go API.
2. User submits a job description and resume text.
3. Go backend calls the Python analysis service.
4. Python returns a fit score, strengths, gaps, and summary.
5. Go persists the application and returns the analysis result.
6. Frontend visualizes application health and fit trends.

## Why this is good for a CV
- Multi-service architecture
- Different languages in one production-style system
- Authentication and persistence
- Data visualization
- Dockerized developer experience
- Clean story for recruiters and interviewers
