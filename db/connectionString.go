package db

type ConnnetionString string

func (cs *ConnnetionString) String() string {
	return string(*cs)
}
