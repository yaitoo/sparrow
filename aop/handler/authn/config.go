package authn

//Config 登錄驗證設定
type Config struct {
	Authentication map[string][]string
}

//IsIdentityUser 验证是否通过
func (c *Config) IsIdentityUser(user *IdentityUser, route string) bool {
	if c == nil {
		return true
	}

	rules, ok := c.Authentication[route]

	if !ok {
		rules, _ = c.Authentication["*"]
	}

	for _, r := range rules {
		if r == "*" {
			return true
		}

		if r == "?" {
			return user != nil
		}

		if user != nil && user.ID == r {

		}
	}

	return true

}
