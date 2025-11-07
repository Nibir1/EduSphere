package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Course struct {
	Code             string
	Name             string
	Language         string
	GradingScale     string
	Organiser        string
	LearningOutcomes string
	Prerequisites    string
	TeacherName      string
	TeacherEmail     string
	CourseLink       string
	CreatedAt        time.Time
}

func main() {
	connStr := "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	defer db.Close()

	courses := []Course{
		{
			Code:             "TIES4911",
			Name:             "Deep Learning for Cognitive Computing for Developers",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Students will learn how to build Machine Intelligence based solutions using TensorFlow and cloud-based AI services (IBM, Google, Microsoft).`,
			Prerequisites:    `Basic Python programming knowledge and familiarity with AI or data mining concepts.`,
			TeacherName:      "Oleksiy Khriyenko",
			TeacherEmail:     "oleksiy.khriyenko@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-d582f31e-a375-48c2-9932-3af18c9ae8c5/brochure",
		},
		{
			Code:             "TIES4480",
			Name:             "Software Engineering Project",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Students gain hands-on experience managing real-world software projects using agile methodologies and team collaboration tools.`,
			Prerequisites:    `Basic programming and software design experience.`,
			TeacherName:      "Mikko Lehtinen",
			TeacherEmail:     "mikko.lehtinen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-xxxxxx-project",
		},
		{
			Code:             "TIES4567",
			Name:             "Data Mining and Knowledge Discovery",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science and Information Systems",
			LearningOutcomes: `Learn how to apply data mining techniques to extract useful knowledge and patterns from large datasets.`,
			Prerequisites:    `Basic statistics, linear algebra, and programming skills.`,
			TeacherName:      "Laura Nieminen",
			TeacherEmail:     "laura.nieminen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-data-mining",
		},
		{
			Code:             "TIES4720",
			Name:             "Information Systems Strategy and Digital Transformation",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Systems Science",
			LearningOutcomes: `Understand how organizations leverage information systems to drive digital transformation and strategic advantage.`,
			Prerequisites:    `Introductory course in information systems or business management.`,
			TeacherName:      "Kaisa Saarinen",
			TeacherEmail:     "kaisa.saarinen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-strategy",
		},
		{
			Code:             "TIES4030",
			Name:             "Artificial Intelligence in Practice",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Develop AI-driven solutions using modern libraries and explore ethical considerations in AI deployment.`,
			Prerequisites:    `Familiarity with Python, machine learning frameworks, and data preprocessing.`,
			TeacherName:      "Janne Mäkelä",
			TeacherEmail:     "janne.makela@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-ai-practice",
		},
	}

	for _, c := range courses {
		var exists bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM courses WHERE code=$1)`, c.Code).Scan(&exists)
		if err != nil {
			log.Printf("Failed checking existence for %s: %v", c.Code, err)
			continue
		}
		if exists {
			log.Printf("Skipping existing course: %s", c.Code)
			continue
		}

		_, err = db.Exec(`
			INSERT INTO courses (
				code, name, language, grading_scale, organiser,
				learning_outcomes, prerequisites, teacher_name, teacher_email, course_link, created_at
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		`,
			c.Code, c.Name, c.Language, c.GradingScale, c.Organiser,
			c.LearningOutcomes, c.Prerequisites, c.TeacherName, c.TeacherEmail, c.CourseLink, time.Now(),
		)
		if err != nil {
			log.Printf("Error inserting course %s: %v", c.Code, err)
		} else {
			log.Printf("Inserted course: %s", c.Name)
		}
	}

	fmt.Println("✅ Course seeding completed successfully.")
}
