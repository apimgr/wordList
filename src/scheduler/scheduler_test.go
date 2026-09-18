package scheduler

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestAddRemoveEnableDisableTask(t *testing.T) {
	s := New()
	s.AddTask("t1", time.Hour, func() error { return nil })

	tasks := s.GetTasks()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if !tasks[0].Enabled {
		t.Error("expected task to be enabled by default")
	}

	s.DisableTask("t1")
	tasks = s.GetTasks()
	if tasks[0].Enabled {
		t.Error("expected task to be disabled")
	}

	s.EnableTask("t1")
	tasks = s.GetTasks()
	if !tasks[0].Enabled {
		t.Error("expected task to be re-enabled")
	}

	// no-op on unknown task
	s.EnableTask("unknown")
	s.DisableTask("unknown")

	s.RemoveTask("t1")
	if len(s.GetTasks()) != 0 {
		t.Error("expected task to be removed")
	}
}

func TestRunNow(t *testing.T) {
	s := New()
	var ran int32
	s.AddTask("t1", time.Hour, func() error {
		atomic.AddInt32(&ran, 1)
		return nil
	})

	if err := s.RunNow("t1"); err != nil {
		t.Fatalf("RunNow() error = %v", err)
	}
	if atomic.LoadInt32(&ran) != 1 {
		t.Errorf("task ran %d times, want 1", ran)
	}

	// unknown task is a no-op, returns nil
	if err := s.RunNow("unknown"); err != nil {
		t.Errorf("RunNow(unknown) error = %v, want nil", err)
	}

	// error propagates
	wantErr := errors.New("boom")
	s.AddTask("t2", time.Hour, func() error { return wantErr })
	if err := s.RunNow("t2"); !errors.Is(err, wantErr) {
		t.Errorf("RunNow() error = %v, want %v", err, wantErr)
	}
}

func TestStartStop(t *testing.T) {
	s := New()
	s.Start()
	// starting twice should be a no-op, not panic
	s.Start()
	s.Stop()
	// stopping twice should be a no-op, not panic
	s.Stop()
}

func TestRunDueTasks(t *testing.T) {
	s := New()
	done := make(chan struct{}, 1)
	s.AddTask("due", time.Hour, func() error {
		done <- struct{}{}
		return nil
	})

	// force the task to be immediately due
	s.mu.Lock()
	s.tasks["due"].NextRun = time.Now().Add(-time.Minute)
	s.mu.Unlock()

	s.runDueTasks()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("expected due task to run")
	}
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		in   string
		want time.Duration
	}{
		{"minutely", time.Minute},
		{"hourly", time.Hour},
		{"daily", 24 * time.Hour},
		{"weekly", 7 * 24 * time.Hour},
		{"monthly", 30 * 24 * time.Hour},
		{"2h", 2 * time.Hour},
		{"not-a-duration", 24 * time.Hour},
	}

	for _, tt := range tests {
		if got := ParseInterval(tt.in); got != tt.want {
			t.Errorf("ParseInterval(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
