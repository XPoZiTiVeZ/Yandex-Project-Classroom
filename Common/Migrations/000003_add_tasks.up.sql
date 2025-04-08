CREATE TABLE IF NOT EXISTS tasks (
  task_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  course_id UUID NOT NULL REFERENCES courses(course_id),
  title TEXT NOT NULL,
  content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS task_submissions (
  solution TEXT,
  submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
  student_id UUID NOT NULL REFERENCES users(user_id),
  task_id UUID NOT NULL REFERENCES tasks(task_id),
  PRIMARY KEY (student_id, task_id)
);