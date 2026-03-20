package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeltaCalcInRange(t *testing.T) {
	assert := assert.New(t)
	matchPayload := MatchPayload{
		UserID: 3,
		Result: "win",
		Crowns: 2,
	}

	testDeltaResult := calcDelta(matchPayload)
	minDelta := 26
	maxDelta := 34
	assert.GreaterOrEqual(testDeltaResult, int64(minDelta), "(%v) not greater than or equal to minDelta", testDeltaResult)
	assert.LessOrEqual(testDeltaResult, int64(maxDelta), "(%v) not less or equal tot maxDelta", testDeltaResult)
}

func TestDeltaOnLossPayload(t *testing.T) {
	assert := assert.New(t)

	mp := MatchPayload{
		UserID: 3,
		Result: "loss",
		Crowns: 2,
	}

	testDeltaResult := calcDelta(mp)

	assert.Negative(testDeltaResult)
}

// func TestErrors(t *testing.T) {
// 	assert := assert.New(t)
//
// 	mpOne := MatchPayload{
// 		UserID: 3,
// 		Result: "loss",
// 		Crowns: 3,
// 	}
//
// 	mpTwo := MatchPayload{
// 		UserID: 3,
// 		Result: "win",
// 		Crowns: 0,
// 	}
//
//
//
// }
