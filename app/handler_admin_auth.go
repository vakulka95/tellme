package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *apiserver) webAdminLogin(c *gin.Context) {
	var (
		login    = c.PostForm("login")
		password = c.PostForm("password")
	)

	token, _, err := s.login(login, password)
	if err != nil {
		if err.Error() == ErrUserIsBlocked.Error() {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": "Користувач не активний"})
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Невірний логін або пароль"})
		return
	}

	c.SetCookie(authCookieKey, token, int(s.config.AccessTokenDuration.Seconds()), "", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/admin/requisition")
}

func (s *apiserver) webAdminGetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (s *apiserver) webAdminLogout(c *gin.Context) {
	c.SetCookie(authCookieKey, "", -1, "", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/admin")
}
