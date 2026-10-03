package task

type addRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Points      int64  `json:"points"`
}

type addResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Points      int64  `json:"points"`
}
