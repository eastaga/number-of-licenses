package helpers

import (
	"encoding/csv"
	"errors"
	"log"
	"math"
	"os"
	"strings"
)

func getRecords() ([][]string, error) {
	file, err := os.Open(os.Getenv(CsvFileEnvVar))

	if err != nil {
		return nil, err
	}

	// Close the file
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	return records, err
}

func parseRecords(records [][]string) map[string]map[string]int {
	// map for storing computer id
	computerIDMap := make(map[string]struct{})
	// map for storing user id and type of machine mapping
	userMachineMap := make(map[string]map[string]int)

	for _, eachRecord := range records {
		// Don't process incomplete set
		if len(eachRecord) != 5 {
			continue
		}
		computerID := eachRecord[0]
		userID := eachRecord[1]
		appID := eachRecord[2]
		computerType := eachRecord[3]

		// Process only required appid
		if appID != os.Getenv(APPIDEnvVar) {
			log.Printf("skipping app id %s", appID)
			continue
		}

		// Don't process duplicate records
		if _, ok := computerIDMap[computerID]; ok {
			log.Printf("Computer ID %s already processed so skip it\n", computerID)
			continue
		}

		computerIDMap[computerID] = struct{}{}

		machineCount := make(map[string]int)

		if _, ok := userMachineMap[userID]; ok {
			machineCount = userMachineMap[userID]
		}

		computerTypeKey := strings.TrimSpace(computerType)
		machineCount[computerTypeKey]++
		userMachineMap[userID] = machineCount
	}
	return userMachineMap
}

func totalLicense(usersMachine map[string]map[string]int) int {
	totalLicenses := 0
	for _, machineTally := range usersMachine {
		// We need separate license for each desktop
		totalLicensePerUser := machineTally[DesktopKey]

		remainingLaptop := machineTally[LaptopKey] - totalLicensePerUser
		// Handle case when laptops are more than desktop
		if remainingLaptop > 0 {
			totalLicensePerUser += int(math.Floor(math.Round(float64(remainingLaptop) / 2.0)))
		}
		totalLicenses += totalLicensePerUser
	}

	return totalLicenses
}

func GetTotalLicenses() int {
	records, err := getRecords()
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	return totalLicense(parseRecords(records))
}

func Validate() error {
	_, appExists := os.LookupEnv(APPIDEnvVar)
	csvFile, csvExists := os.LookupEnv(CsvFileEnvVar)

	if !appExists {
		return errors.New("app ID is missing")
	}

	if !csvExists {
		return errors.New("CSV File is missing")
	}

	if !strings.HasSuffix(csvFile, ".csv") {
		return errors.New("file should have .csv extension")
	}

	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		return errors.New("file does not exists")
	}

	return nil
}
