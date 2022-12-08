package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Node interface {
	name() string
	size() int
	parent() *Dir
	print(indent int)
}

type File struct {
	filename  string
	filesize  int
	parentDir *Dir
}

func (f *File) name() string {
	return f.filename
}

func (f *File) size() int {
	return f.filesize
}

func (f *File) parent() *Dir {
	return f.parentDir
}

func (f *File) print(indent int) {
	fmt.Println(strings.Repeat(" ", indent), f.filesize, " ", f.filename)
}

func createFile(name string, size int, parent *Dir) File {
	return File{
		filename:  name,
		filesize:  size,
		parentDir: parent,
	}
}

type Dir struct {
	dirname   string
	nodes     []Node
	parentDir *Dir
}

func (d *Dir) name() string {
	return d.dirname
}

func (d *Dir) size() int {
	var total int
	for _, d := range d.nodes {
		total = total + d.size()
	}
	return total
}

func (d *Dir) parent() *Dir {
	return d.parentDir
}

func (d *Dir) print(indent int) {
	fmt.Println(strings.Repeat(" ", indent), "- ", d.dirname)
	indent = indent + 2
	for _, n := range d.nodes {
		n.print(indent)
	}
}

func (d *Dir) addNode(n Node) {
	d.nodes = append(d.nodes, n)
}

func createDir(n string, parentDir *Dir) Dir {
	return Dir{
		dirname:   n,
		nodes:     make([]Node, 0),
		parentDir: parentDir,
	}
}

func parse(lines []string) Dir {
	topMost := createDir("/", nil)
	var currentDir *Dir = &topMost

	for i := 1; i < len(lines); i++ {
		l := lines[i]

		if !strings.HasPrefix(l, "$") {
			fmt.Errorf("the line %v is not a proper command", l)
		}

		command := l[2:]
		if strings.HasPrefix(command, "cd") {
			dir := command[3:]
			if dir == "/" {
				currentDir = &topMost
			} else if dir == ".." {
				currentDir = currentDir.parent()
			} else {
				newDir := createDir(dir, currentDir)
				currentDir.addNode(&newDir)
				currentDir = &newDir
			}
		} else if command == "ls" {
			i++
			for i < len(lines) && !strings.HasPrefix(lines[i], "$") {
				split := strings.Split(lines[i], " ")
				if split[0] != "dir" {
					size, _ := strconv.Atoi(split[0])
					file := createFile(split[1], size, currentDir)
					currentDir.addNode(&file)
				}
				i++
			}
			i--
		} else {
			fmt.Errorf("unknown command %v", command)
		}
	}
	return topMost
}

type Visitor interface {
	visitFile(f *File)
	visitDir(f *Dir)
}

func visit(v Visitor, n Node) {
	if d, ok := n.(*Dir); ok {
		v.visitDir(d)
	} else if f, ok := n.(*File); ok {
		v.visitFile(f)
	} else {
		fmt.Errorf("unknown node %v", n)
	}
}

func descend(v Visitor, n Node) {
	if d, ok := n.(*Dir); ok {
		for _, i := range d.nodes {
			visit(v, i)
		}
	}
}

type SizeLessThan100k struct {
	totalSize int
}

func (v *SizeLessThan100k) visitFile(f *File) {}
func (v *SizeLessThan100k) visitDir(d *Dir) {
	size := d.size()
	if d.size() <= 100000 {
		v.totalSize += size
	}
	descend(v, d)
}

type FindDirToDelete struct {
	targetSize int
	foundSize  int
}

func (v *FindDirToDelete) visitFile(f *File) {}
func (v *FindDirToDelete) visitDir(d *Dir) {
	size := d.size()
	if size >= v.targetSize {
		if size < v.foundSize {
			v.foundSize = size
		}
	}
	descend(v, d)
}

func main() {
	rawContent, err := os.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	content := string(rawContent)
	lines := strings.Split(content, "\n")

	tree := parse(lines)

	part1 := SizeLessThan100k{totalSize: 0}
	descend(&part1, &tree)
	fmt.Println("Part 1: ", part1.totalSize)

	part2 := FindDirToDelete{
		targetSize: 30000000 - 70000000 + tree.size(),
		foundSize:  math.MaxInt,
	}
	descend(&part2, &tree)
	fmt.Println("Part 2: ", part2.foundSize)

}
