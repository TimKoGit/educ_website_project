package api

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Group struct {
	ID   int
	Name string
	Code string
}

func (a *Api) fetchGroupsByUserId(ctx *fiber.Ctx) error {
	log.Println("fetchGroupsByUserId")
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

	rows, err := a.db.Query(`
            SELECT g.name, g.id
            FROM groups g
            JOIN users_in_groups uig ON g.id = uig.groupid
            WHERE uig.userid = $1
        `, userID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer rows.Close()

	var groups []map[string]string
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.Name, &group.ID); err != nil {
			log.Println(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
		}
		id_string := strconv.Itoa(group.ID)
		groupMap := map[string]string{
			"id":   id_string,
			"name": group.Name,
		}
		groups = append(groups, groupMap)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	return ctx.JSON(groups)
}

func (a *Api) addNewGroup(ctx *fiber.Ctx) error {
	log.Println("addNewGroup")
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

	var group Group
	if err := ctx.BodyParser(&group); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	if ok, err := checkCorrectWord(group.Name, "Имя"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if ok, err := checkCorrectWord(group.Code, "Код"); !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	tx, err := a.db.BeginTx(ctx.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer tx.Rollback()

	row := tx.QueryRow(`
        SELECT id
        FROM groups
        WHERE code = $1
    `, group.Code)
	var groupWithCode int
	if err := row.Scan(&groupWithCode); err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	if groupWithCode != 0 {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Группа с такми кодом уже есть"})
	}

	err = tx.QueryRow(`
			INSERT INTO groups (name, code)
			VALUES ($1, $2)
			RETURNING id
		`, group.Name, group.Code).Scan(&group.ID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	row = tx.QueryRow(`
			INSERT INTO users_in_groups (userid, groupid)
			VALUES ($1, $2)
		`, userID, group.ID)
	if err := row.Err(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	tx.Commit()

	log.Println("Group", group.ID, "created")

	return ctx.JSON(fiber.Map{"id": group.ID, "name": group.Name})
}

func (a *Api) joinGroupByCode(ctx *fiber.Ctx) error {
	log.Println("joinGroupByCode")
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

	var group Group
	if err := ctx.BodyParser(&group); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	tx, err := a.db.BeginTx(ctx.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	row := tx.QueryRow(`
			SELECT id, name
			FROM groups
			WHERE code = $1
		`, group.Code)
	if err := row.Scan(&group.ID, &group.Name); err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Группа не найдена"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	row = tx.QueryRow(`
        SELECT id
        FROM users_in_groups
        WHERE userid = $1 AND groupid = $2
    `, userID, group.ID)
	var userInGroupID int
	if err := row.Scan(&userInGroupID); err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	if userInGroupID != 0 {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Пользователь уже состоит в группе"})
	}

	row = tx.QueryRow(`
			INSERT INTO users_in_groups (userid, groupid)
			VALUES ($1, $2)
		`, userID, group.ID)
	if err := row.Err(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("User", userID, "joined group", group.ID)

	return ctx.JSON(fiber.Map{"id": group.ID, "name": group.Name})
}

func (a *Api) deleteGroup(ctx *fiber.Ctx) error {
	log.Println("deleteGroup")
	sess, err := a.session.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	userID := sess.Get("id")
	role := sess.Get("role")
	log.Println(userID, role)
	if userID == nil || role == nil {
		log.Println("unauthorized")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неавторизованный запрос"})
	}

	if role != "teacher" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Недостаточно прав"})
	}

	groupID := ctx.Params("groupid")
	if groupID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	tx, err := a.db.BeginTx(ctx.Context(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        DELETE FROM users_in_groups
        WHERE groupid = $1
    `, groupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	_, err = tx.Exec(`
   		DELETE FROM submissions
    	WHERE taskid IN (SELECT t.id FROM tasks t JOIN contests c ON c.id = t.contestid WHERE c.groupid = $1)
	`, groupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	_, err = tx.Exec(`
   		DELETE FROM tasks
    	WHERE contestid IN (SELECT id FROM contests WHERE groupid = $1)
	`, groupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	_, err = tx.Exec(`
        DELETE FROM contests
        WHERE groupid = $1
    `, groupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	_, err = tx.Exec(`
        DELETE FROM groups
        WHERE id = $1
    `, groupID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Group", groupID, "deleted")

	return ctx.JSON(fiber.Map{"success": true})
}
