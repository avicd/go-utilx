package datax

type LinkedNode[T any] struct {
	Item T
	Pre  *LinkedNode[T]
	Next *LinkedNode[T]
}

func (it *LinkedNode[T]) Append(item T) *LinkedNode[T] {
	next := &LinkedNode[T]{
		Item: item,
		Pre:  it,
	}
	if it != nil {
		it.Next = next
	}
	return next
}

func (it *LinkedNode[T]) Remove() {
	if it == nil {
		return
	}
	if it.Pre != nil {
		it.Pre.Next = it.Next
	}
	if it.Next != nil {
		it.Next.Pre = it.Pre
	}
	it.Next = nil
	it.Pre = nil
}

func (it *LinkedNode[T]) MoveTo(ref *LinkedNode[T]) {
	if it == ref || it == nil || ref == nil {
		return
	}
	it.Remove()
	ref.Insert(it)
}

func (it *LinkedNode[T]) Insert(node *LinkedNode[T]) {
	if it == nil || node == nil {
		return
	}
	if it.Pre != nil {
		it.Pre.Next = node
	}
	node.Pre = it.Pre
	it.Pre = node
	node.Next = it
}

func (it *LinkedNode[T]) InsertAfter(node *LinkedNode[T]) {
	if it == nil || node == nil {
		return
	}
	if it.Next == nil {
		it.Next = node
		node.Pre = it
	} else {
		it.Next.Insert(node)
	}
}
