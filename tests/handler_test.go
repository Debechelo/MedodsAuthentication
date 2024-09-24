package tests

var mockRefreshTokenDB = map[string]struct {
	hashedToken string
	ipAddress   string
}{
	"user1": {hashedToken: "2a$10$uQw.fLw7SyZmvjO6ygzIpef4Je1lW8vlLPxmuIfPr25PzAXL8gK7C", ipAddress: "127.0.0.1"},
}
