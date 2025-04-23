package minLogic

import (
	"math/rand"
	"time"
)

type Hello struct {
	Hello string
}

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetRandomGPS() int {
	return rand.Intn(500000000) + 100000000
}
