package server

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/mocks"
	"github.com/dilly3/dice-game-api/service"
	"github.com/dilly3/dice-game-api/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {

	user := db.User{}

	user.Password = "test"
	user.Username = "mickgo"
	db.HashPassword(&user)
	user.ID = 4
	loginRequest := db.LoginDto{
		Username: user.Username,
		Password: "test",
	}

	tests := []struct {
		name       string
		Address    string
		method     string
		status     int
		body       any
		assertions func(*util.ResponseDto, int, string)
		Expected   string
		Stubs      func(*mocks.MockIGameService)
	}{
		{
			body:     loginRequest,
			name:     "login success",
			Address:  "/login",
			method:   http.MethodPost,
			status:   200,
			Expected: "login successful",
			Stubs: func(mockservice *mocks.MockIGameService) {
				mockservice.EXPECT().GetUserByUsername(context.Background(), loginRequest.Username).AnyTimes().Return(user, nil)

			},
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, message, resp.Message)
				require.Nil(t, resp.Errors)
				require.Equal(t, code, resp.Status)

			},
		},
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	FiberEngine = fiber.New()
	sv := mocks.NewMockIGameService(ctrl)

	service.DefaultGameService = sv

	FiberEngine.Post("/login", Login())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			tc.Stubs(sv)
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
