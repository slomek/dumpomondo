package main

import (
	"log"
	"strconv"
	"strings"
)

func defaultDumpSports() string {
	return strings.Join(
		[]string{
			strconv.Itoa(sportRunning),
			strconv.Itoa(sportCycling2),
			strconv.Itoa(sportCycling),
			strconv.Itoa(sportHiking),
			strconv.Itoa(sportWalking),
		},
		",",
	)
}

func mustParseSportsFlag(str string) []int {
	ss := strings.Split(str, ",")

	out := make([]int, 0, len(ss))
	for _, s := range ss {
		id, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Failed to convert sport ID %q to a number: %v", s, err)
		}

		out = append(out, id)
	}

	return out
}

const (
	sportRunning          = 0
	sportCycling2         = 1
	sportCycling          = 2
	sportHiking           = 16
	sportWalking          = 18
	sportSwimming         = 20
	sportIndoorCycling    = 21
	sportOther            = 22
	sportAerobics         = 23
	sportBasketball       = 26
	sportFootball         = 35
	sportPilates          = 38
	sportTableTennis      = 42
	sportWeight           = 46
	sportYoga             = 47
	sportGymnastics       = 49
	sportCircuitTraining2 = 58
	sportCircuitTraining  = 87
	sportTreadmillRunning = 88
	sportIndoorRowing     = 98
	sportStretching       = 103
)

var all = []int{
	sportRunning,
	sportCycling2,
	sportCycling,
	sportHiking,
	sportWalking,
	sportSwimming,
	sportIndoorCycling,
	sportOther,
	sportAerobics,
	sportBasketball,
	sportFootball,
	sportPilates,
	sportTableTennis,
	sportWeight,
	sportYoga,
	sportGymnastics,
	sportCircuitTraining2,
	sportCircuitTraining,
	sportTreadmillRunning,
	sportIndoorRowing,
	sportStretching,
}

func knownSport(sID int) bool {
	return contains(all, sID)
}

func contains(all []int, el int) bool {
	for _, i := range all {
		if i == el {
			return true
		}
	}
	return false
}
