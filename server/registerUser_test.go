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
	_ "github.com/lib/pq"
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
		Stubs      func(*mocks.MockIGameService, *db.PGXDB, service.IGameService)
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
			Stubs: func(mockservice *mocks.MockIGameService, dbx *db.PGXDB, service service.IGameService) {
				defer func() {
					dbx.DeleteWallet(context.Background(), "mickgo")
					dbx.DeleteUser(context.Background(), "mickgo")
				}()
				mockservice.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(arg db.RegisterUserDto) (db.User, error) {
					userParam := db.RegisterUserDto{
						Firstname:       "michael",
						Lastname:        "meghan",
						Username:        "mickgo",
						Password:        "test",
						ConfirmPassword: "test",
					}
					dbUser, err := service.CreateUser(userParam)
					require.Nil(t, err)
					require.Equal(t, dbUser.Username, "mickgo")
					dbUser, err = dbx.GetUserByUsername(context.Background(), "mickgo")
					require.Nil(t, err)
					require.Equal(t, dbUser.Username, "mickgo")
					//dbUser, _, err := dbx.CreateUserTX(context.Background(), userParam)
					return dbUser, err
				}).AnyTimes()

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
				Firstname:       "michael0",
				Lastname:        "meghan0",
				Username:        "mickgo12",
				Password:        "test",
				ConfirmPassword: "teeet",
			},
			name:     "password mismatch",
			Address:  "/register",
			status:   400,
			method:   http.MethodPost,
			Expected: "passwords do not match",
			Stubs: func(mockservice *mocks.MockIGameService, dbx *db.PGXDB, service service.IGameService) {
				mockservice.EXPECT().CreateUser(gomock.Any()).AnyTimes().DoAndReturn(func(arg db.RegisterUserDto) (db.User, error) {
					userParam := db.RegisterUserDto{
						Firstname:       "michael0",
						Lastname:        "meghan0",
						Username:        "mickgo12",
						Password:        "test",
						ConfirmPassword: "teeet",
					}
					dbUser, err := service.CreateUser(userParam)

					return dbUser, err

				})
				mockservice.EXPECT().GetUserByUsername(context.Background(), gomock.Any()).DoAndReturn(func(arg0 context.Context, arg string) (db.User, error) {
					dbUser, err := service.GetUserByUsername(context.Background(), "mickgo12")
					require.Nil(t, err)
					require.Zero(t, dbUser.ID)
					require.Empty(t, dbUser.Username)

					return dbUser, err
				}).AnyTimes()
			},
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
		},
		{
			body: &db.RegisterUserDto{

				Username: "mickgo23",
			},
			name:     "missing fields",
			Address:  "/register",
			status:   400,
			method:   http.MethodPost,
			Expected: "some fields missing in user data",
			Stubs: func(mockservice *mocks.MockIGameService, dbx *db.PGXDB, service service.IGameService) {
				mockservice.EXPECT().CreateUser(gomock.Any()).AnyTimes().DoAndReturn(func(arg db.RegisterUserDto) (db.User, error) {
					userParam := db.RegisterUserDto{}
					dbUser, err := service.CreateUser(userParam)
					return dbUser, err
				})
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
	var err error
	var tx *db.PGXDB
	service.DefaultGameService = sv
	tx, err = db.SetupTestDb("../.env")
	if err != nil {
		t.Fail()
	}

	srv := service.NewGameService(tx)
	FiberEngine.Post("/register", Register())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			tc.Stubs(sv, tx, srv)
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
