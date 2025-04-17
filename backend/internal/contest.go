package api

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Contest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Time     string `json:"time"`
	Duration string `json:"duration"`
	GroupID  int    `json:"groupid"`
}

type ContestRequest struct {
	ID       int
	Name     string
	Time     string
	Duration string
	GroupID  string
}

func (a *Api) fetchContestsByGroupId(ctx *fiber.Ctx) error {
	log.Println("fetchContestsByGroupId")
	groupID := ctx.Query("groupid")
	if groupID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный идентификатор группы"})
	}

	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	userID := sess.Get("id")
	if userID == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	var userInGroupID int
	err = a.db.QueryRow(`
        SELECT id
        FROM users_in_groups
        WHERE userid = $1 AND groupid = $2
    `, userID, groupIDInt).Scan(&userInGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	rows, err := a.db.Query(`
        SELECT id, name, time, duration, groupid
        FROM contests
        WHERE groupid = $1
    `, groupIDInt)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer rows.Close()

	var contests []map[string]interface{}
	for rows.Next() {
		var c Contest
		if err := rows.Scan(&c.ID, &c.Name, &c.Time, &c.Duration, &c.GroupID); err != nil {
			log.Println(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
		}
		contestMap := map[string]interface{}{
			"id":       c.ID,
			"name":     c.Name,
			"time":     c.Time,
			"duration": c.Duration,
			"groupid":  c.GroupID,
		}
		contests = append(contests, contestMap)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return ctx.JSON(contests)
}

func (a *Api) addNewContest(ctx *fiber.Ctx) error {
	log.Println("addNewContest")
	var req ContestRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	groupIDInt, err := strconv.Atoi(req.GroupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный идентификатор группы"})
	}

	contest := Contest{Name: req.Name, Time: req.Time, Duration: req.Duration, GroupID: groupIDInt}

	if contest.Name == "" || contest.Time == "" || contest.Duration == "" || contest.GroupID == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Все поля обязательны для заполнения"})
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

	var userInGroupID int
	err = a.db.QueryRow(`
        SELECT id
        FROM users_in_groups
        WHERE userid = $1 AND groupid = $2
    `, userID, contest.GroupID).Scan(&userInGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	err = a.db.QueryRow(`
        INSERT INTO contests (name, time, duration, groupid)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, contest.Name, contest.Time, contest.Duration, contest.GroupID).Scan(&contest.ID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Contest", contest.ID, "added")

	return ctx.JSON(fiber.Map{"id": contest.ID, "name": contest.Name, "time": contest.Time, "duration": contest.Duration, "groupid": contest.GroupID})
}

func (a *Api) fetchContestById(ctx *fiber.Ctx) error {
	contestID := ctx.Params("contestid")
	if contestID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	// Convert contestID to int
	contestIDInt, err := strconv.Atoi(contestID)
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
	if userID == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	var contest Contest
	err = a.db.QueryRow(`
        SELECT id, name, time, duration, groupid
        FROM contests
        WHERE id = $1
    `, contestIDInt).Scan(&contest.ID, &contest.Name, &contest.Time, &contest.Duration, &contest.GroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Контест не найден"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	var userInGroupID int
	err = a.db.QueryRow(`
        SELECT id
        FROM users_in_groups
        WHERE userid = $1 AND groupid = $2
    `, userID, contest.GroupID).Scan(&userInGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе, связанной с контестом"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return ctx.JSON(contest)
}

func (a *Api) deleteContest(ctx *fiber.Ctx) error {
	contestID := ctx.Params("contestid")
	if contestID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	contestIDInt, err := strconv.Atoi(contestID)
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
	if userID == nil || role == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	if role != "teacher" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Недостаточно прав"})
	}

	var contest Contest
	err = a.db.QueryRow(`
        SELECT id, name, time, duration, groupid
        FROM contests
        WHERE id = $1
    `, contestIDInt).Scan(&contest.ID, &contest.Name, &contest.Time, &contest.Duration, &contest.GroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Контест не найден"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	var userInGroupID int
	err = a.db.QueryRow(`
        SELECT id
        FROM users_in_groups
        WHERE userid = $1 AND groupid = $2
    `, userID, contest.GroupID).Scan(&userInGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе, связанной с контестом"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	_, err = a.db.Exec(`
        DELETE FROM contests
        WHERE id = $1
    `, contestIDInt)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Contest", contestIDInt, "deleted")

	return ctx.JSON(fiber.Map{"success": true})
}
