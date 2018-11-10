package auth

func VerifyToken(token string) bool {
	return token == "aronic"
}
