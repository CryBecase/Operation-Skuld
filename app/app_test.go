package app

import (
	"context"
	"testing"

	"skuld/env"
	"skuld/xlogger"
)

type mockServer struct{}

func (m mockServer) Run() error { return nil }

func (m mockServer) Close(ctx context.Context) error { return nil }

type mockLogger struct{}

func (m mockLogger) Debug(msg string, kvs ...interface{}) {}

func (m mockLogger) Info(msg string, kvs ...interface{}) {}

func (m mockLogger) Warn(msg string, kvs ...interface{}) {}

func (m mockLogger) Error(msg string, kvs ...interface{}) {}

func (m mockLogger) Fatal(msg string, kvs ...interface{}) {}

func TestNew(t *testing.T) {
	type testCase struct {
		name        string
		env         env.Env
		logger      xlogger.Logger
		svr         Server
		expectPanic bool
	}

	logger := &mockLogger{}
	svr := &mockServer{}

	testTable := []testCase{
		{
			name:        "invalid env",
			env:         "bad",
			logger:      logger,
			svr:         svr,
			expectPanic: true,
		},
		{
			name:        "nil logger",
			env:         env.New("local"),
			logger:      nil,
			svr:         svr,
			expectPanic: true,
		},
		{
			name:        "nil svr",
			env:         env.New("local"),
			logger:      logger,
			svr:         nil,
			expectPanic: true,
		},
		{
			name:        "success",
			env:         env.New("local"),
			logger:      logger,
			svr:         svr,
			expectPanic: false,
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if v.expectPanic && r == nil {
					t.Fatalf("expect panic, but no panic")
				}
				if !v.expectPanic && r != nil {
					t.Fatalf("expect no panic, but panic")
				}
			}()
			_ = New(env.NewInfo("", "", 0, v.env), v.logger, v.svr)
		})
	}
}
