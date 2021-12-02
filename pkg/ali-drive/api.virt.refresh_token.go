package ali_drive

func (y *Ali) RefreshToken(refreshToken string) *ApiAcquireTokenResponse {
	return y.AcquireToken(&ApiAcquireTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
}
