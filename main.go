package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	email    = os.Getenv("ENDOMONDO_EMAIL")
	password = os.Getenv("ENDOMONDO_PASSWORD")

	gpxSports []int
	skipJSON  bool
	skipGPX   bool
)

func init() {
	gpxSportsStr := flag.String("gpx-sports", defaultDumpSports(), "sport IDs to be dumped to GPX")
	flag.BoolVar(&skipJSON, "skip-json", false, "skip dumping workout data to JSON")
	flag.BoolVar(&skipGPX, "skip-gpx", false, "skip dumping workout data to GPX")

	flag.Parse()

	if !skipGPX {
		gpxSports = mustParseSportsFlag(*gpxSportsStr)
	}
}

func main() {
	userID, cookies := mustLogin(email, password)
	fmt.Printf("== Dumping data for user %d\n", userID)

	ctx := context.Background()

	var (
		dumpDir = filepath.Join("dump", strconv.Itoa(userID))
		jsonDir = filepath.Join(dumpDir, "raw")
		gpxDir  = filepath.Join(dumpDir, "gpx")
	)

	if !skipJSON {
		fmt.Printf("== Dumping data to JSON (%s)\n", jsonDir)
		mustCreateDirectory(jsonDir)
		mustListWorkouts(ctx, userID, 100, 0, cookies)
		fmt.Printf("==\n\n")
	}

	if !skipGPX {
		fmt.Printf("== Dumping data to GPX (%s)\n", gpxDir)
		mustCreateDirectory(gpxDir)
		mustDumpGPX(ctx, userID, jsonDir, gpxSports, cookies)
	}
}

func mustCreateDirectory(dumpDir string) {
	if err := os.MkdirAll(dumpDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
}
