// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package mock

type MockConn struct {
	ReadMessageCalled chan bool
	ReadMessageOutput struct {
		Ret0 chan int
		Ret1 chan []byte
		Ret2 chan error
	}
	WriteMessageCalled chan bool
	WriteMessageInput  struct {
		Arg0 chan int
		Arg1 chan []byte
	}
	WriteMessageOutput struct {
		Ret0 chan error
	}
	CloseCalled chan bool
	CloseOutput struct {
		Ret0 chan error
	}
}

func NewMockConn() *MockConn {
	m := &MockConn{}
	m.ReadMessageCalled = make(chan bool, 100)
	m.ReadMessageOutput.Ret0 = make(chan int, 100)
	m.ReadMessageOutput.Ret1 = make(chan []byte, 100)
	m.ReadMessageOutput.Ret2 = make(chan error, 100)
	m.WriteMessageCalled = make(chan bool, 100)
	m.WriteMessageInput.Arg0 = make(chan int, 100)
	m.WriteMessageInput.Arg1 = make(chan []byte, 100)
	m.WriteMessageOutput.Ret0 = make(chan error, 100)
	m.CloseCalled = make(chan bool, 100)
	m.CloseOutput.Ret0 = make(chan error, 100)
	return m
}
func (m *MockConn) ReadMessage() (int, []byte, error) {
	m.ReadMessageCalled <- true
	return <-m.ReadMessageOutput.Ret0, <-m.ReadMessageOutput.Ret1, <-m.ReadMessageOutput.Ret2
}
func (m *MockConn) WriteMessage(arg0 int, arg1 []byte) error {
	m.WriteMessageCalled <- true
	m.WriteMessageInput.Arg0 <- arg0
	m.WriteMessageInput.Arg1 <- arg1
	return <-m.WriteMessageOutput.Ret0
}
func (m *MockConn) Close() error {
	m.CloseCalled <- true
	return <-m.CloseOutput.Ret0
}
