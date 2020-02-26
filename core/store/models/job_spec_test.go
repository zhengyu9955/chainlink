package models_test

import (
	"testing"
	"time"

	"chainlink/core/adapters"
	"chainlink/core/assets"
	"chainlink/core/internal/cltest"
	"chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestNewInitiatorFromRequest(t *testing.T) {
	t.Parallel()

	job := cltest.NewJob()
	tests := []struct {
		name     string
		initrReq models.InitiatorRequest
		jobSpec  models.JobSpec
		want     models.Initiator
	}{
		{
			name: models.InitiatorWeb,
			initrReq: models.InitiatorRequest{
				Type: models.InitiatorWeb,
			},
			jobSpec: job,
			want: models.Initiator{
				Type:      models.InitiatorWeb,
				JobSpecID: job.ID,
			},
		},
		{
			name: models.InitiatorWeb,
			initrReq: models.InitiatorRequest{
				Type: models.InitiatorFluxMonitor,
				InitiatorParams: models.InitiatorParams{
					IdleThreshold: models.Duration(5 * time.Second),
					Precision:     2,
					Threshold:     5,
				},
			},
			jobSpec: job,
			want: models.Initiator{
				Type:      models.InitiatorFluxMonitor,
				JobSpecID: job.ID,
				InitiatorParams: models.InitiatorParams{
					IdleThreshold:   models.Duration(5 * time.Second),
					PollingInterval: models.FluxMonitorDefaultInitiatorParams.PollingInterval,
					Precision:       2,
					Threshold:       5,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			res := models.NewInitiatorFromRequest(
				test.initrReq,
				test.jobSpec,
			)
			assert.Equal(t, test.want, res)
		})
	}
}

func TestNewJobFromRequest(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * 7")
	require.NoError(t, store.CreateJob(&j1))

	jsr := models.JobSpecRequest{
		Initiators: cltest.BuildInitiatorRequests(t, j1.Initiators),
		Tasks:      cltest.BuildTaskRequests(t, j1.Tasks),
		StartAt:    j1.StartAt,
		EndAt:      j1.EndAt,
		MinPayment: assets.NewLink(5),
	}

	j2 := models.NewJobFromRequest(jsr)
	require.NoError(t, store.CreateJob(&j2))

	fetched1, err := store.FindJob(j1.ID)
	assert.NoError(t, err)
	assert.Len(t, fetched1.Initiators, 1)
	assert.Len(t, fetched1.Tasks, 1)
	assert.Nil(t, fetched1.MinPayment)

	fetched2, err := store.FindJob(j2.ID)
	assert.NoError(t, err)
	assert.Len(t, fetched2.Initiators, 1)
	assert.Len(t, fetched2.Tasks, 1)
	assert.Equal(t, assets.NewLink(5), fetched2.MinPayment)
}

func TestJobSpec_Save(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	befCreation := time.Now()
	j1 := cltest.NewJobWithSchedule("* * * * 7")
	aftCreation := time.Now()
	assert.True(t, true, j1.CreatedAt.After(aftCreation), j1.CreatedAt.Before(befCreation))
	assert.False(t, false, j1.CreatedAt.IsZero())

	befInsertion := time.Now()
	assert.NoError(t, store.CreateJob(&j1))
	aftInsertion := time.Now()
	assert.True(t, true, j1.CreatedAt.After(aftInsertion), j1.CreatedAt.Before(befInsertion))

	initr := j1.Initiators[0]

	j2, err := store.FindJob(j1.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	assert.Equal(t, initr.Schedule, j2.Initiators[0].Schedule)
}

func TestJobSpec_NewRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithSchedule("1 * * * *")
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp", `{"a":1}`)}

	run := cltest.NewJobRun(job)

	assert.Equal(t, job.ID, run.JobSpecID)
	assert.Equal(t, 1, len(run.TaskRuns))

	taskRun := run.TaskRuns[0]
	assert.Equal(t, "noop", taskRun.TaskSpec.Type.String())
	adapter, _ := adapters.For(taskRun.TaskSpec, store.Config, store.ORM)
	assert.NotNil(t, adapter)
	assert.JSONEq(t, `{"type":"NoOp","a":1}`, taskRun.TaskSpec.Params.String())

	assert.Equal(t, job.Initiators[0], run.Initiator)
}

func TestJobEnded(t *testing.T) {
	t.Parallel()

	endAt := cltest.ParseNullableTime(t, "3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		endAt   null.Time
		current time.Time
		want    bool
	}{
		{"no end at", null.Time{Valid: false}, endAt.Time, false},
		{"before end at", endAt, endAt.Time.Add(-time.Nanosecond), false},
		{"at end at", endAt, endAt.Time, false},
		{"after end at", endAt, endAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.EndAt = test.endAt

			assert.Equal(t, test.want, job.Ended(test.current))
		})
	}
}

func TestJobSpec_Started(t *testing.T) {
	t.Parallel()

	startAt := cltest.ParseNullableTime(t, "3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		startAt null.Time
		current time.Time
		want    bool
	}{
		{"no start at", null.Time{Valid: false}, startAt.Time, true},
		{"before start at", startAt, startAt.Time.Add(-time.Nanosecond), false},
		{"at start at", startAt, startAt.Time, true},
		{"after start at", startAt, startAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.StartAt = test.startAt

			assert.Equal(t, test.want, job.Started(test.current))
		})
	}
}

func TestNewTaskType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		errored bool
	}{
		{"basic", "NoOp", "noop", false},
		{"special characters", "-_-", "-_-", false},
		{"invalid character", "NoOp!", "", true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := models.NewTaskType(test.input)

			if test.errored {
				assert.Error(t, err)
			} else {
				assert.Equal(t, models.TaskType(test.want), got)
				assert.NoError(t, err)
			}
		})
	}
}
