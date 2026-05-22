//ff:func feature=session type=helper control=sequence
//ff:what Creates an empty Session instance
package session

func New() *Session {
	return &Session{}
}
