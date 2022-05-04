package persistance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestCouchDB(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// setup suite
func (suite *TestSuite) SetupSuite() {

}

func (suite *TestSuite) BeforeTest(suiteName, testName string) {
	// delete db if exist
	_, err := os.Stat(testName + ".bolt")
	if err == nil {
		os.Remove(testName + ".bolt")
	}
	Open(testName + ".bolt")
}

func (suite *TestSuite) AfterTest(suiteName, testName string) {
	// delete db if exist
	_, err := os.Stat(testName + ".bolt")
	if err == nil {
		os.Remove(testName + ".bolt")
	}
}

func (suite *TestSuite) TearDownSuite() {

}
