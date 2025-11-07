package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/w-h-a/demo-go/api/user"
	"github.com/w-h-a/demo-go/cmd"
	userservice "github.com/w-h-a/demo-go/internal/service/user"
)

func TestUser_Integration(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) == 0 {
		t.Log("SKIPPING INTEGRATION TEST")
		return
	}

	ur, err := cmd.InitUserRepo("postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable")
	require.NoError(t, err)

	n, err := cmd.InitNotifier()
	require.NoError(t, err)

	userService := userservice.New(ur, n)
	err = userService.Start()
	require.NoError(t, err)
	defer userService.Stop()

	srv, err := cmd.InitHttpServer(":4000", userService)
	require.NoError(t, err)
	err = srv.Start()
	require.NoError(t, err)
	defer srv.Stop()

	t.Run("CreateUser_Success", func(t *testing.T) {
		// Arrange
		body := `{"name":"Integration Test", "email":"integ@test.com"}`
		req, _ := http.NewRequest("POST", "http://localhost:4000/api/users", strings.NewReader(body))
		var u user.User

		// Act
		rsp, err := http.DefaultClient.Do(req)
		defer rsp.Body.Close()

		// Assert
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rsp.StatusCode)
		err = json.NewDecoder(rsp.Body).Decode(&u)
		assert.NotEmpty(t, u.ID)
		assert.Equal(t, "Integration Test", u.Name)
		assert.Equal(t, "integ@test.com", u.Email)
	})
}
