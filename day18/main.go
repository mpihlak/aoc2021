package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	isLeaf          bool
	value           int
	left, right, up *Node
}

func findMatchingBracketPos(s string) int {
	count := 0
	for pos, c := range s {
		if c == '[' {
			count += 1
		} else if c == ']' {
			count -= 1
			if count == 0 {
				return pos
			}
		}
	}
	return -1
}

func parseTree(s string, parent *Node) *Node {
	// must start with '[' and end with a matching ']'
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	this := Node{up: parent}

	// Now we have a pair of <node>, <node>
	// Parse the left node and advance "s"
	if strings.HasPrefix(s, "[") {
		endPos := findMatchingBracketPos(s)
		this.left = parseTree(s[:endPos+1], &this)
		s = s[endPos+2:]
	} else {
		val, err := strconv.Atoi(s[:1])
		if err != nil {
			panic(fmt.Sprintf("%v: %v", s[:1], err))
		}
		this.left = &Node{isLeaf: true, value: val, up: &this}
		s = s[2:]
	}

	// Parse the right node
	if strings.HasPrefix(s, "[") {
		this.right = parseTree(s, &this)
	} else {
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Sprintf("%v: %v", s[:1], err))
		}
		this.right = &Node{isLeaf: true, value: val, up: &this}
	}

	return &this
}

func formatTree(n *Node) string {
	buf := ""

	if n != nil {
		if n.isLeaf {
			buf += fmt.Sprintf("%v", n.value)
		} else {
			buf += fmt.Sprintf("[")
			buf += formatTree(n.left)
			buf += fmt.Sprintf(",")
			buf += formatTree(n.right)
			buf += fmt.Sprintf("]")
		}
	}

	return buf
}

func addToFirstLeaf(to *Node, val int, leftFirst bool) bool {
	if to == nil {
		return false
	}
	if to.isLeaf {
		to.value += val
		return true
	}

	var first, second *Node

	if leftFirst {
		first, second = to.left, to.right
	} else {
		first, second = to.right, to.left
	}

	if addToFirstLeaf(first, val, leftFirst) {
		return true
	}
	if addToFirstLeaf(second, val, leftFirst) {
		return true
	}

	return false
}

// Explode the leftmost leaf node at level 4, modify the tree in place
// Return true if exploded
func Explode(n *Node, level int) bool {
	if n == nil {
		return false
	}

	fmt.Printf("level=%v node=%p %+v\n", level, n, n)
	// XXX: Are we choosing the right pair, if there's multiple?

	if level == 4 {
		if n.left != nil && n.right != nil && n.left.isLeaf && n.right.isLeaf {
			leftValue := n.left.value
			rightValue := n.right.value

			n.isLeaf = true
			n.value = 0
			n.left = nil
			n.right = nil

			// Climb up the tree, add leftValue to the first left leaf found
			prev := n
			for b := n.up; b != nil; b = b.up {
				if b.left != prev {
					addToFirstLeaf(b.left, leftValue, false)
					break
				}
				prev = b
			}

			// Climb up the tree, add rightValue to the first right leaf found
			prev = n
			for b := n.up; b != nil; b = b.up {
				if b.right != prev {
					addToFirstLeaf(b.right, rightValue, true)
					break
				}
				prev = b
			}

			return true
		}
	} else {
		if Explode(n.left, level+1) {
			return true
		}
		if Explode(n.right, level+1) {
			return true
		}
	}

	return false
}

func Split(n *Node) bool {
	if n == nil {
		return false
	}

	if n.isLeaf && n.value >= 10 {
		l := n.value / 2
		r := n.value - l
		n.left = &Node{isLeaf: true, up: n, value: l}
		n.right = &Node{isLeaf: true, up: n, value: r}
		n.isLeaf = false
		n.value = 0
		return true
	} else {
		if Split(n.left) {
			return true
		}
		if Split(n.right) {
			return true
		}
	}

	return false
}

func Add(a, b *Node) *Node {
	n := &Node{left: a, right: b}
	a.up = n
	b.up = n
	return n
}

func Reduce(t *Node) {
	fmt.Println("Reducing", formatTree(t))
	for {
		reduceMore := false

		if Explode(t, 0) {
			fmt.Println("Exploded", formatTree(t))
			reduceMore = true
		}
		if Split(t) {
			fmt.Println("Split", formatTree(t))
			reduceMore = true
		}

		if !reduceMore {
			break
		}
	}
}

