package api

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	db "github.com/nibir1/go-fiber-postgres-REST-boilerplate/db/sqlc"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/token"
	"github.com/nibir1/go-fiber-postgres-REST-boilerplate/util"

	"github.com/gofiber/swagger"
	_ "github.com/nibir1/go-fiber-postgres-REST-boilerplate/docs"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	app        *fiber.App
	validate   *validator.Validate

	uploadsDir   string
	summariesDir string
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	app := fiber.New(fiber.Config{})

	app.Use(logger.New())

	allowedOrigins := config.AllowedOrigins
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length, Content-Type",
		AllowCredentials: true,
	}))

	validate := validator.New()
	validate.RegisterValidation("currency", validCurrency)

	server := &Server{
		config:       config,
		store:        store,
		tokenMaker:   tokenMaker,
		app:          app,
		validate:     validate,
		uploadsDir:   "./uploads",
		summariesDir: "./summaries",
	}

	// ensure folders
	_ = os.MkdirAll(server.uploadsDir, 0o755)
	_ = os.MkdirAll(server.summariesDir, 0o755)

	server.setUpRoutes()
	return server, nil
}

func (server *Server) setUpRoutes() {
	app := server.app

	// Swagger if you like
	app.Get("/swagger/*", swagger.HandlerDefault)

	// PUBLIC
	api := app.Group("/api")
	api.Post("/users", server.createUser)
	api.Post("/users/login", server.loginUser)

	// PROTECTED
	auth := api.Group("/", authMiddlewareFiber(server.tokenMaker))

	// ====== NEW EDU-SPHERE ROUTES ======
	// transcripts
	auth.Post("/transcripts/upload", server.uploadTranscript)
	auth.Get("/transcripts", server.listTranscripts)
	auth.Get("/transcripts/:id", server.getTranscript)

	// recommendations
	auth.Post("/recommendations", server.createRecommendation) // body: { transcript_id }
	auth.Post("/recommendations/generate", server.generateRecommendations)
	auth.Get("/recommendations", server.listRecommendations)
	auth.Get("/recommendations/:id", server.getRecommendation)

	// summaries
	auth.Post("/summaries", server.createSummaryPDF) // body: { recommendation_id }
	auth.Get("/summaries", server.listSummaries)
	auth.Get("/summaries/:id/download", server.downloadSummaryPDF)
	auth.Delete("/summaries/:id", server.deleteSummary)

	// chat (simple)
	auth.Post("/chat", server.chatOnce)
}

func (server *Server) Start(address string) error {
	return server.app.Listen(address)
}

func errorResponse(err error) fiber.Map {
	return fiber.Map{"error": err.Error()}
}
