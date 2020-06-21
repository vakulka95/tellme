package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *apiserver) webAdminLogin(c *gin.Context) {
	var (
		login    = c.PostForm("login")
		password = c.PostForm("password")
		captcha  = c.PostForm("g-recaptcha-response")
	)

	err := s.checkGoogleCaptcha(captcha)
	if err != nil {
		log.Printf("(ERR) Google captcha check failed: %v", err)
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Поставте позначку, щоб увійти в систему"})
		return
	}

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
	c.Redirect(http.StatusFound, "/admin")
}

func (s *apiserver) webAdminGetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (s *apiserver) webAdminLogout(c *gin.Context) {
	c.SetCookie(authCookieKey, "", -1, "", "", false, true)
	c.Redirect(http.StatusFound, "/admin")
}

func (s *apiserver) checkGoogleCaptcha(c string) error {
	req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("secret", s.config.GoogleCaptchaSecret)
	q.Add("response", c)
	req.URL.RawQuery = q.Encode()

	var (
		client         = &http.Client{}
		googleResponse map[string]interface{}
	)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &googleResponse)
	if err != nil {
		return err
	}

	if !(googleResponse["success"].(bool)) {
		return fmt.Errorf("invalid check google captcha response")
	}

	return nil
}
