-- EXTENSIONES
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ACADEMIC PERIODS
CREATE TABLE academic_periods (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL UNIQUE,             -- 2026-1, 2026-2
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  is_active BOOLEAN DEFAULT false
);

-- USERS
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- ROLES
CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL              -- ADMIN, PROFESOR, ESTUDIANTE
);

-- USER ROLES
CREATE TABLE user_roles (
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  role_id INT REFERENCES roles(id),
  PRIMARY KEY (user_id, role_id)
);

-- COURSES / PROJECTS
CREATE TABLE courses (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  type TEXT NOT NULL,                    -- COURSE | PROJECT
  professor_id UUID REFERENCES users(id),
  period_id UUID REFERENCES academic_periods(id),
  start_date DATE,
  end_date DATE,
  status TEXT NOT NULL,                  -- ACTIVE | CLOSED
  created_at TIMESTAMP DEFAULT NOW(),

  CHECK (type IN ('COURSE', 'PROJECT')),
  CHECK (status IN ('ACTIVE', 'CLOSED'))
);

-- ASSIGNMENTS
CREATE TABLE assignments (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id),
  course_id UUID REFERENCES courses(id),
  role_type TEXT NOT NULL,               -- MONITOR | ASSISTANT
  contracted_hours INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CHECK (role_type IN ('MONITOR', 'ASSISTANT')),
  CHECK (contracted_hours > 0),
  UNIQUE (user_id, course_id, role_type)
);

-- WEEKS
CREATE TABLE weeks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  period_id UUID REFERENCES academic_periods(id),  -- NUEVO: relación con período
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  is_active BOOLEAN DEFAULT false,
  is_closed BOOLEAN DEFAULT false,
  CHECK (start_date < end_date)
);

-- TASKS
CREATE TABLE tasks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
  week_id UUID REFERENCES weeks(id),
  title TEXT NOT NULL,
  description TEXT,
  hours NUMERIC(4,2) NOT NULL,
  status TEXT NOT NULL,                  -- PENDING | COMPLETED
  observations TEXT,
  attachment_url TEXT,
  is_late BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW(),
  CHECK (status IN ('PENDING', 'COMPLETED')),
  CHECK (hours > 0)
);

-- REPORTS
CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  assignment_id UUID REFERENCES assignments(id) ON DELETE CASCADE,
  week_id UUID REFERENCES weeks(id),
  total_hours NUMERIC(5,2) DEFAULT 0,
  file_url TEXT,
  generated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE (assignment_id, week_id)
);

-- NOTIFICATIONS
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW()
);

-- INDICES
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_assignments_user ON assignments(user_id);
CREATE INDEX idx_assignments_course ON assignments(course_id);
CREATE INDEX idx_tasks_assignment ON tasks(assignment_id);
CREATE INDEX idx_tasks_week ON tasks(week_id);
CREATE INDEX idx_reports_assignment ON reports(assignment_id);
CREATE INDEX idx_reports_week ON reports(week_id);
CREATE INDEX idx_weeks_period ON weeks(period_id);  -- NUEVO: índice para la nueva relación

-- DATA INICIAL
INSERT INTO roles (name) VALUES
('ADMIN'),
('PROFESOR'),
('ESTUDIANTE');

WITH new_user AS (
  INSERT INTO users (name, email, password_hash)
  VALUES ('Test User', 'test@test.com', '$2a$10$mMiZQtzpsQ8KKKc1DYeNP.F.txM6F9xcJvSEFS8NEtEmWXife.0JO')
  RETURNING id
)
INSERT INTO user_roles (user_id, role_id)
SELECT id, 1 FROM new_user;

INSERT INTO academic_periods (name, start_date, end_date, is_active)
VALUES ('2026-1', '2026-01-15', '2026-06-15', true);