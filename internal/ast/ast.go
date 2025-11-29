package ast

type FluxFile struct {
	Vars     map[string]string
	Tasks    []Task
	Profiles []Profile
	Includes []string
}

type Task struct {
	Name   string
	Deps   []string
	Run    []string
	Env    map[string]string
	Watch  []string
	Matrix *Matrix
	Docker bool
	Remote string
}

type Profile struct {
	Name string
	Env  map[string]string
}

type Matrix struct {
	Dimensions map[string][]string
}

type Expr interface {
	exprNode()
}

type StringLiteral struct {
	Value string
}

type NumberLiteral struct {
	Value string
}

type Identifier struct {
	Value string
}

type ShellExpr struct {
	Command string
}

func (s *StringLiteral) exprNode() {}
func (n *NumberLiteral) exprNode() {}
func (i *Identifier) exprNode()    {}
func (s *ShellExpr) exprNode()     {}

func NewFluxFile() *FluxFile {
	return &FluxFile{
		Vars:     make(map[string]string),
		Tasks:    []Task{},
		Profiles: []Profile{},
		Includes: []string{},
	}
}

func NewTask(name string) Task {
	return Task{
		Name:   name,
		Deps:   []string{},
		Run:    []string{},
		Env:    make(map[string]string),
		Watch:  []string{},
		Matrix: nil,
		Docker: false,
		Remote: "",
	}
}

func NewProfile(name string) Profile {
	return Profile{
		Name: name,
		Env:  make(map[string]string),
	}
}

func NewMatrix() *Matrix {
	return &Matrix{
		Dimensions: make(map[string][]string),
	}
}
