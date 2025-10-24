package wxapi

type GetAccessTokenResult struct {
	*Result
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
