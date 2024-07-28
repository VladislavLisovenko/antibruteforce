package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const delay = 5 * time.Second

func TestFeatures(t *testing.T) {
	fmt.Printf("wait %s for service availability...", delay)
	time.Sleep(delay)

	suite := godog.TestSuite{
		Name:                "integration",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:    "progress",
			Paths:     []string{"features"},
			Randomize: 0,
		},
	}

	if st := suite.Run(); st != 0 {
		t.Fatal("failed to run feature tests")
	}
}

func InitializeScenario(s *godog.ScenarioContext) {
	test := new(apiTest)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	s.Step(`^The response should match text:$`, test.theResponseShouldMatchText)
}
