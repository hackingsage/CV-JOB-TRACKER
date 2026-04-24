package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pgxpool.New(ctx, databaseURL)
}

func RunMigrations(pool *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS applications (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			company TEXT NOT NULL,
			role TEXT NOT NULL,
			status TEXT NOT NULL,
			job_url TEXT,
			job_description TEXT NOT NULL,
			resume_text TEXT NOT NULL,
			fit_score INT NOT NULL DEFAULT 0,
			strengths TEXT[] NOT NULL DEFAULT '{}',
			gaps TEXT[] NOT NULL DEFAULT '{}',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(context.Background(), q); err != nil {
			return err
		}
	}
	return nil
}

func SeedDemoData(pool *pgxpool.Pool) error {
	queries := []string{
		`INSERT INTO users (id, name, email, password_hash)
		 VALUES ('11111111-1111-1111-1111-111111111111', 'Demo User', 'demo@careerflow.dev', '$2b$12$FqqYjP1ClXHHAKPShEdaR.NndstkEuj/MQRtBQcYwtzOQNJoU6gI.')
		 ON CONFLICT (email) DO NOTHING;`,
		`INSERT INTO applications (id, user_id, company, role, status, job_url, job_description, resume_text, fit_score, strengths, gaps)
		 VALUES
		 ('22222222-2222-2222-2222-222222222221', '11111111-1111-1111-1111-111111111111', 'Stripe', 'Backend Engineer', 'Interview', 'https://stripe.com/jobs', 'Go microservices APIs distributed systems postgres kafka docker', 'Built Go APIs, Python services, Dockerized apps, PostgreSQL-backed platforms', 88, ARRAY['go','postgres','docker','apis'], ARRAY['kafka']),
		 ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'Notion', 'Platform Engineer', 'Applied', 'https://www.notion.so/careers', 'Platform engineering reliability observability python backend cloud', 'Built backend systems in Python and Go with observability focus', 79, ARRAY['python','backend','platform'], ARRAY['cloud','reliability']),
		 ('22222222-2222-2222-2222-222222222223', '11111111-1111-1111-1111-111111111111', 'Airbnb', 'Full Stack Engineer', 'Draft', 'https://careers.airbnb.com', 'React typescript backend product engineering experimentation', 'React dashboards, APIs, analytics tools, experimentation mindset', 72, ARRAY['react','product','analytics'], ARRAY['typescript','backend depth'])
		 ON CONFLICT (id) DO NOTHING;`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(context.Background(), q); err != nil {
			return err
		}
	}
	return nil
}
