package server

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/util"
	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	defer deleteTestUsers([]string{"testUser4"})
	createuser("testUser4")

	tests := []struct {
		name       string
		Address    string
		method     string
		status     int
		body       any
		assertions func(*util.ResponseDto, int, string)
		Expected   string
		delete     func(string)
	}{
		{
			body: db.LoginDto{
				Username: "testUser4",
				Password: "test",
			},
			name:     "login success",
			Address:  "/login",
			method:   http.MethodPost,
			status:   200,
			Expected: "login successful",

			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, message, resp.Message)
				require.Nil(t, resp.Errors)
				require.Equal(t, code, resp.Status)

			},
			delete: func(username string) {
				db.DefaultGameRepo.DeleteTransactionByUsername(context.Background(), username)
				db.DefaultGameRepo.DeleteWallet(context.Background(), username)
				db.DefaultGameRepo.DeleteUser(context.Background(), username)

			},
		},

		{
			body: db.LoginDto{
				Username: "Hello",
				Password: "Hello",
			},
			name:     "login fail",
			Address:  "/login",
			method:   http.MethodPost,
			status:   fiber.StatusBadRequest,
			Expected: "username or password incorrect",

			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, message, resp.Message)
				require.Nil(t, resp.Errors)
				require.Equal(t, code, resp.Status)

			},
		},
		{
			body: db.RegisterUserDto{
				Firstname: "hello",
				Lastname:  "hello",
			},
			name:     "bad credential",
			Address:  "/login",
			method:   http.MethodPost,
			status:   fiber.StatusBadRequest,
			Expected: "bad login credentials",

			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, message, resp.Message)
				require.Nil(t, resp.Errors)
				require.Equal(t, code, resp.Status)

			},
		},
	}

	FiberEngine.Post("/login", Login())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			resp, err := ExecuteRequest(tc.method, tc.Address, tc.body, "")
			if err != nil {
				t.Fail()
			}

			res := &util.ResponseDto{}
			err = json.Unmarshal(resp, res)
			if err != nil {
				t.Fail()
			}
			tc.assertions(res, tc.status, tc.Expected)

		})
	}

}
