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

func TestRegister(t *testing.T) {

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
			body: &db.RegisterUserDto{
				Firstname:       "michael",
				Lastname:        "meghan",
				Username:        "mickgo",
				Password:        "test",
				ConfirmPassword: "test",
			},
			name:     "success",
			Address:  "/register",
			status:   201,
			method:   http.MethodPost,
			Expected: "user created successfully",
			Stubs: func(mockservice *mocks.MockIGameService) {
				mockservice.EXPECT().CreateUser(gomock.Any()).AnyTimes().Return(db.User{}, nil)
				mockservice.EXPECT().GetUserByUsername(context.Background(), gomock.Any()).AnyTimes().Return(db.User{}, nil)
			},
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
		},
		{
			body: &db.RegisterUserDto{
				Firstname:       "michael",
				Lastname:        "meghan",
				Username:        "mickgo",
				Password:        "test",
				ConfirmPassword: "teeet",
			},
			name:     "password mismatch",
			Address:  "/register",
			status:   400,
			method:   http.MethodPost,
			Expected: "passwords do not match",
			Stubs: func(mockservice *mocks.MockIGameService) {
				mockservice.EXPECT().CreateUser(gomock.Any()).AnyTimes().Return(db.User{}, nil)
				mockservice.EXPECT().GetUserByUsername(context.Background(), gomock.Any()).AnyTimes().Return(db.User{}, nil)
			},
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
		},
		{
			body:     &db.RegisterUserDto{},
			name:     "missing fields",
			Address:  "/register",
			status:   400,
			method:   http.MethodPost,
			Expected: "some fields missing in user data",
			Stubs: func(mockservice *mocks.MockIGameService) {
				mockservice.EXPECT().CreateUser(gomock.Any()).AnyTimes().Return(db.User{}, nil)
				mockservice.EXPECT().GetUserByUsername(context.Background(), gomock.Any()).AnyTimes().Return(db.User{}, nil)
			},
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
		},
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	FiberEngine = fiber.New()
	sv := mocks.NewMockIGameService(ctrl)

	service.DefaultGameService = sv

	FiberEngine.Post("/register", Register())

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
