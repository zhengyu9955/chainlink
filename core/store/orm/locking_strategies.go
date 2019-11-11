package orm

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// NewLockingStrategy returns the locking strategy for a particular dialect
// to ensure exlusive access to the orm.
func NewLockingStrategy(dialect DialectName, dbpath string) (LockingStrategy, error) {
	switch dialect {
	case DialectPostgres:
		return NewPostgresLockingStrategy(dbpath)
	case DialectSqlite:
		return NewFileLockingStrategy(dbpath)
	}

	return nil, fmt.Errorf("unable to create locking strategy for dialect %s and path %s", dialect, dbpath)
}

// LockingStrategy employs the locking and unlocking of an underlying
// resource for exclusive access, usually a file or database.
type LockingStrategy interface {
	Open(timeout time.Duration) error
	Close(timeout time.Duration) error
	Lock(timeout time.Duration) error
	Unlock(timeout time.Duration) error
}

// FileLockingStrategy uses a file lock on disk to ensure exclusive access.
type FileLockingStrategy struct {
	path     string
	fileLock *flock.Flock
	m        *sync.Mutex
}

// NewFileLockingStrategy creates a new instance of FileLockingStrategy
// at the passed path.
func NewFileLockingStrategy(dbpath string) (LockingStrategy, error) {
	directory := filepath.Dir(dbpath)
	lockPath := filepath.Join(directory, "chainlink.lock")
	return &FileLockingStrategy{
		path:     lockPath,
		fileLock: flock.New(lockPath),
		m:        &sync.Mutex{},
	}, nil
}

func (s *FileLockingStrategy) Open(timeout time.Duration) error {
	return nil
}

func (s *FileLockingStrategy) Close(timeout time.Duration) error {
	return nil
}

// Lock returns immediately and assumes is always unlocked.
func (s *FileLockingStrategy) Lock(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	var err error
	locked := make(chan struct{})
	go func() {
		err = s.fileLock.Lock()
		close(locked)
	}()
	select {
	case <-locked:
	case <-normalizedTimeout(timeout):
		return fmt.Errorf("file locking strategy timed out for %s", s.path)
	}
	return err
}

func normalizedTimeout(timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return make(chan time.Time) // never time out
	}
	return time.After(timeout)
}

// Unlock is a noop.
func (s *FileLockingStrategy) Unlock(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()
	return s.fileLock.Unlock()
}

// PostgresLockingStrategy uses a postgres advisory lock to ensure exclusive
// access.
type PostgresLockingStrategy struct {
	db     *sql.DB
	conn   *sql.Conn
	path   string
	m      *sync.Mutex
	locked bool

	calledOpen  bool
	calledClose bool
	opens       [][]string
	closes      [][]string
}

// NewPostgresLockingStrategy returns a new instance of the PostgresLockingStrategy.
func NewPostgresLockingStrategy(path string) (LockingStrategy, error) {
	return &PostgresLockingStrategy{
		m:    &sync.Mutex{},
		path: path,
	}, nil
}

const postgresAdvisoryLockID int64 = 1027321974924625846

func (s *PostgresLockingStrategy) Open(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.calledOpen = true
	// s.opens = append(s.opens, getStack())

	ctx := context.Background()
	if timeout != 0 {
		// var cancel context.CancelFunc
		// ctx, cancel = context.WithTimeout(ctx, timeout)
		// defer cancel()
	}

	db, err := sql.Open(string(DialectPostgres), s.path)
	if err != nil {
		panic("failed on .Open")
		return err
	}
	s.db = db

	// `database/sql`.DB does opaque connection pooling, but PG advisory locks are per-connection
	conn, err := db.Conn(ctx)
	if err != nil {
		panic("failed on .Open2")
		return err
	}
	if conn == nil {
		panic("nil")
	}
	s.conn = conn

	return nil
}

