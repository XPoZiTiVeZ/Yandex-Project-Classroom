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