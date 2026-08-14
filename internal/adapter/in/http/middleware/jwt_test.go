package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupProtectedRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", RequireAuth(secret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return router
}

func signToken(secret string, expiresAt time.Time, method jwt.SigningMethod) string {
	claims := jwt.MapClaims{"sub": 1, "exp": expiresAt.Unix()}
	token := jwt.NewWithClaims(method, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestRequireAuth(t *testing.T) {
	secret := "test-secret"
	router := setupProtectedRouter(secret)

	t.Run("Missing Header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Malformed Header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "not-a-bearer-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		token := signToken("wrong-secret", time.Now().Add(time.Hour), jwt.SigningMethodHS256)

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Expired Token", func(t *testing.T) {
		token := signToken(secret, time.Now().Add(-time.Hour), jwt.SigningMethodHS256)

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Valid Token", func(t *testing.T) {
		token := signToken(secret, time.Now().Add(time.Hour), jwt.SigningMethodHS256)

		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
