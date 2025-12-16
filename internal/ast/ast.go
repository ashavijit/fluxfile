package ast

type FluxFile struct {
	Vars      map[string]string
	Tasks     []Task
	Profiles  []Profile
	Includes  []string
	Templates []Template
	Groups    []TaskGroup
	Aliases   map[string]string
}

type Task struct {
	Name        string
	Desc        string
	Deps        []string
	Parallel    bool
	If          string
	Run         []string
	Env         map[string]string
	Watch       []string
	WatchIgnore []string
	Matrix      *Matrix
	Cache       bool
	Inputs      []string
	Outputs     []string
	Docker      bool
	Remote      string
	Profile     string
	Secrets     []string
	Pre         []Precondition
	Retries     int
	RetryDelay  string
	Timeout     string
	Prompt      string
	Notify      NotifyConfig
	Alias       string
	Extends     string
	Before      []string
	After       []string
}

type NotifyConfig struct {
	Success string
	Failure string
}

type Precondition struct {
	Type  string
	Value string
}

type Profile struct {
	Name string
	Env  map[string]string
}

type Matrix struct {
	Dimensions map[string][]string
}

type Template struct {
	Name       string
	Desc       string
	Deps       []string
	Env        map[string]string
	Cache      bool
	Inputs     []string
	Outputs    []string
	Parallel   bool
	Docker     bool
	Remote     string
	Secrets    []string
	Pre        []Precondition
	Retries    int
	RetryDelay string
	Timeout    string
	Before     []string
	After      []string
}

type TaskGroup struct {
	Name  string
	Tasks []string
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
		Vars:      make(map[string]string),
		Tasks:     []Task{},
		Profiles:  []Profile{},
		Includes:  []string{},
		Templates: []Template{},
		Groups:    []TaskGroup{},
		Aliases:   make(map[string]string),
	}
}

func NewTask(name string) Task {
	return Task{
		Name:        name,
		Desc:        "",
		Deps:        []string{},
		Parallel:    false,
		If:          "",
		Run:         []string{},
		Env:         make(map[string]string),
		Watch:       []string{},
		WatchIgnore: []string{},
		Matrix:      nil,
		Cache:       false,
		Inputs:      []string{},
		Outputs:     []string{},
		Docker:      false,
		Remote:      "",
		Profile:     "",
		Secrets:     []string{},
		Pre:         []Precondition{},
		Retries:     0,
		RetryDelay:  "",
		Timeout:     "",
		Prompt:      "",
		Notify:      NotifyConfig{},
		Alias:       "",
		Extends:     "",
		Before:      []string{},
		After:       []string{},
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

// NewTemplate creates a new Template with default values
func NewTemplate(name string) Template {
	return Template{
		Name:       name,
		Desc:       "",
		Deps:       []string{},
		Env:        make(map[string]string),
		Cache:      false,
		Inputs:     []string{},
		Outputs:    []string{},
		Parallel:   false,
		Docker:     false,
		Remote:     "",
		Secrets:    []string{},
		Pre:        []Precondition{},
		Retries:    0,
		RetryDelay: "",
		Timeout:    "",
		Before:     []string{},
		After:      []string{},
	}
}

// NewTaskGroup creates a new TaskGroup with default values
func NewTaskGroup(name string) TaskGroup {
	return TaskGroup{
		Name:  name,
		Tasks: []string{},
	}
}
