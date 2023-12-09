package datax

type LinkedList[T any] struct {
	head *LinkedNode[T]
	end  *LinkedNode[T]
	size int
}

func (it *LinkedList[T]) First() (T, bool) {
	var ret T
	if it.head != nil {
		ret = it.head.Item
	} else {
		return ret, false
	}
	return ret, true
}

func (it *LinkedList[T]) Last() (T, bool) {
	var ret T
	if it.end != nil {
		ret = it.end.Item
	} else {
		return ret, false
	}
	return ret, true
}

func (it *LinkedList[T]) Push(item T) int {
	it.end = it.end.Append(item)
	if it.head == nil {
		it.head = it.end
	}
	it.size++
	return it.size
}

func (it *LinkedList[T]) Pop() (T, bool) {
	var ret T
	if it.end == nil {
		return ret, false
	}
	ret = it.end.Item
	if it.head == it.end {
		it.head = nil
		it.end = nil
	} else if it.end.Pre != nil {
		it.end = it.end.Pre
		it.end.Next.Remove()
	} else {
		it.end = nil
	}
	it.size--
	return ret, true
}

func (it *LinkedList[T]) Unshift(item T) int {
	node := &LinkedNode[T]{Item: item}
	if it.head == nil {
		it.head = node
		it.end = node
	} else {
		it.head.Insert(node)
		it.head = node
	}
	it.size++
	return it.size
}

func (it *LinkedList[T]) Shift() (T, bool) {
	var ret T
	if it.head == nil {
		return ret, false
	}
	ret = it.end.Item
	if it.head == it.end {
		it.head = nil
		it.end = nil
	} else if it.head.Next != nil {
		it.head = it.head.Next
		it.head.Pre.Remove()
	} else {
		it.head = nil
	}
	it.size--
	return ret, true
}

func (it *LinkedList[T]) Len() int {
	return it.size
}

func (it *LinkedList[T]) nodeOf(index int) *LinkedNode[T] {
	var node *LinkedNode[T]
	if index < (it.size >> 1) {
		node = it.head
		for i := 0; i < index; i++ {
			node = node.Next
		}
	} else {
		node = it.end
		for i := it.size - 1; i > index; i-- {
			node = node.Pre
		}
	}
	return node
}

func (it *LinkedList[T]) Get(index int) (T, bool) {
	var ret T
	if it.size < 1 || index < 0 || index >= it.size {
		return ret, false
	} else {
		ret = it.nodeOf(index).Item
	}
	return ret, true
}

func (it *LinkedList[T]) Remove(index int) int {
	if it.size < 1 || index < 0 || index >= it.size {
		return it.size
	} else {
		node := it.nodeOf(index)
		if it.head == node {
			it.head = it.head.Next
		}
		if it.end == node {
			it.end = it.end.Pre
		}
		node.Remove()
		it.size--
	}
	return it.size
}

func (it *LinkedList[T]) Clear() {
	it.head = nil
	it.end = nil
	it.size = 0
}

func (it *LinkedList[T]) ForEach(fn func(i int, item T)) {
	index := 0
	for p := it.head; p != nil; p = p.Next {
		fn(index, p.Item)
		index++
	}
}
