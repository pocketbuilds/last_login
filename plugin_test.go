package last_login

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

const testDataDir = "./test/pb_data/"

func TestPlugin(t *testing.T) {
	setupTestApp := func(t testing.TB) *tests.TestApp {
		testApp, err := tests.NewTestApp(testDataDir)
		if err != nil {
			t.Fatal(err)
		}
		(&Plugin{
			// test config will go here
			FieldName: "last_login",
		}).Init(testApp)
		return testApp
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "create record",
			Method:         http.MethodPost,
			URL:            "/api/collections/users/auth-with-password",
			TestAppFactory: setupTestApp,
			Body: strings.NewReader(`{
				"identity": "test@example.com",
				"password": "test12345"
			}`),
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"last_login":`,
			},
			NotExpectedContent: []string{
				`"last_login":""`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordUpdate": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
