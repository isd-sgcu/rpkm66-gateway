package qr

type IContext interface {
	JSON(int, interface{})
	UserID() string
	Bind(interface{}) error
	ID() (string, error)
}
