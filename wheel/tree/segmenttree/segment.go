package segmenttree

import (
	"math"
)

type SegmentTree struct {
	nodes []int // 树节点
	size  int   // 原始数组大小
}

func calcTreeSize(dataSize int) int {
	return 1<<uint(math.Ceil(math.Log2(float64(dataSize)))+1) - 1
}

func NewSegmentTree(from []int) *SegmentTree {
	treeSize := calcTreeSize(len(from))
	nodes := make([]int, treeSize)

	tree := new(SegmentTree)
	tree.nodes = nodes
	tree.size = len(from)

	tree.build(from, 0, 0, tree.size-1)

	return tree
}

func (t *SegmentTree) build(data []int, root, istart, iend int) {
	if istart == iend {
		// 叶子节点
		t.nodes[root] = data[istart]
		return
	}
	// 左儿子，右儿子
	//2*root+1, 2*root + 2
	mid := (istart + iend) / 2
	t.build(data, 2*root+1, istart, mid)
	t.build(data, 2*root+2, mid+1, iend)

	leftMin := t.nodes[2*root+1]
	rightMin := t.nodes[2*root+2]

	if leftMin < rightMin {
		t.nodes[root] = leftMin
	} else {
		t.nodes[root] = rightMin
	}
}

type query struct {
	left, right int
	nodes       []int
}

func (q *query) rangeMinimum(root, istart, iend int) int {
	// 查询区间和当前节点区间没有交集
	if q.left > iend || q.right < istart {
		return math.MaxInt64
	}
	// 这个的情形是：我们在不断递归查找的过程中，一定会找到 left == istart,  right == iend 的区间，至于两者是否同时满足，则不一定
	// case study:
	// left, right = [55, 60]
	// istart, iend = [55, 57]
	// istart, iend = [58, 60]
	// 上面两个是实际二分出来的区间
	if q.left <= istart && q.right >= iend {
		return q.nodes[root]
	}

	bisect := (istart + iend) / 2
	leftMin := q.rangeMinimum(2*root+1, istart, bisect)
	rightMin := q.rangeMinimum(2*root+2, bisect+1, iend)
	if leftMin < rightMin {
		return leftMin
	}
	return rightMin
}

func (t *SegmentTree) RangeMinQuery(left, right int) int {
	// 根节点 0 存储的是 [0, len(from)-1] 中最小值
	// 根节点 1 存储的是 [0, (len(from)-1)/2] 中最小值
	// 根节点 2 存储的是 [(len(from)-1)/2+1, len(from)-1] 中最小值

	if left > right {
		left, right = right, left
	}
	return (&query{left: left, right: right, nodes: t.nodes}).rangeMinimum(0, 0, t.size-1)
}

func (q *query) update(root, istart, iend, value int) {
	// 查询区间和当前节点区间没有交集
	if q.left > iend || q.right < istart {
		return
	}

	if q.left == istart && q.right == iend{ // 最末的叶子节点了
		q.nodes[root] = value
		return
	}

	if value < q.nodes[root]  {
		// 更新最小值
		q.nodes[root] = value
	}

	mid := (istart + iend) / 2
	if q.left <= mid { // 下标在左侧
		q.update(2*root+1, istart, mid, value)
	} else {
		q.update(2*root+2, mid+1, iend, value)
	}

	return
}

func (t *SegmentTree) Update(value, index int) error {
	//
	if index >= t.size || index < 0 {
		return EBADINDEX
	}
	(&query{left: index, right: index, nodes: t.nodes}).update(0, 0, t.size-1, value)
	return nil
}