func Magnitude(t *Node) int64 {
	if t == nil {
		return 0
	}
	if t.isLeaf {
		return int64(t.value)
	}
	return 3*Magnitude(t.left) + 2*Magnitude(t.right)
}

func TestMagnitude() {
	t := parseTree("[[[[5,0],[7,4]],[5,5]],[6,6]]", nil)
	if Magnitude(t) != 1137 {
		panic("Magnitude test failed")
	}
}

func TestExplosions() {
	explosions := [][2]string{
		{"[[[[[9,8],1],2],3],4]", "[[[[0,9],2],3],4]"},
		{"[7,[6,[5,[4,[3,2]]]]]", "[7,[6,[5,[7,0]]]]"},
		{"[[6,[5,[4,[3,2]]]],1]", "[[6,[5,[7,0]]],3]"},
		{"[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]"},
		{"[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[7,0]]]]"},
	}

	for _, testCase := range explosions {
		n := parseTree(testCase[0], nil)
		Explode(n, 0)
		res := formatTree(n)

		if res != testCase[1] {
			fmt.Printf("Exploding %s is %s (expected %s)\n", testCase[0], res, testCase[1])
			os.Exit(1)
		}
	}
}

func TestSplitsAndExplosions() {
	a := parseTree("[[[[4,3],4],4],[7,[[8,4],9]]]", nil)
	b := parseTree("[1,1]", nil)
	t := Add(a, b)
	Explode(t, 0)
	if formatTree(t) != "[[[[0,7],4],[7,[[8,4],9]]],[1,1]]" {
		panic("Explosion failed")
	}
	Explode(t, 0)
	if formatTree(t) != "[[[[0,7],4],[15,[0,13]]],[1,1]]" {
		panic("Secondd Explode failed")
	}
	Split(t)
	if formatTree(t) != "[[[[0,7],4],[[7,8],[0,13]]],[1,1]]" {
		panic("Split failed")
	}
	Split(t)
	if formatTree(t) != "[[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]" {
		panic("Split failed")
	}
	Explode(t, 0)
	if formatTree(t) != "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]" {
		panic("Third Explode failed")
	}

	if Explode(t, 0) {
		panic("Should not have exploded")
	}
	if Split(t) {
		panic("Should not have split")
	}
}

func TestReduce() {
	t := parseTree("[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]", nil)
	Reduce(t)
	if formatTree(t) != "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]" {
		panic("Reduce failed")
	}
}

func TestAdd() {
	numbers := []string{"[1,1]", "[2,2]", "[3,3]", "[4,4]", "[5,5]", "[6,6]"}
	var t *Node
	for _, s := range numbers {
		n := parseTree(s, nil)
		if t == nil {
			t = n
		} else {
			t = Add(t, n)
			Reduce(t)
		}
	}
	if formatTree(t) != "[[[[5,0],[7,4]],[5,5]],[6,6]]" {
		panic("Add failed")
	}
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	line := "[[[[4,0],[0,4]],[[7,7],[6,0]]],[[[6,6],[5,0]],[[6,6],[8,[[5,6],8]]]]]"
	t := parseTree(line, nil)
	fmt.Println(formatTree(t))
	Explode(t, 0)
	fmt.Println(formatTree(t))

	//line = "[[[6,6],[5,0]],[[6,6],[8,[[5,6],8]]]]"
	//t = parseTree(line, nil)
	//Explode(t, 0)
	//fmt.Println(formatTree(t))
	os.Exit(1)

	TestExplosions()
	TestSplitsAndExplosions()
	TestReduce()
	TestAdd()
	TestMagnitude()

	buf, _ := os.ReadFile(fileName)
	input := string(buf)

	var sum *Node

	fmt.Println("Summing it up")
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}

		n := parseTree(line, nil)
		fmt.Println(formatTree(n))

		if sum == nil {
			sum = n
		} else {
			sum = Add(sum, n)
			Reduce(sum)
		}
	}

	fmt.Println()
	fmt.Println("Sum =", formatTree(sum))
	fmt.Println("Magnitude =", Magnitude(sum))
}
