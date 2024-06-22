package user

type RegReq struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Pwd      string `form:"pwd" binding:"required"`
}

type RegResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Email string `form:"email" binding:"required,email"`
	Pwd   string `form:"pwd" binding:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResp struct {
	AccessToken string `json:"access_token"`
}
