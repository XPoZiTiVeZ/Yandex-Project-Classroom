CREATE TABLE IF NOT EXISTS lectures (
 lecture_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 course_id UUID NOT NULL REFERENCES courses(course_id),
 title TEXT NOT NULL,
 content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS lecture_materials (
 material_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
 object_url TEXT NOT NULL,
 lecture_id UUID NOT NULL REFERENCES lectures(lecture_id)
);