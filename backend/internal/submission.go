package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Submission struct {
	ID          int    `json:"id"`
	StudentName string `json:"studentName"`
	StudentId   int
	TaskID      int    `json:"taskid"`
	TaskName    string `json:"taskName"`
	GroupID     int    `json:"groupid"`
	ContestID   int    `json:"contestid"`
	Status      string `json:"status"`
}

type SubmissionPhoto struct {
	ID          string
	URL         string
	StudentName string
	TaskName    string
}

func (a *Api) fetchUncheckedSubmissionsByGroupId(ctx *fiber.Ctx) error {
	log.Println("fetchUncheckedSubmissionsByGroupId")
	groupID := ctx.Params("groupid")
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
		SELECT s.id, u.firstname || ' ' || u.surname, s.userid, t.id, t.name, c.groupid, t.contestid, s.status
		FROM submissions s
		JOIN users u ON s.userid = u.id
		JOIN tasks t ON s.taskid = t.id
		JOIN contests c ON t.contestid = c.id
		WHERE c.groupid = $1 AND s.status = 'unchecked'
	`, groupIDInt)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer rows.Close()

	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.StudentName, &s.StudentId, &s.TaskID, &s.TaskName, &s.GroupID, &s.ContestID, &s.Status)
		if err != nil {
			log.Println(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
		}
		submissions = append(submissions, s)
	}

	log.Println(submissions)

	return ctx.JSON(submissions)
}

func (a *Api) fetchUncheckedSubmissionsByContestId(ctx *fiber.Ctx) error {
	log.Println("fetchUncheckedSubmissionsByContestId")
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

	var userInGroupID int
	err = a.db.QueryRow(`
		SELECT id
		FROM users_in_groups ug
		JOIN contests c ON ug.groupid = c.groupid
		WHERE userid = $1 AND c.id = $2
	`, userID, contestIDInt).Scan(&userInGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not in group")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Пользователь не состоит в группе"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	rows, err := a.db.Query(`
		SELECT s.id, u.firstname || ' ' || u.surname, s.userid, t.id, t.name, c.groupid, t.contestid, s.status
		FROM submissions s
		JOIN users u ON s.userid = u.id
		JOIN tasks t ON s.taskid = t.id
		JOIN contests c ON t.contestid = c.id
		WHERE c.id = $1 AND s.status = 'unchecked'
	`, contestIDInt)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}
	defer rows.Close()

	submissions := make([]Submission, 0)
	for rows.Next() {
		var s Submission
		err = rows.Scan(&s.ID, &s.StudentName, &s.StudentId, &s.TaskID, &s.TaskName, &s.GroupID, &s.ContestID, &s.Status)
		if err != nil {
			log.Println(err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
		}
		submissions = append(submissions, s)
	}

	log.Println(submissions)

	return ctx.JSON(submissions)
}

func (a *Api) fetchSubmissionPictureById(ctx *fiber.Ctx) error {
	log.Println("fetchSubmissionPictureById")
	submissionID := ctx.Params("submissionid")
	if submissionID == "" {
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

	var submission SubmissionPhoto
	err = a.db.QueryRow(`
        SELECT s.id, url
        FROM submissions s
        WHERE s.id = $1
    `, submissionID).Scan(&submission.ID, &submission.URL)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Посылка не найдена"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	if _, err := os.Stat(submission.URL); os.IsNotExist(err) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Скриншот не найден"})
	}

	log.Println(submission)

	if err := ctx.SendFile(submission.URL); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось отправить файл"})
	}

	return nil
}

func (a *Api) fetchSubmissionDetailsById(ctx *fiber.Ctx) error {
	log.Println("fetchSubmissionDetailsById")
	submissionID := ctx.Params("submissionid")
	if submissionID == "" {
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

	var submission SubmissionPhoto
	err = a.db.QueryRow(`
        SELECT s.id, u.firstname || ' ' || u.surname, t.name
        FROM submissions s
        JOIN users u on s.userid = u.id
        JOIN tasks t on s.taskid = t.id
        WHERE s.id = $1
    `, submissionID).Scan(&submission.ID, &submission.StudentName, &submission.TaskName)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Посылка не найдена"})
		}
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println(submission)

	return ctx.JSON(fiber.Map{
		"studentName": submission.StudentName,
		"taskName":    submission.TaskName,
	})
}

func (a *Api) acceptSubmission(ctx *fiber.Ctx) error {
	log.Println("acceptSubmission")
	submissionID := ctx.Params("submissionid")
	if submissionID == "" {
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

	_, err = a.db.Exec(`
        UPDATE submissions
        SET status = 'accepted'
        WHERE id = $1
    `, submissionID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Submission", submissionID, "accepted")

	return ctx.JSON(fiber.Map{"success": true})
}

func (a *Api) declineSubmission(ctx *fiber.Ctx) error {
	log.Println("declineSubmission")
	submissionID := ctx.Params("submissionid")
	if submissionID == "" {
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

	_, err = a.db.Exec(`
        UPDATE submissions
        SET status = 'declined'
        WHERE id = $1
    `, submissionID)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка сервера"})
	}

	log.Println("Submission", submissionID, "declined")

	return ctx.JSON(fiber.Map{"success": true})
}

func (a *Api) uploadSubmission(ctx *fiber.Ctx) error {
	log.Println("uploadSubmission")
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
	if role != "student" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Недостаточно прав"})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректный запрос"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Файл не найден"})
	}

	file := files[0]
	if file.Header.Get("Content-Type") != "image/png" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Файл должен быть PNG"})
	}

	dir := "../submission_photos"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось создать директорию"})
	}

	filename := fmt.Sprintf("%d_%s", userID, file.Filename)
	filepath := filepath.Join(dir, filename)
	if err := ctx.SaveFile(file, filepath); err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось сохранить файл"})
	}

	_, err = a.db.Exec(`
        INSERT INTO submissions (userid, taskid, status, url, created_at)
        VALUES ($1, $2, 'unchecked', $3, $4)
    `, userID, taskID, filepath, time.Now())
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Не удалось сохранить посылку в базу данных"})
	}

	log.Println("Submission added:", filename)

	return ctx.JSON(fiber.Map{"success": true})
}
