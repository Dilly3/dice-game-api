package server

import (
	"context"
	"encoding/json"
	"net/http"

	"testing"

	db "github.com/dilly3/dice-game-api/db/sqlc"
	"github.com/dilly3/dice-game-api/service"
	"github.com/dilly3/dice-game-api/util"

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
		delete     func(string)
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

			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
			delete: func(username string) {
				db.DefaultGameRepo.DeleteWallet(context.Background(), username)
				db.DefaultGameRepo.DeleteUser(context.Background(), username)

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

			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
			delete: func(username string) {
				db.DefaultGameRepo.DeleteWallet(context.Background(), username)
				db.DefaultGameRepo.DeleteUser(context.Background(), username)

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

			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
			delete: func(username string) {
				db.DefaultGameRepo.DeleteWallet(context.Background(), username)
				db.DefaultGameRepo.DeleteUser(context.Background(), username)

			},
		},
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	var err error
	var tx *db.PGXDB

	tx, err = db.SetupTestDb("../.env")
	if err != nil {
		t.Fail()
	}
	//defer tx.DB.Close()
	db.DefaultGameRepo = tx

	service.DefaultGameService = service.NewGameService(tx)
	FiberEngine.Post("/register", Register())

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			//tc.Stubs(sv, tx, srv)
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
			tc.delete(tc.body.(*db.RegisterUserDto).Username)

		})
	}

}
