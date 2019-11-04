package orm

func (o *ORM) LockingStrategy() LockingStrategy {
	return o.lockingStrategy
}
