package util

import (
	"time"

	"github.com/dilly3/dice-game-api/game"
)

type ResponseBuilder struct {
	Message       string
	Data          any
	Errors        []string
	Status        int
	JackpotNumber int
	TimeStamp     string
	Balance       string
	Assts         string
	GameScore     any
	SessionState  bool
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{GameScore: nil, Errors: nil, Data: nil}
}

func (r *ResponseBuilder) SetMessage(message string) *ResponseBuilder {
	r.Message = message
	return r
}

func (r *ResponseBuilder) SetData(data interface{}) *ResponseBuilder {
	r.Data = data
	return r
}

func (r *ResponseBuilder) SetErrors(errors []string) *ResponseBuilder {
	r.Errors = errors
	return r
}

func (r *ResponseBuilder) SetStatus(status int) *ResponseBuilder {
	r.Status = status
	return r
}

func (r *ResponseBuilder) SetJackpotNumber(jackpotNumber int) *ResponseBuilder {
	r.JackpotNumber = jackpotNumber
	return r
}

func (r *ResponseBuilder) SetTimeStamp() *ResponseBuilder {
	r.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	return r
}

func (r *ResponseBuilder) SetBalance(bal string) *ResponseBuilder {
	r.Balance = bal
	return r
}

func (r *ResponseBuilder) SetAssets() *ResponseBuilder {
	r.Assts = "sats"
	return r
}

func (r *ResponseBuilder) SetGameScore(score1, score2 int) *ResponseBuilder {
	r.GameScore = game.GameScore{
		RollNumber1: score1,
		RollNumber2: score2,
	}
	return r
}

func (r *ResponseBuilder) SetSessionState(state bool) *ResponseBuilder {
	r.SessionState = state
	return r
}

func (r *ResponseBuilder) Build() *ResponseDto {
	return &ResponseDto{
		Message:       r.Message,
		Data:          r.Data,
		Errors:        r.Errors,
		Status:        r.Status,
		JackpotNumber: r.JackpotNumber,
		TimeStamp:     r.TimeStamp,
		Balance:       r.Balance,
		Assts:         r.Assts,
		GameScore:     r.GameScore,
		SessionState:  r.SessionState,
	}
}
