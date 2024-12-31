package gapp

import "fmt"

type Sse struct {
	c       *Content
	isClose bool
	key     string
}

func (s *Sse) Send(data string, event ...string) error {
	if len(event) == 1 {
		_, e := s.c.w.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", data, event[0])))
		if e == nil {
			s.c.flush()
		} else {
			s.isClose = true
		}
		return e
	} else {
		_, e := s.c.w.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
		if e == nil {
			s.c.flush()
		} else {
			s.isClose = true
		}
		return e
	}
}
func (s *Sse) IsClose() bool {
	return s.isClose
}
func (s *Sse) GetKey() string {
	return s.key
}
func (s *Sse) GetUserId() uint64 {
	return s.c.GetUserID()
}
