package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

func (s *apiserver) webAdminReviewList(c *gin.Context) {
	const logPref = "webAdminReviewItem"

	qlp := &representation.QueryListParams{}

	if err := c.BindQuery(qlp); err != nil {
		log.Printf("(ERR) %s: bind query error: %v", logPref, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	var (
		username string
		role     = iRole.(string)
	)
	if role == UserRoleExpert {

		iUserID, ok := c.Get("userID")
		if !ok {
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}

		userID := iUserID.(string)
		if userID != qlp.ExpertID {
			qlp.ExpertID = userID
		}

		expert, err := s.repository.GetExpert(userID)
		if err != nil {
			log.Printf("(ERR) %s: failed to fetch expeert details: %v", logPref, err)
		}
		username = expert.Username
	}

	list, err := s.repository.GetReviewList(representation.QueryReviewAPItoPersistence(qlp))
	if err != nil {
		log.Printf("(ERR) %s: failed to fetch review list: %v", logPref, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "review_list.html",
		gin.H{
			"metadata": gin.H{
				"title":     "Відгуки",
				"username":  username,
				"logged_in": true,
				"role":      role,
				"status":    iStatus.(string),
			},
			"data":       representation.ReviewListPersistenceToAPI(list),
			"pagination": qlp.GeneratePagination(list.Total),
			"queries": gin.H{
				"status":    qlp.Status,
				"expert_id": qlp.ExpertID,
			},
		},
	)
}

func (s *apiserver) webAdminReviewItem(c *gin.Context) {
	const logPref = "webAdminReviewItem"

	reviewID := c.Param("reviewId")
	reviewRes, err := s.repository.GetReview(reviewID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: review [%s] not found", logPref, reviewID)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get review view [%s]: %v", logPref, reviewID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	role := iRole.(string)
	if role == UserRoleExpert {

		iUserID, ok := c.Get("userID")
		if !ok {
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}

		if iUserID.(string) != reviewRes.ExpertID {
			c.Redirect(http.StatusFound, "/admin/profile")
			return
		}
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "review_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      role,
				"status":    iStatus.(string),
			},
			"review": representation.ReviewItemPersistenceToAPI(reviewRes),
		},
	)
}
