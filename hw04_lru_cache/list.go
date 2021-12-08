package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length      int
	front, back *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: l.front, Prev: nil}
	if l.front != nil {
		l.front.Prev = newItem
	}
	if l.back == nil {
		l.back = newItem
	}
	l.front = newItem
	l.length++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: nil, Prev: l.back}
	if l.back != nil {
		l.back.Next = newItem
	}
	if l.front == nil {
		l.front = newItem
	}
	l.back = newItem
	l.length++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.front == i {
		l.front = i.Next
	}
	if l.back == i {
		l.back = i.Prev
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i != l.front {
		l.Remove(i)
		i.Next = l.front
		l.front.Prev = i
		l.front = i
		l.length++
	}
}

func NewList() List {
	return new(list)
}
