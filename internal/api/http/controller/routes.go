package controller

const (
	PathAuthRegister = "/api/v1/auth/register"
	PathAuthLogin    = "/api/v1/auth/login"

	PathUserMe = "/api/v1/users/me"

	PathFamily        = "/api/v1/family"
	PathFamilyMembers = "/api/v1/family/members"

	PathTasks   = "/api/v1/tasks"
	PathTaskAdd = "/api/v1/task/add"

	PathSeasons = "/api/v1/seasons"
	PathSeason  = "/api/v1/seasons/{season_id}"

	PathSeasonTasks       = "/api/v1/seasons/{season_id}/tasks"
	PathSeasonActivate    = "/api/v1/seasons/{season_id}/activate"
	PathSeasonComplete    = "/api/v1/seasons/{season_id}/complete"
	PathSeasonResults     = "/api/v1/seasons/{season_id}/results"
	PathSeasonCompletions = "/api/v1/seasons/{season_id}/completions"

	PathSeasonTaskSubmit      = "/api/v1/season-tasks/{season_task_id}/submit"
	PathTaskCompletionApprove = "/api/v1/task-completions/{completion_id}/approve"
	PathTaskCompletionReject  = "/api/v1/task-completions/{completion_id}/reject"
)
