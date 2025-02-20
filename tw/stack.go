package tw

type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) Get(index int) T {
	return s.values[index]
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop() T {
	value := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return value
}

func (s *Stack[T]) Peek() T {
	return s.values[len(s.values)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.values)
}
