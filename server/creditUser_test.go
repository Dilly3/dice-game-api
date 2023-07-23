package server

import (
	"encoding/json"
	"net/http"
	"testing"

	db "github.com/dilly3/dice-game-api/db/sqlc"

	"github.com/dilly3/dice-game-api/util"
	"github.com/stretchr/testify/require"
)

func Test_CreditWallet(t *testing.T) {

	var err error
	var tx *db.PGXDB

	defer deleteTestUsers([]string{"testUser7", "testUser4"})
	createuserandcredit("testUser4", 70)
	createuser("testUser7")

	tests := []struct {
		name       string
		Address    string
		method     string
		status     int
		body       any
		assertions func(*util.ResponseDto, int, string)
		Expected   string
		username   string
	}{
		{
			name:     "success",
			username: "testUser7",
			Address:  "/credit",
			method:   http.MethodGet,
			status:   200,
			Expected: "wallet credited",
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Equal(t, message, resp.Message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)
			},
		},

		{
			name:     "cant credit",
			Address:  "/credit",
			method:   http.MethodGet,
			username: "testUser4",
			status:   500,
			Expected: "you still have up to 35 sats",
			assertions: func(resp *util.ResponseDto, code int, message string) {
				require.Contains(t, resp.Message, message)
				require.Nil(t, resp.Errors)
				require.Equal(t, resp.Status, code)

			},
		},
	}

	tx, err = db.SetupTestDb("../.env")

	if err != nil {
		t.Fail()
	}
	defer tx.DB.Close()

	FiberEngine.Get("/credit", CreditWallet())

	for i := 0; i < len(tests); i++ {
		//time.Sleep(time.Second * 3)
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {

			resp, err := ExecuteRequest(tc.method, tc.Address, tc.body, tc.username)
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
