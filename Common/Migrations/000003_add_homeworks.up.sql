CREATE TABLE IF NOT EXISTS homeworks (
  homework_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  course_id UUID NOT NULL REFERENCES courses(course_id),
  title TEXT NOT NULL,
  content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS homework_submissions (
  solution TEXT,
  submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
  student_id UUID NOT NULL REFERENCES users(user_id),
  homework_id UUID NOT NULL REFERENCES homeworks(homework_id),
  PRIMARY KEY (student_id, homework_id)
);