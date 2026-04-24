INSERT INTO users (id, name, email, password_hash)
VALUES ('11111111-1111-1111-1111-111111111111', 'Demo User', 'demo@careerflow.dev', '$2a$10$7EqJtq98hPqEX7fNZaFWoOHiqP8Zl9wQ6jvB8zV7h1Y7v8zWwQm7K')
ON CONFLICT (email) DO NOTHING;

INSERT INTO applications (id, user_id, company, role, status, job_url, job_description, resume_text, fit_score, strengths, gaps)
VALUES
('22222222-2222-2222-2222-222222222221', '11111111-1111-1111-1111-111111111111', 'Stripe', 'Backend Engineer', 'Interview', 'https://stripe.com/jobs', 'Go microservices APIs distributed systems postgres kafka docker', 'Built Go APIs, Python services, Dockerized apps, PostgreSQL-backed platforms', 88, ARRAY['go','postgres','docker','apis'], ARRAY['kafka']),
('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'Notion', 'Platform Engineer', 'Applied', 'https://www.notion.so/careers', 'Platform engineering reliability observability python backend cloud', 'Built backend systems in Python and Go with observability focus', 79, ARRAY['python','backend','platform'], ARRAY['cloud','reliability']),
('22222222-2222-2222-2222-222222222223', '11111111-1111-1111-1111-111111111111', 'Airbnb', 'Full Stack Engineer', 'Draft', 'https://careers.airbnb.com', 'React typescript backend product engineering experimentation', 'React dashboards, APIs, analytics tools, experimentation mindset', 72, ARRAY['react','product','analytics'], ARRAY['typescript','backend depth'])
ON CONFLICT (id) DO NOTHING;
