CREATE TABLE IF NOT EXISTS users (
 user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 email TEXT NOT NULL UNIQUE,
 password_hash BYTEA NOT NULL,
 is_superuser BOOLEAN DEFAULT FALSE,
 first_name TEXT NOT NULL,
 last_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS courses (
 course_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 teacher_id UUID NOT NULL REFERENCES users(user_id),
 title TEXT NOT NULL,
 description TEXT NOT NULL,
 visibility BOOLEAN NOT NULL DEFAULT FALSE,
 start_time TIMESTAMP,
 end_time TIMESTAMP,
 created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS enrollments (
 course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
 student_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
 enrolled_at TIMESTAMP NOT NULL DEFAULT NOW(),
 PRIMARY KEY (course_id, student_id)
);

CREATE TABLE IF NOT EXISTS tasks (
  task_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  course_id UUID NOT NULL REFERENCES courses(course_id),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS task_submissions (
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  student_id UUID NOT NULL REFERENCES users(user_id),
  task_id UUID NOT NULL REFERENCES tasks(task_id),
  PRIMARY KEY (student_id, task_id)
);

CREATE TABLE IF NOT EXISTS lessons (
 lesson_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 course_id UUID NOT NULL REFERENCES courses(course_id),
 title TEXT NOT NULL,
 content TEXT NOT NULL,
 created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
