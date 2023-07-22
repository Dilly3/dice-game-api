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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {

	user := db.User{}

	user.Password = "test"
	user.Username = "test"
	user.ID = 4
	loginRequest := db.LoginDto{
		Username: user.Username,
		Password: user.Password,
	}
	loginRequest2 := db.LoginDto{
		Username: "hello",
		Password: "hello",
	}

	tests := []struct {
		name       string
		Address    string
		method     string
		status     int
		body       any
		assertions func(*util.ResponseDto, int, string)
		Expected   string
		Stubs      func(*mocks.MockIGameService, service.IGameService)
	}{
		{
			body:     loginRequest,
			name:     "login success",
			Address:  "/login",
			method:   http.MethodPost,
			status:   200,
			Expected: "login successful",
			Stubs: func(mockservice *mocks.MockIGameService, serv service.IGameService) {
				mockservice.EXPECT().GetUserByUsername(context.Background(), loginRequest.Username).DoAndReturn(func(arg0 context.Context, arg string) (db.User, error) {

					dbUser, err := serv.GetUserByUsername(context.Background(), user.Username)
					assert.NoError(t, err)
					assert.Equal(t, user.Username, dbUser.Username)

					return dbUser, err
				}).AnyTimes()

			},
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, message, resp.Message)
				require.Nil(t, resp.Errors)
				require.Equal(t, code, resp.Status)

			},
		},

		{
			body:     loginRequest2,
			name:     "login fail",
			Address:  "/login",
			method:   http.MethodPost,
			status:   fiber.StatusBadRequest,
			Expected: "username or password incorrect",
			Stubs: func(mockservice *mocks.MockIGameService, serv service.IGameService) {
				mockservice.EXPECT().GetUserByUsername(context.Background(), loginRequest2.Username).DoAndReturn(func(arg0 context.Context, arg string) (db.User, error) {

					dbUser, err := serv.GetUserByUsername(context.Background(), loginRequest2.Username)
					assert.Error(t, err)
					assert.Zero(t, dbUser.ID)

					return dbUser, err
				}).AnyTimes()

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
	tx, err := db.SetupTestDb("../.env")
	if err != nil {
		t.Fail()
	}
	srv := service.NewGameService(tx)

	FiberEngine.Post("/login", Login())
	srv.CreateUser(db.RegisterUserDto{
		Username:        "test",
		Firstname:       "test",
		Lastname:        "test",
		Password:        "test",
		ConfirmPassword: "test",
	})

	defer func() {
		tx.DeleteUser(context.Background(), "test")
		tx.DeleteWallet(context.Background(), "test")
	}()

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			tc.Stubs(sv, srv)
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
