package api

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	Password  string
	Firstname string
	Surname   string
	Role      string
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
	// return hashedPassword == password
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
	// return password, nil
}

func (a *Api) login(ctx *fiber.Ctx) error {
	username := ctx.Query("username")
	password := ctx.Query("password")

	if ok, err := checkCorrectWord(password, "Password"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if ok, err := checkCorrectWord(username, "Username"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	var user User
	err := a.db.QueryRow("SELECT id, password, firstname, surname, role FROM users WHERE username = $1", username).Scan(&user.ID, &user.Password, &user.Firstname, &user.Surname, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверные данные"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	if !verifyPassword(user.Password, password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверные данные"})
	}

	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	sess.Set("id", user.ID)
	sess.Set("role", user.Role)
	if err := sess.Save(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("User", user.ID, "logged in")

	return ctx.JSON(fiber.Map{
		"firstname": user.Firstname,
		"surname":   user.Surname,
		"role":      user.Role,
		"id":        user.ID,
	})
}

func (a *Api) register(ctx *fiber.Ctx) error {
	username := ctx.Query("username")
	password := ctx.Query("password")
	firstname := ctx.Query("firstname")
	surname := ctx.Query("surname")
	role := ctx.Query("role")

	if ok, err := checkCorrectWord(username, "Имя пользователя"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if ok, err := checkCorrectWord(password, "Пароль"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if ok, err := checkCorrectWord(firstname, "Имя"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if ok, err := checkCorrectWord(surname, "Фамилия"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if ok, err := checkCorrectWord(role, "Роль"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if role != "teacher" && role != "student" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверная роль"})
	}

	var existingUser User
	err := a.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&existingUser.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	if existingUser.ID != 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Имя пользователя уже занято"})
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	var userID int
	err = a.db.QueryRow(
		"INSERT INTO users (username, password, firstname, surname, role) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		username, hashedPassword, firstname, surname, role,
	).Scan(&userID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	sess.Set("id", userID)
	sess.Set("role", role)
	if err := sess.Save(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("User", username, "registered")

	return ctx.JSON(fiber.Map{
		"firstname": firstname,
		"surname":   surname,
		"role":      role,
		"id":        userID,
	})
}
