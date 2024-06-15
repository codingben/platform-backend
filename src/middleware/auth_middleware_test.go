package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockTokenProvider struct {
	Username string
	Token    string
	Err      error
}

func (m MockTokenProvider) ObtainToken(username, password string, logger *zap.Logger, ctx *gin.Context) (string, error) {
	return m.Token, m.Err
}

func (m MockTokenProvider) ObtainUsername(token string, logger *zap.Logger) (string, error) {
	return m.Username, m.Err
}

func TestTokenAuthMiddleware(t *testing.T) {
	_ = os.Setenv(envKubeAPIServer, "https://example.com/api")
	_ = os.Setenv(envInsecureSkipVerify, "true")

	logger, _ := zap.NewDevelopment()
	router := gin.New()

	// Set up a middleware that injects the logger into the context
	router.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	router.Use(TokenAuthMiddleware(MockTokenProvider{Token: "valid_token", Username: "user", Err: nil}))
	router.GET("/ping", func(c *gin.Context) {
		_, ok := c.Get("kubeClient")
		if !ok {
			t.Error("Expected kubeClient to be set in context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "kubeClient not set in context"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	type args struct {
		authHeader string
	}
	type want struct {
		expectedStatus int
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ShouldSuccessWithValidToken": {
			args: args{
				authHeader: httpBearerTokenPrefix + " valid_token",
			},
			want: want{
				expectedStatus: http.StatusOK,
			},
		},
		"ShouldFailWithoutAuthorizationHeader": {
			args: args{
				authHeader: "",
			},
			want: want{
				expectedStatus: http.StatusUnauthorized,
			},
		},
		"ShouldFailWithInvalidAuthorizationHeader": {
			args: args{
				authHeader: "InvalidToken invalid_token",
			},
			want: want{
				expectedStatus: http.StatusUnauthorized,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			if tc.args.authHeader != "" {
				req.Header.Set(httpAuthorizationHeader, tc.args.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.want.expectedStatus {
				t.Errorf("Expected status code %d; got %d", tc.want.expectedStatus, w.Code)
			}
		})
	}
}