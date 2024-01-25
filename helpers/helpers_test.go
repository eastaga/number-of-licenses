package helpers

import (
	"errors"
	"github.com/onsi/gomega"
	"os"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc          string
		appID         string
		csvFile       string
		fileExists    bool
		expectedError error
	}{
		{
			desc:          "app id is not set",
			expectedError: errors.New("app ID is missing"),
		},
		{
			desc:          "csv file is not set",
			appID:         DummyAppID,
			expectedError: errors.New("CSV File is missing"),
		},
		{
			desc:          "type of file is not csv",
			appID:         DummyAppID,
			csvFile:       "test.txt",
			expectedError: errors.New("file should have .csv extension"),
		},
		{
			desc:          "csv file does not exist",
			appID:         DummyAppID,
			csvFile:       "test.csv",
			expectedError: errors.New("file does not exists"),
		},
		{
			desc:          "no error",
			appID:         DummyAppID,
			csvFile:       "test.csv",
			fileExists:    true,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			g := gomega.NewWithT(t)

			// clean up first
			os.Unsetenv(CsvFileEnvVar)
			os.Unsetenv(APPIDEnvVar)
			// test setup
			if test.appID != "" {
				os.Setenv(APPIDEnvVar, test.appID)
			}

			if test.csvFile != "" {
				os.Setenv(CsvFileEnvVar, test.csvFile)
			}

			if test.fileExists {
				os.Create(test.csvFile)
			}

			// test the function
			err := Validate()
			if test.expectedError != nil {
				g.Expect(err).NotTo(gomega.BeNil())
				g.Expect(err.Error()).To(gomega.Equal(test.expectedError.Error()))
				return
			}
			g.Expect(err).To(gomega.BeNil())

			// cleanup
			defer func() {
				os.Unsetenv(CsvFileEnvVar)
				os.Unsetenv(APPIDEnvVar)
				if test.fileExists {
					if err := os.Remove(test.csvFile); err != nil {
						t.Logf("unable to delete file")
					}
				}
			}()
		})
	}
}

func TestTotalLicense(t *testing.T) {
	tests := []struct {
		desc                  string
		laptopCount           int
		desktopCount          int
		expectedTotalLicenses int
	}{
		{
			desc:                  "desktop and laptops are equal",
			laptopCount:           3,
			desktopCount:          3,
			expectedTotalLicenses: 3,
		},
		{
			desc:                  "desktop are more than laptop",
			laptopCount:           3,
			desktopCount:          4,
			expectedTotalLicenses: 4,
		},
		{
			desc:                  "laptops are more than desktop",
			laptopCount:           4,
			desktopCount:          3,
			expectedTotalLicenses: 4,
		},
		{
			desc:                  "laptops(odd number remaining) are more than desktop",
			laptopCount:           6,
			desktopCount:          3,
			expectedTotalLicenses: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			g := gomega.NewWithT(t)
			// create user data
			testMap := make(map[string]map[string]int)

			testMachineMap := make(map[string]int)
			testMachineMap[DesktopKey] = test.desktopCount
			testMachineMap[LaptopKey] = test.laptopCount
			testMap[DummyUserID] = testMachineMap

			totalLicenses := totalLicense(testMap)
			g.Expect(totalLicenses).To(gomega.Equal(test.expectedTotalLicenses))
		})
	}
}

func getMap() map[string]map[string]int {
	mp := make(map[string]int)
	mp[DesktopKey] = 2
	mp[LaptopKey] = 1

	testMp := make(map[string]map[string]int)
	testMp[DummyUserID] = mp

	return testMp
}
func TestParseRecords(t *testing.T) {
	_ = os.Setenv(APPIDEnvVar, DummyAppID)
	comment := "test comment"
	testEntry1 := []string{"1", DummyUserID, DummyAppID, LaptopKey, comment}
	testEntry2 := []string{"2", DummyUserID, DummyAppID, "DESKTOP", comment}
	testEntry3 := []string{"3", DummyUserID, DummyAppID, DesktopKey, comment}
	unwantedApp := []string{"3", DummyUserID, "123", LaptopKey, comment}
	incompleteEntry := []string{"4", DummyUserID, DummyAppID, LaptopKey}

	tests := []struct {
		desc            string
		records         [][]string
		expectedMapping map[string]map[string]int
	}{
		{
			desc:            "incomplete record",
			records:         [][]string{testEntry1, testEntry2, testEntry3, incompleteEntry},
			expectedMapping: getMap(),
		},
		{
			desc:            "records contain unwanted app id",
			records:         [][]string{testEntry1, testEntry2, testEntry3, unwantedApp},
			expectedMapping: getMap(),
		},
		{
			desc:            "records contain duplicate computer id",
			records:         [][]string{testEntry1, testEntry2, testEntry3, testEntry3},
			expectedMapping: getMap(),
		},
		{
			desc:            "valid record",
			records:         [][]string{testEntry1, testEntry2, testEntry3},
			expectedMapping: getMap(),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			g := gomega.NewWithT(t)
			actualMap := parseRecords(test.records)
			g.Expect(reflect.DeepEqual(actualMap, test.expectedMapping)).To(gomega.BeTrue())
		})
	}

}
