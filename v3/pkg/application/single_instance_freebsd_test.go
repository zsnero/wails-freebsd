//go:build freebsd && !server

package application

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestValidUniqueID(t *testing.T) {
	tests := []struct {
		name     string
		uniqueID string
		wantErr  bool
	}{
		{name: "reverse dns id", uniqueID: "com.myapp.myapplication"},
		{name: "hyphens and digits", uniqueID: "net.my-company.my-app2"},
		{name: "empty", uniqueID: "", wantErr: true},
		{name: "forward slash", uniqueID: "com.myapp/../etc", wantErr: true},
		{name: "backslash", uniqueID: `com.myapp\evil`, wantErr: true},
		{name: "bare dot", uniqueID: ".", wantErr: true},
		{name: "bare dotdot", uniqueID: "..", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validUniqueID(tt.uniqueID)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("validUniqueID(%q) = nil, want error", tt.uniqueID)
				}
				return
			}
			if err != nil {
				t.Fatalf("validUniqueID(%q) = %v, want nil", tt.uniqueID, err)
			}
		})
	}
}

// TestFreeBSDLockExcludesSecondInstance checks that the flock actually
// arbitrates: the second acquire must fail with alreadyRunningError so the
// caller hands off and exits, and must not report a generic error that the
// caller would treat as a fatal startup failure.
func TestFreeBSDLockExcludesSecondInstance(t *testing.T) {
	uniqueID := "com.wails.test.singleinstance"

	first, err := newPlatformLock(&singleInstanceManager{})
	if err != nil {
		t.Fatalf("newPlatformLock: %v", err)
	}
	if err := first.acquire(uniqueID); err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	t.Cleanup(first.release)

	// Confirm the lock file really was created where getLockPath says, and that
	// it records this pid, so the "is it a live lock or a stale file" question
	// has an answer on disk.
	path := getLockPath(uniqueID)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected a lock file at %q: %v", path, err)
	}
	if got, err := os.ReadFile(path); err == nil {
		if !strings.Contains(string(got), strconv.Itoa(os.Getpid())) {
			t.Errorf("lock file %q contains %q, want the current pid %d", path, got, os.Getpid())
		}
	}

	second, err := newPlatformLock(&singleInstanceManager{})
	if err != nil {
		t.Fatalf("newPlatformLock (second): %v", err)
	}
	err = second.acquire(uniqueID)
	if !errors.Is(err, alreadyRunningError) {
		t.Fatalf("second acquire = %v, want alreadyRunningError", err)
	}
	second.release()
}

// TestFreeBSDLockReleaseAllowsReacquire guards the release path: if the lock
// file were left locked, the next launch of the same app would see itself as a
// second instance forever.
func TestFreeBSDLockReleaseAllowsReacquire(t *testing.T) {
	uniqueID := "com.wails.test.reacquire"

	first, err := newPlatformLock(&singleInstanceManager{})
	if err != nil {
		t.Fatalf("newPlatformLock: %v", err)
	}
	if err := first.acquire(uniqueID); err != nil {
		t.Fatalf("first acquire: %v", err)
	}
	first.release()

	if _, err := os.Stat(getLockPath(uniqueID)); !os.IsNotExist(err) {
		t.Errorf("release left %q behind (stat err = %v), want it removed", getLockPath(uniqueID), err)
	}

	second, err := newPlatformLock(&singleInstanceManager{})
	if err != nil {
		t.Fatalf("newPlatformLock (second): %v", err)
	}
	defer second.release()
	if err := second.acquire(uniqueID); err != nil {
		t.Fatalf("acquire after release = %v, want nil", err)
	}
}

// TestFreeBSDLockRejectsPathTraversal guards the rejection against a
// regression that would let a UniqueID place the lock file outside the
// temporary directory it belongs in.
func TestFreeBSDLockRejectsPathTraversal(t *testing.T) {
	for _, uniqueID := range []string{"../escape", "sub/dir", `..\escape`} {
		lock, err := newPlatformLock(&singleInstanceManager{})
		if err != nil {
			t.Fatalf("newPlatformLock: %v", err)
		}
		if err := lock.acquire(uniqueID); err == nil {
			lock.release()
			t.Errorf("acquire(%q) = nil, want an error", uniqueID)
		}
	}
}
