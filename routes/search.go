package routes

import (
	"net/http"

	"github.com/Devansh3712/tsuki/database"
	"github.com/Devansh3712/tsuki/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var searchLimit = 10

type search struct {
	models.User
	Followers int
	Following int
	Posts     int
	Follows   any
}

func SearchUser(c *gin.Context) {
	session := sessions.Default(c)
	switch c.Request.Method {
	case "GET":
		searchLimit = 10
		session.Delete("search")
		session.Save()
		c.HTML(http.StatusOK, "search.tmpl.html", nil)
	case "POST":
		id := session.Get("userId")
		if c.PostForm("search") != "" {
			session.Set("search", c.PostForm("search"))
			session.Save()
		}
		keyword := session.Get("search").(string)

		searchResult := database.ReadUsers(keyword, searchLimit)
		var users []search
		for _, result := range searchResult {
			user := search{
				User:      result,
				Followers: database.ReadFollowersCount(result.Id),
				Following: database.ReadFollowingCount(result.Id),
				Posts:     database.ReadPostsCount(result.Id),
			}
			if id != nil && id.(string) != result.Id {
				user.Follows = database.Followed(id.(string), result.Id)
			}
			users = append(users, user)
		}
		c.HTML(http.StatusOK, "search.tmpl.html", gin.H{
			"users": users,
		})
	}
}

func LoadMoreUsers(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("userId")
	keyword := session.Get("search").(string)
	searchLimit += 10

	searchResult := database.ReadUsers(keyword, searchLimit)
	var users []search
	for _, result := range searchResult {
		user := search{
			User:      result,
			Followers: database.ReadFollowersCount(result.Id),
			Following: database.ReadFollowingCount(result.Id),
			Posts:     database.ReadPostsCount(result.Id),
		}
		if id != nil && id.(string) != result.Id {
			user.Follows = database.Followed(id.(string), result.Id)
		}
		users = append(users, user)
	}
	c.HTML(http.StatusOK, "search.tmpl.html", gin.H{
		"users": users,
	})
}

func ToggleSearchFollow(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("userId")
	if id == nil {
		c.HTML(http.StatusUnauthorized, "error.tmpl.html", gin.H{
			"error":   "401 Unauthorized",
			"message": "User not logged in.",
		})
		return
	}
	username := c.Param("username")
	toFollow := database.ReadUserByName(username)
	database.ToggleFollow(id.(string), toFollow.Id)
	c.Redirect(http.StatusPermanentRedirect, "/search")
}