func (s *PostgresLockingStrategy) Close(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.calledClose = true
	// s.closes = append(s.closes, getStack())
	//    callingTest := callingTest()
	// if _, ok := s.closes[callingTest]; ok {
	for _, f := range getStack() {
		fmt.Printf(".Close: %v\n", f)
	}
	fmt.Println("")
	// 	fmt.Printf("==================================================\n")
	// 	frmz := getStack()
	// 	for _, f := range frmz {
	// 		fmt.Printf("%v\n", f)
	// 	}
	// } else {
	// 	s.closes[callingTest] = getStack()
	// }

	// fmt.Printf("----- releasing all locks [%v] %v\n", s.name, s.db)

	// if s.db == nil {
	// 	return nil
	// }

	ctx := context.Background()
	if timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, err := s.conn.ExecContext(ctx, "SELECT pg_advisory_unlock_all()")
	if err != nil {
		panic("lock failed on .Close")
		return errors.Wrapf(ErrReleaseLockFailed, "postgres advisory locking strategy failed on .Close, timeout set to %v: %v", displayTimeout(timeout), err)
	}

	connErr := s.conn.Close()
	if connErr == sql.ErrConnDone {
		connErr = nil
	}
	dbErr := s.db.Close()
	if dbErr == sql.ErrConnDone {
		dbErr = nil
	}
	s.db = nil
	s.conn = nil
	return multierr.Combine(
		connErr,
		dbErr,
	)
}

// Lock uses a blocking postgres advisory lock that times out at the passed
// timeout.
func (s *PostgresLockingStrategy) Lock(timeout time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.conn == nil {
		// callingTest := callingTest()
		fmt.Printf("DB ~> %v\n", s.path)
		// fmt.Printf("CALLING TEST ~> %v\n", callingTest)
		fmt.Printf("CALLING TEST called open? ~> %v\n", s.calledOpen)
		fmt.Printf("CALLING TEST called close? ~> %v\n", s.calledClose)
		// for _, open := range s.opens {
		// 	for _, f := range open {
		// 		fmt.Printf("%v\n", f)
		// 	}
		// 	fmt.Printf("--------------------------------------------------\n")
		// }
		// fmt.Printf("==================================================\n")
		// for _, cls := range s.closes {
		// 	for _, f := range cls {
		// 		fmt.Printf("%v\n", f)
		// 	}
		// 	fmt.Printf("--------------------------------------------------\n")
		// }

		// fmt.Printf("**************************************************\n")

		// fs := getStack()
		// for _, f := range fs {
		// 	fmt.Printf("%v\n", f)
		// }

		panic("you must call PostgresLockingStrategy#Open first")
	}

	ctx := context.Background()
	if timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, err := s.conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", postgresAdvisoryLockID)
	if err != nil {
		return errors.Wrapf(ErrNoAdvisoryLock, "postgres advisory locking strategy failed on .Lock, timeout set to %v: %v", displayTimeout(timeout), err)
	}
	return nil
}

// Unlock unlocks the locked postgres advisory lock.
func (s *PostgresLockingStrategy) Unlock(timeout time.Duration) error {
	return nil
}

func logTest(thisFunc string) {
	fmt.Printf("[%v] %v\n", callingTest(), thisFunc)
}

func callingTest() string {
	// Ask runtime.Callers for up to 10 pcs, including runtime.Callers itself.
	pc := make([]uintptr, 30)
	n := runtime.Callers(0, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return ""
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()
		// To keep this example's output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		if strings.Contains(frame.File, "_test.go") {
			return frame.Function
		}
		if !more {
			break
		}
	}
	return "???"
}

func getStack() []string {
	// Ask runtime.Callers for up to 10 pcs, including runtime.Callers itself.
	pc := make([]uintptr, 30)
	n := runtime.Callers(0, pc)
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return nil
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frmz := runtime.CallersFrames(pc)
	var frames []string

	// Loop to get frmz.
	// A fixed number of pcs can expand to an indefinite number of Frames.
	for {
		frame, more := frmz.Next()
		frames = append(frames, fmt.Sprintf("%v", frame))
		if !more {
			break
		}
	}
	return frames
}
