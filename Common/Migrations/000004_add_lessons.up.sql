CREATE TABLE IF NOT EXISTS lessons (
 lesson_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 course_id UUID NOT NULL REFERENCES courses(course_id),
 title TEXT NOT NULL,
 content TEXT NOT NULL
);
