package family

type createRequest struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type familyResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
	OwnerID  string `json:"owner_id"`
	Role     string `json:"role"`
}

type createMemberRequest struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type memberResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

type membersResponse struct {
	Items []memberResponse `json:"items"`
}
