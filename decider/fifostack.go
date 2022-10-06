package decider

type element struct {
	next  *element
	value *Configuration
}

type fifoStack struct {
	first  *element
	last   *element
	length int
}

func (s *fifoStack) push(value *Configuration) {
	newElement := &element{nil, value}
	if s.length == 0 {
		s.first = newElement
		s.last = newElement

	} else {
		s.last.next = newElement
		s.last = newElement
	}
	s.length++
}
func (s *fifoStack) pop() *Configuration {
	poppedElement := s.first
	s.first = poppedElement.next
	s.length--
	if s.length <= 0 {
		s.length = 0
		s.last = nil
	}
	return poppedElement.value
}
