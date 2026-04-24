package repository

import (
	"context"

	"careerflow/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.DB.Exec(ctx, `INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)`, user.ID, user.Name, user.Email, user.PasswordHash)
	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow(ctx, `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func (r *Repository) CreateApplication(ctx context.Context, app models.Application) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO applications (id, user_id, company, role, status, job_url, job_description, resume_text, fit_score, strengths, gaps)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, app.ID, app.UserID, app.Company, app.Role, app.Status, app.JobURL, app.JobDescription, app.ResumeText, app.FitScore, app.Strengths, app.Gaps)
	return err
}

func (r *Repository) ListApplications(ctx context.Context, userID string) ([]models.Application, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, user_id, company, role, status, job_url, job_description, resume_text, fit_score, strengths, gaps, created_at, updated_at FROM applications WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []models.Application{}
	for rows.Next() {
		var app models.Application
		if err := rows.Scan(&app.ID, &app.UserID, &app.Company, &app.Role, &app.Status, &app.JobURL, &app.JobDescription, &app.ResumeText, &app.FitScore, &app.Strengths, &app.Gaps, &app.CreatedAt, &app.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, rows.Err()
}
