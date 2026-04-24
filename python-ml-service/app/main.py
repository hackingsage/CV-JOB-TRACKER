from collections import Counter
from fastapi import FastAPI
from pydantic import BaseModel
import re

app = FastAPI(title="CareerFlow Analyzer", version="1.0.0")

class AnalyzeRequest(BaseModel):
    jobDescription: str
    resumeText: str

class AnalyzeResponse(BaseModel):
    fitScore: int
    strengths: list[str]
    gaps: list[str]
    summary: str


def tokenize(text: str) -> list[str]:
    return re.findall(r"[a-zA-Z][a-zA-Z0-9+#.-]+", text.lower())


def top_keywords(text: str, limit: int = 12) -> list[str]:
    stop = {
        "the", "and", "for", "with", "that", "this", "from", "have", "will", "your",
        "about", "into", "role", "team", "you", "our", "are", "but", "not", "use"
    }
    words = [w for w in tokenize(text) if w not in stop and len(w) > 2]
    return [word for word, _ in Counter(words).most_common(limit)]


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/analyze", response_model=AnalyzeResponse)
def analyze(req: AnalyzeRequest):
    jd_keywords = top_keywords(req.jobDescription, 15)
    resume_tokens = set(tokenize(req.resumeText))
    matched = [kw for kw in jd_keywords if kw in resume_tokens]
    missing = [kw for kw in jd_keywords if kw not in resume_tokens]

    score = 35
    if jd_keywords:
        score += int((len(matched) / len(jd_keywords)) * 65)
    score = max(0, min(score, 100))

    strengths = matched[:5] if matched else ["General profile alignment detected"]
    gaps = missing[:5] if missing else ["No major keyword gaps detected"]
    summary = (
        f"Resume matches {len(matched)} of the top {len(jd_keywords)} job keywords. "
        f"Focus on strengthening: {', '.join(gaps[:3])}."
    )

    return AnalyzeResponse(
        fitScore=score,
        strengths=strengths,
        gaps=gaps,
        summary=summary,
    )
