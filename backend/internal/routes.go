package api

func (a *Api) Register() {
	a.app.Get("/login", a.login)
	a.app.Get("/groups", a.fetchGroupsByUserId)
	a.app.Post("/groups", a.addNewGroup)
	a.app.Post("/join_group", a.joinGroupByCode)
	a.app.Delete("/groups/:groupid", a.deleteGroup)
	a.app.Get("/contests_by_groupid", a.fetchContestsByGroupId)
	a.app.Post("/contests", a.addNewContest)
	a.app.Get("/contests/:contestid", a.fetchContestById)
	a.app.Delete("/contests/:contestid", a.deleteContest)
	a.app.Get("/submissions/unchecked_by_groupid/:groupid", a.fetchUncheckedSubmissionsByGroupId)
	a.app.Get("/submissions/unchecked_by_contestid/:contestid", a.fetchUncheckedSubmissionsByContestId)
	a.app.Get("/submissions_picture/:submissionid", a.fetchSubmissionPictureById)
	a.app.Get("/submissions_details/:submissionid", a.fetchSubmissionDetailsById)
	a.app.Patch("/submissions/accept/:submissionid", a.acceptSubmission)
	a.app.Patch("/submissions/decline/:submissionid", a.declineSubmission)
	a.app.Get("/tasks/:contestid", a.fetchTasksByContestID)
	a.app.Delete("/tasks/:taskid", a.deleteTask)
	a.app.Post("/tasks", a.addTask)
	a.app.Post("/tasks/:taskid/upload", a.uploadSubmission)
	a.app.Post("/register", a.register)
}
