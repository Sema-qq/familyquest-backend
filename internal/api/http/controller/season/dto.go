package season

type createRequest struct {
	Title    string `json:"title"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type seasonResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	Status   string `json:"status"`
}

type seasonsResponse struct {
	Items []seasonResponse `json:"items"`
}

type seasonDetailResponse struct {
	ID       string               `json:"id"`
	Title    string               `json:"title"`
	StartsAt string               `json:"starts_at"`
	EndsAt   string               `json:"ends_at"`
	Status   string               `json:"status"`
	Tasks    []seasonTaskResponse `json:"tasks"`
}

type seasonTaskResponse struct{}
