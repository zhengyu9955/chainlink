package web_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	"chainlink/core/store/models"
	"chainlink/core/store/presenters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var endAt = time.Now().AddDate(0, 10, 0).Round(time.Second).UTC()
var endAtISO8601 = endAt.Format(time.RFC3339)

func TestServiceAgreementsController_Create(t *testing.T) {
	t.Parallel()

	base := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	base = strings.Replace(base, "2019-10-19T22:17:19Z", endAtISO8601, 1)
	tests := []struct {
		name                string
		input               string
		wantLogSubscription bool
		wantCode            int
	}{
		{"success", base, true, http.StatusOK},
		{"fails validation", cltest.MustJSONDel(t, base, "payment"), false, http.StatusUnprocessableEntity},
		{"invalid JSON", "{", false, http.StatusUnprocessableEntity},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
			defer cleanup()

			if test.wantLogSubscription {
				app.EthMock.RegisterSubscription("logs")
			}

			require.NoError(t, app.Start())

			client := app.NewHTTPClient()

			resp, cleanup2 := client.Post("/v2/service_agreements", bytes.NewBufferString(test.input))
			defer cleanup2()

			cltest.AssertServerResponse(t, resp, test.wantCode)
			if test.wantCode == http.StatusOK {
				responseSA := models.ServiceAgreement{}

				err := cltest.ParseJSONAPIResponse(t, resp, &responseSA)
				require.NoError(t, err)
				assert.NotEqual(t, "", responseSA.ID)
				assert.NotEqual(t, "", responseSA.Signature.String())

				createdSA := cltest.FindServiceAgreement(t, app.Store, responseSA.ID)
				assert.NotEqual(t, "", createdSA.ID)
				assert.NotEqual(t, "", createdSA.Signature.String())
				assert.Equal(t, endAt, createdSA.Encumbrance.EndAt.Time)

				var jobids []*models.ID
				for _, j := range app.JobSubscriber.Jobs() {
					jobids = append(jobids, j.ID)
				}
				assert.Contains(t, jobids, createdSA.JobSpec.ID)
				app.EthMock.EventuallyAllCalled(t)
			}
		})
	}
}

func TestServiceAgreementsController_Create_isIdempotent(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	app.EthMock.RegisterSubscription("logs")

	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	base := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
	base = strings.Replace(base, "2019-10-19T22:17:19Z", endAtISO8601, 1)
	reader := bytes.NewBuffer([]byte(base))

	resp, cleanup := client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	response1 := models.ServiceAgreement{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &response1))

	reader = bytes.NewBuffer([]byte(base))
	resp, cleanup = client.Post("/v2/service_agreements", reader)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)
	response2 := models.ServiceAgreement{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &response2))

	assert.Equal(t, response1.ID, response2.ID)
	assert.Equal(t, response1.JobSpec.ID, response2.JobSpec.ID)
	app.EthMock.EventuallyAllCalled(t)
}

func TestServiceAgreementsController_Show(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()

	input := cltest.MustReadFile(t, "testdata/hello_world_agreement.json")
	sa, err := cltest.ServiceAgreementFromString(string(input))
	require.NoError(t, err)
	require.NoError(t, app.Store.CreateServiceAgreement(&sa))

	resp, cleanup := client.Get("/v2/service_agreements/" + sa.ID)
	defer cleanup()
	cltest.AssertServerResponse(t, resp, http.StatusOK)

	normalizedInput := cltest.NormalizedJSON(t, input)
	parsed := presenters.ServiceAgreement{}
	cltest.ParseJSONAPIResponse(t, resp, &parsed)
	assert.Equal(t, normalizedInput, parsed.RequestBody)
}
