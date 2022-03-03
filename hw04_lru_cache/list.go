package hw04lrucache

import (
	"errors"
	"fmt"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem) error
	MoveToFront(i *ListItem) error
}

var ErrorEmptyList = errors.New("empty list")

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}
	if l.head != nil {
		l.head.Prev = node
	}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
	l.len++
	return l.head
}

func (l *list) PopFront() (interface{}, error) {
	if l.head == nil {
		return nil, fmt.Errorf("PopFront: %w", ErrorEmptyList)
	}

	tmp := l.head
	l.head = l.head.Next
	if l.head != nil {
		l.head.Prev = nil
	}
	if tmp == l.tail {
		l.tail = nil
	}
	tmp.Next = nil
	l.len--
	return tmp.Value, nil
}

func (l *list) PushBack(v interface{}) *ListItem {
	node := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}
	if l.tail != nil {
		l.tail.Next = node
	}
	l.tail = node
	if l.head == nil {
		l.head = node
	}

	l.len++
	return node
}

func (l *list) PopBack() (interface{}, error) {
	if l.tail == nil {
		return nil, fmt.Errorf("PopBack: %w", ErrorEmptyList)
	}
	tmp := l.tail
	l.tail = l.tail.Prev
	if l.tail != nil {
		l.tail.Next = nil
	}
	if tmp == l.head {
		l.head = nil
	}
	tmp.Prev = nil
	l.len--

	return tmp.Value, nil
}

func (l *list) Remove(i *ListItem) error {
	if l.len == 0 {
		return fmt.Errorf("remove: %w, val: %v", ErrorEmptyList, i)
	}

	if l.head == i {
		if _, err := l.PopFront(); err != nil {
			return fmt.Errorf("remove: %w, val: %v", err, i)
		}
		return nil
	}

	if l.tail == i {
		if _, err := l.PopBack(); err != nil {
			return fmt.Errorf("remove: %w, val: %v", err, i)
		}
		return nil
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Next = nil
	i.Prev = nil
	l.len--
	return nil
}

func (l *list) MoveToFront(i *ListItem) error {
	if l.len == 0 {
		return ErrorEmptyList
	}

	if i == l.head || l.len == 1 {
		return nil
	}
	if err := l.Remove(i); err != nil {
		return fmt.Errorf("MoveToFront: %w", err)
	}

	i.Next = l.head
	l.head.Prev = i
	l.head = i

	return nil
}

func (l *list) String() string {
	if l.len == 0 {
		return ""
	}
	var (
		currItem = l.head
		res      = make([]interface{}, 0, l.len)
	)

	for {
		if currItem == nil {
			break
		}
		res = append(res, currItem.Value)
		currItem = currItem.Next
	}
	return fmt.Sprint(res)
}

func NewList() List {
	return new(list)
}
