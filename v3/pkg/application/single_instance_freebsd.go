//go:build freebsd && !server

package application

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// freebsdLock implements platformLock with an flock(2) lock file.
//
// Linux arbitrates the lock by claiming a D-Bus name, and the same bus carries
// the second instance's arguments across to the first. FreeBSD desktops do ship
// D-Bus, but there is no session bus a bare process can count on reaching, and
// the desktop portals the rest of the Linux backend leans on do not exist there.
// flock is in libc, needs no daemon, and is released by the kernel if the
// holder crashes, so it is the reliable half of the Linux design.
//
// The consequence is that the relay is one-way-less: a second launch still
// detects the running instance and exits, but its arguments are dropped and
// OnSecondInstanceLaunch does not fire. Carrying them would need a unix socket
// listening on a path the first instance advertises, which is a larger change
// than this milestone takes on.
type freebsdLock struct {
	file    *os.File
	path    string
	manager *singleInstanceManager
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &freebsdLock{manager: manager}, nil
}

// validUniqueID rejects anything that would let a UniqueID escape the lock
// directory. UniqueID is documented as a dot-separated identifier such as
// com.myapp.myapplication, and it reaches us as a path component here, so a
// value containing a separator must not be trusted to stay inside it.
func validUniqueID(uniqueID string) error {
	if uniqueID == "" {
		return errors.New("UniqueID is required for single instance lock")
	}
	if strings.ContainsAny(uniqueID, `/\`) || uniqueID == "." || uniqueID == ".." {
		return fmt.Errorf("UniqueID %q must not contain a path separator", uniqueID)
	}
	return nil
}

func (l *freebsdLock) acquire(uniqueID string) error {
	if err := validUniqueID(uniqueID); err != nil {
		return err
	}

	path := getLockPath(uniqueID)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("could not open the single instance lock file %q: %w", path, err)
	}

	if err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		file.Close()
		if errors.Is(err, unix.EWOULDBLOCK) || errors.Is(err, unix.EAGAIN) {
			// The kernel holds the lock, so the other instance is alive rather
			// than merely having crashed without cleaning up.
			return alreadyRunningError
		}
		return fmt.Errorf("could not lock %q: %w", path, err)
	}

	// Record the holder's pid so `fuser`/`lsof` point at something useful, and
	// so a stale file is recognisable as stale.
	if err := file.Truncate(0); err == nil {
		if _, err := file.WriteAt([]byte(strconv.Itoa(os.Getpid())), 0); err != nil {
			l.file, l.path = file, path
			return fmt.Errorf("could not record the pid in %q: %w", path, err)
		}
	}

	l.file = file
	l.path = path
	return nil
}

func (l *freebsdLock) release() {
	if l.file == nil {
		return
	}
	// The lock is dropped by closing the descriptor too, but unlocking
	// explicitly keeps the intent obvious next to the matching Flock above.
	_ = unix.Flock(int(l.file.Fd()), unix.LOCK_UN)
	l.file.Close()
	if l.path != "" {
		os.Remove(l.path)
	}
	l.file = nil
	l.path = ""
}

func (l *freebsdLock) notify(data string) error {
	// data is the serialized (and possibly encrypted) SecondInstanceData. There
	// is no channel to hand it over, and the caller exits immediately after
	// this returns regardless of the error, so report the loss once and let the
	// second instance stand down rather than leaving two windows running.
	if l.manager != nil && l.manager.app != nil {
		l.manager.app.warning(
			"single instance: a second instance exited without relaying its arguments, "+
				"because relaying them needs a session bus that FreeBSD does not provide here; "+
				"OnSecondInstanceLaunch will not fire (payload was %d bytes)", len(data))
	}
	return nil
}
