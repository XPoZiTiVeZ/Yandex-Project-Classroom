package repository

import (
	models "Classroom/Courses/internal/repo/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r *CourseRepo) CreateCourse(ctx context.Context, teacherID, title, description string, visibility bool, startTime, endTime *time.Time) (string, error) {
	var courseID string

	query := `INSERT INTO courses 
		(course_id, teacher_id, title, description, visibility, start_time, end_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING course_id`

	err := r.db.QueryRow(
		ctx, query,
		teacherID,
		title,
		description,
		visibility,
		startTime,
		endTime,
	).Scan(&courseID)

	if err != nil {
		return "", fmt.Errorf("CreateCourse error: %v", err)
	}

	return courseID, nil
}

func (r *CourseRepo) DeleteCourse(ctx context.Context, courseID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM courses WHERE course_id = $1`, courseID)

	if err != nil {
		return fmt.Errorf("DeleteCourse error: %v", err)
	}

	return nil
}

func (r *CourseRepo) GetCourse(ctx context.Context, courseID string) (*models.Course, error) {
	query := `
        SELECT course_id, teacher_id, title, description, visibility, start_time, end_time
        FROM courses
        WHERE course_id = $1
    `

	var course models.Course
	err := r.db.QueryRow(ctx, query, courseID).Scan(
		&course.CourseID,
		&course.TeacherID,
		&course.Title,
		&course.Description,
		&course.Visibility,
		&course.StartTime,
		&course.EndTime,
	)

	if err != nil {
		return nil, fmt.Errorf("GetCourse error: %v", err)
	}

	return &course, nil
}

func (r *CourseRepo) GetCourses(ctx context.Context, teacherID string) ([]*models.Course, error) {
	var courses []*models.Course

	query := `
	SELECT course_id, teacher_id, title, description, visibility, start_time, end_time
	FROM courses
	WHERE teacher_id = $1
	`

	rows, err := r.db.Query(ctx, query, teacherID)

	if err != nil {
		return nil, fmt.Errorf("GetCourses error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(
			&course.CourseID,
			&course.TeacherID,
			&course.Title,
			&course.Description,
			&course.Visibility,
			&course.StartTime,
			&course.EndTime,
		); err != nil {
			return nil, fmt.Errorf("GetCourses (scan row) error: %v", err)
		}
		courses = append(courses, &course)
	}

	return courses, nil
}

func (r *CourseRepo) UpdateCourse(ctx context.Context, courseID string, title, description *string, visibility *bool, startTime, endTime *time.Time) error {
	query := `
	UPDATE courses SET 
	title = COALESCE($2, title),
	description = COALESCE($3, description),
	visibility = COALESCE($4, visibility),
	start_time = COALESCE($5, start_time),
	end_time = COALESCE($6, end_time)
	WHERE course_id = $1
	`
	result, err := r.db.Exec(
		ctx,
		query,
		courseID,
		title,
		description,
		visibility,
		startTime,
		endTime,
	)

	if err != nil {
		return fmt.Errorf("UpdateCourse error: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("course not found")
	}

	return nil
}

func (r *CourseRepo) EnrollUser(ctx context.Context, courseID, userID string) error {
	query := `INSERT INTO enrollments (course_id, member_id, enrolled_at)
	          VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	_, err := r.db.Exec(
		ctx,
		query,
		courseID,
		userID,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("EnrollUser error: %v", err)
	}

	return nil
}

func (r *CourseRepo) ExpelUser(ctx context.Context, courseID, userID string) error {
	query := "DELETE FROM enrollments WHERE course_id = $1 AND member_id = $2"

	result, err := r.db.Exec(ctx, query, courseID, userID)

	if err != nil {
		return fmt.Errorf("ExpelUser error: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("enrollment not found")
	}

	return nil
}

func (r *CourseRepo) GetUserCourses(ctx context.Context, userID string) ([]*models.Course, error) {
	query := `SELECT 
    courses.course_id, 
    courses.teacher_id, 
    courses.title, 
    courses.description, 
    courses.visibility, 
    courses.start_time, 
    courses.end_time
FROM enrollments
JOIN courses ON courses.course_id = enrollments.course_id
WHERE enrollments.member_id = $1;`

	var courses []*models.Course

	rows, err := r.db.Query(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("GetUserCourses error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Course
		err := rows.Scan(&c.CourseID, &c.TeacherID, &c.Title, &c.Description, &c.Visibility, &c.StartTime, &c.EndTime)

		if err != nil {
			return nil, fmt.Errorf("GetUserCourses (row scan) error: %v", err)
		}

		courses = append(courses, &c)
	}

	return courses, nil
}

func (r *CourseRepo) IsTeacher(ctx context.Context, userID, courseID string) (bool, error) {
    const query = `SELECT teacher_id = $1 FROM courses WHERE course_id = $2`
    
    var isTeacher bool
    err := r.db.QueryRow(ctx, query, courseID).Scan(&isTeacher)
    
    if errors.Is(err, pgx.ErrNoRows) {
        return false, fmt.Errorf("course not found") // Курс не существует
    } else if err != nil {
        return false, fmt.Errorf("failed to verify teacher: %w", err) // Ошибка запроса
    }
    
    return isTeacher, nil
}

func (r *CourseRepo) IsMember(ctx context.Context, courseID, userID string) (bool, error) {
    query := `SELECT true FROM enrollments WHERE course_id = $1 AND student_id = $2`
	
	var exists bool
	err := r.db.QueryRow(ctx, query, courseID).Scan(&exists)
	
	if errors.Is(err, pgx.ErrNoRows) {
        return false, fmt.Errorf("course not found") // Курс не существует
    }	
	
    return true, nil
}

func (r *CourseRepo) GetCourseMembers(ctx context.Context, courseID string, index, limit int32) (int32, []models.Member, error) {
    var courseExists bool
    err := r.db.QueryRow(
        ctx,
        `SELECT EXISTS(SELECT 1 FROM courses WHERE course_id = $1)`,
        courseID,
    ).Scan(&courseExists)
    
    if err != nil {
        return 0, []models.Member{}, fmt.Errorf("failed to check course existence: %w", err)
    }
    if !courseExists {
        return 0, []models.Member{}, fmt.Errorf("course not found")
    }

	var total int32
    err = r.db.QueryRow(
        ctx,
        `SELECT COUNT(*) FROM enrollments WHERE course_id = $1`,
        courseID,
    ).Scan(&total)
    
    if err != nil {
        return 0, nil, fmt.Errorf("failed to count members: %w", err)
    }
    
    offset := index * limit
    rows, err := r.db.Query(
        ctx,
        `SELECT 
            u.user_id, 
            u.email, 
            u.first_name, 
            u.last_name
         FROM enrollments e
         JOIN users u ON e.student_id = u.user_id
         WHERE e.course_id = $1
         ORDER BY u.last_name, u.first_name
         LIMIT $2 OFFSET $3`,
        courseID,
        limit,
        offset,
    )
    
    if err != nil {
        return 0, nil, fmt.Errorf("failed to query members: %w", err)
    }
    defer rows.Close()

	var members []models.Member
    for rows.Next() {
        var m models.Member
        if err := rows.Scan(
            &m.UserID,
            &m.Email,
            &m.FirstName,
            &m.LastName,
        ); err != nil {
            return 0, nil, fmt.Errorf("failed to scan member: %w", err)
        }
        members = append(members, m)
    }

    if err := rows.Err(); err != nil {
        return 0, nil, fmt.Errorf("rows iteration error: %w", err)
    }

    return total, members, nil
}