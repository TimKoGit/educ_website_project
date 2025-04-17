package api

import (
	"database/sql"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/lib/pq"
)

type Api struct {
	app     *fiber.App
	db      *sql.DB
	session *session.Store
}

func NewApi(connStr string) *Api {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024,
	})
	sn := session.New()
	return &Api{app: app, db: db, session: sn}
}

func (api *Api) UseCors(frontendURL string) {
	api.app.Use(logger.New())
	api.app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))
}

func (api *Api) Start(port string) {
	api.app.Listen(port)
}

func (api *Api) Close() {
	api.db.Close()
}

func checkCorrectWord(word, wordName string) (bool, string) {
	if word == "" {
		return false, wordName + " не должен быть пустым"
	}

	var validPattern = regexp.MustCompile(`^[а-яА-Яa-zA-Z0-9!@#\$%\^\&*\)\(+=._ -]+$`)

	if validPattern.MatchString(word) {
		return true, ""
	}

	return false, wordName + " должен состоять из цифр, букв и спецсимволов"
}
