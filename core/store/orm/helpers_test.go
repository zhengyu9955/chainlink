package orm

import (
	"time"
)

func (o *ORM) LockingStrategyHelperSimulateDisconnect() (error, error) {
	err1 := o.lockingStrategy.(*PostgresLockingStrategy).conn.Close()
	err2 := o.lockingStrategy.(*PostgresLockingStrategy).db.Close()
	return err1, err2
}

func (o *ORM) LockingStrategyHelperSimulateReconnect(timeout time.Duration) error {
	return o.lockingStrategy.Open(timeout)
}
