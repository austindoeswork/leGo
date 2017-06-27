package entity

import (
	"math/rand"
	"time"

	"github.com/wardn/uuid"
)

type Identifier interface {
	ID() string
	TypeID() string
}

func UUID() string {
	return uuid.NewNoDash()
}

func Now() time.Time {
	return time.Now().UTC()
}

/////////////////////////////////////////////////////////
// RANDOM
/////////////////////////////////////////////////////////

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RANDstring() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RANDfloat64() float64 {
	return rand.Float64()
}

func RANDint() int {
	return rand.Int()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
