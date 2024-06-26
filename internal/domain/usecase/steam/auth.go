package usecasesteam

func (u *steam) SteamAuth() error {
	return u.r.SteamLogin()
}
