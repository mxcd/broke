package clients

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mxcd/broke/internal/util"
	"github.com/rs/zerolog/log"
)

type JWT struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
}

type KeycloakMockServerConfig struct {
	Port          int                    `yaml:"port"`
	AdminUsername string                 `yaml:"adminUsername"`
	AdminPassword string                 `yaml:"adminPassword"`
	Realm         string                 `yaml:"realm"`
	Data          KeycloakMockServerData `yaml:"data"`
}

type KeycloakMockServerData struct {
	Users  []KeycloakMockServerUser  `json:"users"`
	Groups []KeycloakMockServerGroup `json:"groups"`
}

type KeycloakMockServerUser struct {
	Id            string   `json:"id"`
	Username      string   `json:"username"`
	Enabled       bool     `json:"enabled"`
	EmailVerified bool     `json:"emailVerified"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	Email         string   `json:"email"`
	Groups        []string `json:"groups"`
}

type KeycloakMockServerGroup struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	SubGroups  []string `json:"subGroups"`
	RealmRoles []string `json:"realmRoles"`
}

func StartMockServer(ctx context.Context, config *KeycloakMockServerConfig) *http.Server {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Endpoint to mock GetToken
	router.POST("/realms/:realm/protocol/openid-connect/token", func(c *gin.Context) {

		// get the realm from the request
		realm := c.Param("realm")

		// get the username and password from the request from x-www-form-urlencoded
		username := c.PostForm("username")
		password := c.PostForm("password")
		grantType := c.PostForm("grant_type")

		if grantType != "password" {
			log.Error().Msgf("grant_type '%s' not supported", grantType)
			c.Status(http.StatusBadRequest)
		}

		if username != config.AdminUsername || password != config.AdminPassword {
			log.Error().Msgf("invalid username or password")
			c.Status(http.StatusUnauthorized)
		}

		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusUnauthorized)
		}

		c.JSON(http.StatusOK, JWT{
			AccessToken:      "mock_access_token",
			ExpiresIn:        300,
			RefreshExpiresIn: 600,
			RefreshToken:     "mock_refresh_token",
			TokenType:        "Bearer",
			Scope:            "profile email",
		})
	})

	// Endpoint to mock GetUserCount
	router.GET("/admin/realms/:realm/users/count", func(c *gin.Context) {
		realm := c.Param("realm")
		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusBadRequest)
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, len(config.Data.Users))
	})

	router.GET("/admin/realms/:realm/users", func(c *gin.Context) {
		realm := c.Param("realm")
		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusBadRequest)
		}
		firstString := c.Query("first")
		first, err := strconv.Atoi(firstString)
		if err != nil {
			log.Error().Err(err).Msgf("invalid first")
			c.Status(http.StatusBadRequest)
		}
		maxString := c.Query("max")
		max, err := strconv.Atoi(maxString)
		if err != nil {
			log.Error().Err(err).Msgf("invalid max")
			c.Status(http.StatusBadRequest)
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, util.ListLimitOffset(config.Data.Users, max, first))
	})

	router.GET("/admin/realms/:realm/groups/count", func(c *gin.Context) {
		realm := c.Param("realm")
		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusBadRequest)
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"count": len(config.Data.Groups)})
	})

	// endpoint to get mock groups
	router.GET("/admin/realms/:realm/groups", func(c *gin.Context) {
		// get the realm from the request
		realm := c.Param("realm")
		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusBadRequest)
		}
		firstString := c.Query("first")
		first, err := strconv.Atoi(firstString)
		if err != nil {
			log.Error().Err(err).Msgf("invalid first")
			c.Status(http.StatusBadRequest)
		}
		maxString := c.Query("max")
		max, err := strconv.Atoi(maxString)
		if err != nil {
			log.Error().Err(err).Msgf("invalid max")
			c.Status(http.StatusBadRequest)
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, util.ListLimitOffset(config.Data.Groups, max, first))
	})

	router.GET("/admin/realms/:realm/groups/:groupid", func(c *gin.Context) {
		// get the realm from the request
		realm := c.Param("realm")
		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusBadRequest)
		}

		groupId := c.Param("groupid")
		if groupId == "" {
			log.Error().Msgf("invalid group id")
			c.Status(http.StatusBadRequest)
		}

		var group *KeycloakMockServerGroup = nil
		for _, g := range config.Data.Groups {
			if g.Id == groupId {
				group = &g
				break
			}
		}
		if group == nil {
			log.Error().Msgf("group not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, group)
	})

	router.GET("/admin/realms/:realm/groups/:groupid/members", func(c *gin.Context) {
		// get the realm from the request
		realm := c.Param("realm")
		if realm != config.Realm {
			log.Error().Msgf("invalid realm")
			c.Status(http.StatusBadRequest)
		}

		groupId := c.Param("groupid")
		if groupId == "" {
			log.Error().Msgf("invalid group id")
			c.Status(http.StatusBadRequest)
		}

		var group *KeycloakMockServerGroup = nil
		for _, g := range config.Data.Groups {
			if g.Id == groupId {
				group = &g
				break
			}
		}
		if group == nil {
			log.Error().Msgf("group not found")
			c.Status(http.StatusNotFound)
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, group)
	})

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %s\n", err)
		}
	}()

	util.WaitForServerUp("http://localhost:" + strconv.Itoa(config.Port) + "/health")

	return server
}
