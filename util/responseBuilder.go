package util

import "github.com/dilly3/dice-game-api/game"

type RollResponseBuilder struct {
	Message   string         `json:"message,omitempty"`
	Gamescore game.GameScore `json:"gameScore,omitempty"`
	JackpotNumber int         `json:"jackpot_number,omitempty"`
}

func NewRollResponseBuilder() *RollResponseBuilder {
	return &RollResponseBuilder{}
}
func (r *RollResponseBuilder) SetMessage(message string) *RollResponseBuilder {
	r.Message = message
	return r
}

func (r *RollResponseBuilder) SetGameScore(score1, score2 int) *RollResponseBuilder {
	r.Gamescore = game.GameScore{
		RollNumber1: score1,
		RollNumber2: score2,
	}
	return r
}

func (r *RollResponseBuilder) SetJackpotNumber(jackpotNumber int) *RollResponseBuilder {
	r.JackpotNumber = jackpotNumber
	return r
}

func (r *RollResponseBuilder) Build() RollResponseDto {
	return RollResponseDto{
		Message:   r.Message,
		Gamescore: r.Gamescore,
		JackpotNumber: r.JackpotNumber,
	}
}

type ResponseDtoBuilder struct {
	Message       string      `json:"message,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	Errors        []string    `json:"errors,omitempty"`
	Status        int         `json:"status,omitempty"`
	JackpotNumber int         `json:"jackpot_number,omitempty"`
	TimeStamp     string      `json:"timestamp,omitempty"`
	Balance       string      `json:"balance,omitempty"`
	Assts         string      `json:"assets,omitempty"`
}

func NewResponseDtoBuilder() *ResponseDtoBuilder {
	return &ResponseDtoBuilder{}
}

func (r *ResponseDtoBuilder) SetMessage(message string) *ResponseDtoBuilder {
	r.Message = message
	return r
}

func (r *ResponseDtoBuilder) SetData(data interface{}) *ResponseDtoBuilder {
	r.Data = data
	return r
}

func (r *ResponseDtoBuilder) SetErrors(errors []string) *ResponseDtoBuilder {
	r.Errors = errors
	return r
}

func (r *ResponseDtoBuilder) SetStatus(status int) *ResponseDtoBuilder {
	r.Status = status
	return r
}

func (r *ResponseDtoBuilder) SetJackpotNumber(jackpotNumber int) *ResponseDtoBuilder {
	r.JackpotNumber = jackpotNumber
	return r
}

func (r *ResponseDtoBuilder) SetTimeStamp(timeStamp string) *ResponseDtoBuilder {
	r.TimeStamp = timeStamp
	return r
}

func (r *ResponseDtoBuilder) SetBalance(bal string) *ResponseDtoBuilder {
	r.Balance = bal
	return r
}

func (r *ResponseDtoBuilder) SetAssets(assets string) *ResponseDtoBuilder {
	r.Assts = assets
	return r
}

func (r *ResponseDtoBuilder) Build() *ResponseDto {
	return &ResponseDto{
		Message:       r.Message,
		Data:          r.Data,
		Errors:        r.Errors,
		Status:        r.Status,
		JackpotNumber: r.JackpotNumber,
		TimeStamp:     r.TimeStamp,
	}
}
