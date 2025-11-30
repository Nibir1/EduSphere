// db/seed/seed_courses.go
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
		{
			Code:             "TIES4810",
			Name:             "Cloud Computing and Virtualization",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science",
			LearningOutcomes: `Gain understanding of cloud infrastructure, virtualization technologies, and container orchestration platforms such as Docker and Kubernetes.`,
			Prerequisites:    `Basic knowledge of networking and operating systems.`,
			TeacherName:      "Pekka Rantanen",
			TeacherEmail:     "pekka.rantanen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-cloud-computing",
		},
		{
			Code:             "TIES4625",
			Name:             "Big Data Analytics",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Students learn scalable data processing, distributed computation frameworks, and analytical methods for massive datasets.`,
			Prerequisites:    `Knowledge of Python or Java and fundamentals of databases.`,
			TeacherName:      "Antti Virtanen",
			TeacherEmail:     "antti.virtanen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-big-data",
		},
		{
			Code:             "TIES4150",
			Name:             "Cybersecurity Fundamentals",
			Language:         "English",
			GradingScale:     "Pass/Fail",
			Organiser:        "Department of Computer Science and Information Systems",
			LearningOutcomes: `Understand core principles of cybersecurity including network security, cryptography, and secure software design.`,
			Prerequisites:    `Basic networking and programming knowledge.`,
			TeacherName:      "Matti Korhonen",
			TeacherEmail:     "matti.korhonen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-cybersecurity",
		},
		{
			Code:             "TIES4025",
			Name:             "Machine Learning Applications",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Explore various supervised and unsupervised machine learning algorithms and apply them to real-world datasets.`,
			Prerequisites:    `Introductory statistics and linear algebra.`,
			TeacherName:      "Elina Räsänen",
			TeacherEmail:     "elina.rasanen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-machine-learning",
		},
		{
			Code:             "TIES4335",
			Name:             "Natural Language Processing",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computing Sciences",
			LearningOutcomes: `Students learn techniques for text preprocessing, sentiment analysis, and language modeling using deep learning architectures.`,
			Prerequisites:    `Familiarity with Python, linear algebra, and probability theory.`,
			TeacherName:      "Ville Ojala",
			TeacherEmail:     "ville.ojala@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-nlp",
		},
		{
			Code:             "TIES4200",
			Name:             "Internet of Things (IoT) Systems",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science and Electrical Engineering",
			LearningOutcomes: `Learn to design, develop, and integrate IoT systems including sensor networks, embedded devices, and cloud integration.`,
			Prerequisites:    `Basic electronics and programming experience.`,
			TeacherName:      "Jari Hämäläinen",
			TeacherEmail:     "jari.hamalainen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-iot-systems",
		},
		{
			Code:             "TIES4850",
			Name:             "Ethics and Society in Artificial Intelligence",
			Language:         "English",
			GradingScale:     "Pass/Fail",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Examine the societal, ethical, and philosophical implications of AI technologies and algorithmic decision-making.`,
			Prerequisites:    `Introductory course in AI or computer ethics.`,
			TeacherName:      "Riikka Salonen",
			TeacherEmail:     "riikka.salonen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-ai-ethics",
		},
		{
			Code:             "TIES4990",
			Name:             "Advanced Topics in Neural Networks",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Study advanced deep learning architectures including transformers, GANs, and reinforcement learning frameworks.`,
			Prerequisites:    `Completion of a basic deep learning or machine learning course.`,
			TeacherName:      "Aleksi Leppänen",
			TeacherEmail:     "aleksi.leppanen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-advanced-nn",
		},
		{
			Code:             "TIES4510",
			Name:             "Software Architecture and Design Patterns",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Learn the principles of scalable software architecture, microservices, and common design patterns for robust systems.`,
			Prerequisites:    `Basic knowledge of object-oriented programming and software engineering.`,
			TeacherName:      "Petteri Hiltunen",
			TeacherEmail:     "petteri.hiltunen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-software-architecture",
		},
		{
			Code:             "TIES4900",
			Name:             "Blockchain and Distributed Ledger Technologies",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science and Information Systems",
			LearningOutcomes: `Understand the fundamentals of blockchain, smart contracts, and decentralized applications.`,
			Prerequisites:    `Knowledge of programming and computer networks.`,
			TeacherName:      "Henri Mikkonen",
			TeacherEmail:     "henri.mikkonen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-blockchain",
		},
		{
			Code:             "TIES4820",
			Name:             "Data Visualization and Storytelling",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computing Sciences",
			LearningOutcomes: `Students will learn how to visualize complex datasets and effectively communicate insights using modern visualization tools.`,
			Prerequisites:    `Basic knowledge of statistics and data analysis.`,
			TeacherName:      "Noora Väisänen",
			TeacherEmail:     "noora.vaisanen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-visualization",
		},
		{
			Code:             "TIES4630",
			Name:             "Reinforcement Learning Fundamentals",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Learn key reinforcement learning concepts including Q-learning, policy gradients, and environment design.`,
			Prerequisites:    `Basic knowledge of machine learning and programming.`,
			TeacherName:      "Eero Hakkarainen",
			TeacherEmail:     "eero.hakkarainen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-reinforcement-learning",
		},
		{
			Code:             "TIES4105",
			Name:             "Computer Networks and Distributed Systems",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science",
			LearningOutcomes: `Understand modern computer network protocols, distributed systems design, and cloud communication.`,
			Prerequisites:    `Introductory computer systems course.`,
			TeacherName:      "Sanna Pulkkinen",
			TeacherEmail:     "sanna.pulkkinen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-networks",
		},
		{
			Code:             "TIES4775",
			Name:             "Human–Computer Interaction",
			Language:         "English",
			GradingScale:     "Pass/Fail",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Explore interaction design, usability evaluation, and user experience principles in modern interfaces.`,
			Prerequisites:    `Basic programming and web design knowledge.`,
			TeacherName:      "Marja Koivisto",
			TeacherEmail:     "marja.koivisto@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-hci",
		},
		{
			Code:             "TIES4605",
			Name:             "DevOps and Continuous Integration",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science",
			LearningOutcomes: `Master tools and workflows for CI/CD pipelines, version control, and automated testing in software projects.`,
			Prerequisites:    `Experience with software development and version control.`,
			TeacherName:      "Markus Lappalainen",
			TeacherEmail:     "markus.lappalainen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-devops",
		},
		{
			Code:             "TIES4240",
			Name:             "Computer Vision and Image Analysis",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Learn fundamental and advanced methods in image processing, object detection, and visual recognition.`,
			Prerequisites:    `Basic linear algebra and machine learning.`,
			TeacherName:      "Sami Aalto",
			TeacherEmail:     "sami.aalto@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-computer-vision",
		},
		{
			Code:             "TIES4955",
			Name:             "AI for Robotics and Automation",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Automation and Robotics",
			LearningOutcomes: `Apply AI algorithms to robotic control, sensor fusion, and intelligent automation systems.`,
			Prerequisites:    `Basic robotics or machine learning knowledge.`,
			TeacherName:      "Tuomas Viitanen",
			TeacherEmail:     "tuomas.viitanen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-ai-robotics",
		},
		{
			Code:             "TIES4005",
			Name:             "Programming Paradigms and Languages",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science",
			LearningOutcomes: `Study imperative, functional, and concurrent programming paradigms and language design principles.`,
			Prerequisites:    `Basic programming knowledge.`,
			TeacherName:      "Teemu Hiltula",
			TeacherEmail:     "teemu.hiltula@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-paradigms",
		},
		{
			Code:             "TIES4925",
			Name:             "Edge Computing and Smart Systems",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Understand architecture and deployment of intelligent edge systems and real-time data analytics.`,
			Prerequisites:    `Knowledge of networking and distributed systems.`,
			TeacherName:      "Johanna Lindholm",
			TeacherEmail:     "johanna.lindholm@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-edge-computing",
		},
		{
			Code:             "TIES4860",
			Name:             "AI for Healthcare Applications",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Learn how to apply AI techniques in healthcare contexts including diagnostic imaging and predictive modeling.`,
			Prerequisites:    `Basic machine learning and data analysis.`,
			TeacherName:      "Heidi Korpela",
			TeacherEmail:     "heidi.korpela@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-ai-healthcare",
		},
		{
			Code:             "TIES4970",
			Name:             "Quantum Computing Concepts",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Computer Science",
			LearningOutcomes: `Introduction to quantum computation principles, qubits, and quantum algorithms such as Grover and Shor.`,
			Prerequisites:    `Knowledge of linear algebra and algorithms.`,
			TeacherName:      "Oskari Niemi",
			TeacherEmail:     "oskari.niemi@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-quantum-computing",
		},
		{
			Code:             "TIES4985",
			Name:             "AI Product Design and Innovation",
			Language:         "English",
			GradingScale:     "Pass/Fail",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Learn how to design, prototype, and evaluate AI-powered products that align with human needs and business goals.`,
			Prerequisites:    `Basic AI and product design knowledge.`,
			TeacherName:      "Nina Seppälä",
			TeacherEmail:     "nina.seppala@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-ai-product-design",
		},
		{
			Code:             "TIES4790",
			Name:             "Data Security and Privacy",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Information Systems",
			LearningOutcomes: `Learn how to protect sensitive data, manage access control, and apply privacy-preserving machine learning methods.`,
			Prerequisites:    `Basic understanding of cybersecurity and databases.`,
			TeacherName:      "Arto Välimäki",
			TeacherEmail:     "arto.valimaki@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-data-security",
		},
		{
			Code:             "TIES4960",
			Name:             "AI and Human Creativity",
			Language:         "English",
			GradingScale:     "Pass/Fail",
			Organiser:        "Faculty of Information Technology",
			LearningOutcomes: `Explore how AI systems collaborate with humans in creative domains such as art, music, and design.`,
			Prerequisites:    `Basic understanding of AI tools.`,
			TeacherName:      "Timo Salo",
			TeacherEmail:     "timo.salo@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-ai-creativity",
		},
		{
			Code:             "TIES4875",
			Name:             "Autonomous Systems and Self-Driving Vehicles",
			Language:         "English",
			GradingScale:     "General scale, 0-5",
			Organiser:        "Department of Robotics and Control",
			LearningOutcomes: `Understand key algorithms for perception, motion planning, and decision-making in autonomous systems.`,
			Prerequisites:    `Machine learning or robotics basics.`,
			TeacherName:      "Ilkka Paananen",
			TeacherEmail:     "ilkka.paananen@jyu.fi",
			CourseLink:       "https://sisu.jyu.fi/student/courseunit/otm-autonomous-systems",
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
