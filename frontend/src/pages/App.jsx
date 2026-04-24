import { useMemo, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, RadarChart, PolarGrid, PolarAngleAxis, Radar } from 'recharts'
import { cvBullets, featureCards, roadmap } from '../data/portfolioContent'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
const demoCredentials = {
  email: 'demo@careerflow.dev',
  password: 'password123',
}

const fallbackApps = [
  { company: 'Stripe', role: 'Backend Engineer', fitScore: 88, status: 'Interview' },
  { company: 'Notion', role: 'Platform Engineer', fitScore: 79, status: 'Applied' },
  { company: 'Airbnb', role: 'Full Stack Engineer', fitScore: 72, status: 'Draft' },
]

const radarData = [
  { skill: 'Backend', value: 92 },
  { skill: 'Data', value: 74 },
  { skill: 'Frontend', value: 84 },
  { skill: 'Infra', value: 78 },
  { skill: 'Product', value: 81 },
]

export default function App() {
  const [form, setForm] = useState({
    company: '',
    role: '',
    status: 'Applied',
    jobDescription: '',
    resumeText: '',
  })
  const [result, setResult] = useState(null)
  const [applications, setApplications] = useState(fallbackApps)
  const [loading, setLoading] = useState(false)
  const [bootstrapping, setBootstrapping] = useState(false)
  const [error, setError] = useState('')
  const [demoMode, setDemoMode] = useState('local preview')

  const chartData = useMemo(() => applications.map(app => ({ name: app.company, score: app.fitScore })), [applications])

  const onChange = event => {
    setForm(prev => ({ ...prev, [event.target.name]: event.target.value }))
  }

  const localAnalysis = () => {
    const keywords = ['go', 'python', 'api', 'postgres', 'docker', 'react', 'microservices', 'kafka', 'cloud']
    const resume = form.resumeText.toLowerCase()
    const jobDescription = form.jobDescription.toLowerCase()
    const ranked = keywords.filter(k => jobDescription.includes(k))
    const matches = ranked.filter(k => resume.includes(k))
    const gaps = ranked.filter(k => !resume.includes(k))
    return {
      fitScore: Math.min(100, 42 + matches.length * 9),
      strengths: matches.slice(0, 5),
      gaps: gaps.slice(0, 5),
      summary: `You align strongly on ${matches.slice(0, 3).join(', ') || 'general engineering foundations'}. To improve this application, strengthen ${gaps.slice(0, 2).join(', ') || 'measurable impact statements'}.`
    }
  }

  const loginDemo = async () => {
    const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(demoCredentials)
    })
    if (!response.ok) {
      throw new Error('Demo login failed')
    }
    return response.json()
  }

  const loadApplications = async token => {
    const response = await fetch(`${API_BASE_URL}/api/applications`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!response.ok) {
      throw new Error('Could not load applications')
    }
    const data = await response.json()
    setApplications(data.applications)
  }

  const bootstrapDemo = async () => {
    setBootstrapping(true)
    setError('')
    try {
      const auth = await loginDemo()
      await loadApplications(auth.token)
      setDemoMode('live API mode')
    } catch (err) {
      setDemoMode('local preview')
      setError('Live services are not reachable yet, so the UI is showing polished local demo data.')
    } finally {
      setBootstrapping(false)
    }
  }

  const onSubmit = async event => {
    event.preventDefault()
    setLoading(true)
    setError('')
    try {
      const auth = await loginDemo()
      const response = await fetch(`${API_BASE_URL}/api/applications/analyze`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${auth.token}`,
        },
        body: JSON.stringify(form),
      })
      if (!response.ok) {
        throw new Error('Live analysis failed')
      }
      const data = await response.json()
      setResult({ ...data.application, summary: data.summary })
      await loadApplications(auth.token)
      setDemoMode('live API mode')
    } catch (err) {
      const fallback = localAnalysis()
      setResult(fallback)
      setDemoMode('local preview')
      setError('Live API was unavailable, so I fell back to local client-side scoring.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="shell">
      <section className="hero">
        <div>
          <p className="eyebrow">Portfolio-grade full stack system</p>
          <h1>CareerFlow AI</h1>
          <p className="subcopy">A serious showcase project for software roles. It turns the job search itself into a data product with scoring, persistence, analytics, and service boundaries.</p>
          <div className="hero-actions">
            <a href="#analyzer" className="primary-link">Try the analyzer</a>
            <a href="#cv" className="secondary-link">CV bullets</a>
            <button type="button" onClick={bootstrapDemo} disabled={bootstrapping}>{bootstrapping ? 'Connecting...' : 'Load live demo'}</button>
          </div>
          <p className="mode-pill">Mode: {demoMode}</p>
          {error && <p className="notice">{error}</p>}
        </div>
        <div className="stats">
          <div><strong>Go API</strong><span>Auth, persistence, routing</span></div>
          <div><strong>Python ML-style service</strong><span>Keyword fit scoring engine</span></div>
          <div><strong>React dashboard</strong><span>Charts, UX, storytelling</span></div>
          <div><strong>Docker + CI</strong><span>Portable dev workflow</span></div>
        </div>
      </section>

      <section className="feature-grid">
        {featureCards.map(card => (
          <div className="panel" key={card.title}>
            <h2>{card.title}</h2>
            <p>{card.text}</p>
          </div>
        ))}
      </section>

      <section className="grid">
        <div className="panel">
          <h2>Application fit dashboard</h2>
          <ResponsiveContainer width="100%" height={240}>
            <BarChart data={chartData}>
              <XAxis dataKey="name" />
              <YAxis domain={[0, 100]} />
              <Tooltip />
              <Bar dataKey="score" fill="#7c3aed" radius={[8, 8, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
          <div className="app-list">
            {applications.map(app => (
              <div className="app-row" key={`${app.company}-${app.role}`}>
                <div>
                  <strong>{app.company}</strong>
                  <span>{app.role}</span>
                </div>
                <div>
                  <strong>{app.fitScore}/100</strong>
                  <span>{app.status}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="panel">
          <h2>Profile radar</h2>
          <ResponsiveContainer width="100%" height={240}>
            <RadarChart data={radarData}>
              <PolarGrid />
              <PolarAngleAxis dataKey="skill" />
              <Radar dataKey="value" stroke="#22c55e" fill="#22c55e" fillOpacity={0.45} />
            </RadarChart>
          </ResponsiveContainer>
          <p className="muted">A recruiter can instantly read this as a product-minded engineering dashboard.</p>
        </div>
      </section>

      <section className="grid analyzer-grid" id="analyzer">
        <div className="panel">
          <h2>Analyze a role</h2>
          <form onSubmit={onSubmit} className="form">
            <input name="company" placeholder="Company" value={form.company} onChange={onChange} required />
            <input name="role" placeholder="Role" value={form.role} onChange={onChange} required />
            <select name="status" value={form.status} onChange={onChange}>
              <option>Applied</option>
              <option>Interview</option>
              <option>Offer</option>
              <option>Rejected</option>
            </select>
            <textarea name="jobDescription" placeholder="Paste job description" value={form.jobDescription} onChange={onChange} rows="7" required />
            <textarea name="resumeText" placeholder="Paste your resume text" value={form.resumeText} onChange={onChange} rows="7" required />
            <button type="submit" disabled={loading}>{loading ? 'Analyzing...' : 'Run analysis'}</button>
          </form>
        </div>

        <div className="panel result-panel">
          <h2>Why this sells on a CV</h2>
          <ul className="bullet-list" id="cv">
            {cvBullets.map(item => <li key={item}>{item}</li>)}
          </ul>
          <h3>Roadmap</h3>
          <ul className="bullet-list">
            {roadmap.map(item => <li key={item}>{item}</li>)}
          </ul>
        </div>
      </section>

      {result && (
        <section className="panel result">
          <h2>Analysis result</h2>
          <div className="score">{result.fitScore}<span>/100</span></div>
          <p>{result.summary}</p>
          <div className="result-grid">
            <div>
              <h3>Strengths</h3>
              <ul>{result.strengths.map(item => <li key={item}>{item}</li>)}</ul>
            </div>
            <div>
              <h3>Gaps</h3>
              <ul>{result.gaps.map(item => <li key={item}>{item}</li>)}</ul>
            </div>
          </div>
        </section>
      )}
    </div>
  )
}
