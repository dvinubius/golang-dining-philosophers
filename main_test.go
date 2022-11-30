package main

import (
	"testing"
	"time"
)

const times = 10

// Requirement: no philosopher should starve
// we check that, at the end, the history of table-leavers contains all philosophers
func Test_dine(t *testing.T) {
	run(t, 0, 0)
}

func Test_dine_varyingDurations(t *testing.T) {
	run(t, 25, 25)
	run(t, 250, 250)
	run(t, 0, 25)
	run(t, 25, 0)
}

func run(t *testing.T, eatMs int, thinkMs int) {
	eatTime = time.Duration(eatMs) * time.Millisecond
	thinkTime = time.Duration(thinkMs) * time.Millisecond
	for i := 0; i < times; i++ {
		historyLeave = []int{}
		dine()
		if len(historyLeave) != len(philosophers) {
			t.Errorf("Not all philosophers have eaten enough. Eat time %v\tThink time %v", eatTime, thinkTime)
		}
	}
}
