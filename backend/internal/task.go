package api

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Task struct {
	ID        int    `json:"id"`
	ContestID int    `json:"contestid"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Status    string `json:"status"`
}

type TaskRequest struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	ContestID string `json:"contestid"`
}

func (a *Api) fetchTasksByContestID(ctx *fiber.Ctx) error {
	ContestID := ctx.Params("contestid")
	if ContestID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	ContestIDInt, err := strconv.Atoi(ContestID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный идентификатор контеста"})
	}

	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	userID := sess.Get("id")
	role := sess.Get("role")
	if userID == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	var groupID int
	err = a.db.QueryRow(`
        SELECT c.groupid
        FROM contests c
        JOIN users_in_groups uig ON c.groupid = uig.groupid
        WHERE c.id = $1 AND uig.userid = $2
    `, ContestIDInt, userID).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе, связанной с контестом"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	rows, err := a.db.Query(`
        SELECT id, ContestID, name, url
        FROM tasks
        WHERE ContestID = $1
    `, ContestIDInt)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.ContestID, &task.Name, &task.URL); err != nil {
			log.Println(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
		}

		if role == "student" {
			var submissionStatus string
			err := a.db.QueryRow(`
                SELECT status
                FROM submissions
                WHERE taskid = $1 AND userid = $2
                ORDER BY created_at DESC
                LIMIT 1
            `, task.ID, userID).Scan(&submissionStatus)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
			}
			task.Status = submissionStatus
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return ctx.JSON(tasks)
}

func (a *Api) deleteTask(ctx *fiber.Ctx) error {
	taskID := ctx.Params("taskid")
	if taskID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	userID := sess.Get("id")
	role := sess.Get("role")
	if userID == nil || role == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	if role != "teacher" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Недостаточно прав"})
	}

	var groupID int
	err = a.db.QueryRow(`
        SELECT c.groupid
        FROM tasks t
        JOIN contests c ON t.contestid = c.id
        JOIN users_in_groups uig ON c.groupid = uig.groupid
        WHERE t.id = $1 AND uig.userid = $2
    `, taskID, userID).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе, связанной с задачей"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	tx, err := a.db.BeginTx(ctx.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        DELETE FROM submissions
        WHERE taskid = $1
    `, taskID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	_, err = tx.Exec(`
        DELETE FROM tasks
        WHERE id = $1
    `, taskID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Task", taskID, "deleted")

	return ctx.JSON(fiber.Map{"success": true})
}

func (a *Api) addTask(ctx *fiber.Ctx) error {
	log.Println("addTask")
	var req TaskRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	userID := sess.Get("id")
	role := sess.Get("role")
	if userID == nil || role == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	if role != "teacher" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Недостаточно прав"})
	}

	var groupID int
	err = a.db.QueryRow(`
        SELECT c.groupid
        FROM contests c
        JOIN users_in_groups uig ON c.groupid = uig.groupid
        WHERE c.id = $1 AND uig.userid = $2
    `, req.ContestID, userID).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе, связанной с контестом"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	// Insert the new task into the database
	_, err = a.db.Exec(`
        INSERT INTO tasks (name, url, contestid)
        VALUES ($1, $2, $3)
    `, req.Name, req.URL, req.ContestID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Task added:", req.Name)

	return ctx.JSON(fiber.Map{"success": true})
}
