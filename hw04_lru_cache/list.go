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
	// List Remove me after realization.
	// Place your code here.
	firstNode  *ListItem
	lastNode   *ListItem
	lenghtList int
}

func (l *list) Len() int {
	return l.lenghtList
}

func (l *list) Front() *ListItem {
	return l.firstNode
}

func (l *list) Back() *ListItem {
	return l.lastNode
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := new(ListItem)
	newListItem.Value = v
	if l.Len() == 0 {
		l.firstNode = newListItem
		l.lastNode = newListItem
	} else {
		newListItem.Next = l.firstNode
		l.firstNode.Prev = newListItem
		l.firstNode = newListItem
	}
	l.lenghtList++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := new(ListItem)
	newListItem.Value = v
	if l.Len() == 0 {
		l.firstNode = newListItem
		l.lastNode = newListItem
	} else {
		newListItem.Prev = l.lastNode
		l.lastNode.Next = newListItem
		l.lastNode = newListItem
	}
	l.lenghtList++
	return newListItem
}

func (l *list) Remove(item *ListItem) {
	left := item.Prev
	right := item.Next

	if left != nil {
		left.Next = right
	} else {
		l.firstNode = right
	}

	if right != nil {
		right.Prev = left
	} else {
		l.lastNode = left
	}

	l.lenghtList--
}

func (l *list) MoveToFront(item *ListItem) {
	if item.Prev != nil {
		// itemCopyValue := item.Value
		l.PushFront(item.Value)
		l.Remove(item)
	}
}

func NewList() List {
	return new(list)
}
