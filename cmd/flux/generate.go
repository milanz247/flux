package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// stubData feeds the generator templates.
type stubData struct {
	Name      string // PascalCase, e.g. "Post"
	Snake     string // snake_case, e.g. "post"
	Camel     string // camelCase, e.g. "post"
	Plural    string // naive plural snake, e.g. "posts"
	Timestamp string // for migrations
	PagePath  string // for Vue pages, e.g. "Users/Index"
}

// makeFromStub generates a Go file from one of the named stubs.
func makeFromStub(args []string, kind string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: flux make:%s <Name>", kind)
	}
	name := pascal(args[0])
	data := newStubData(name)

	var path, stub string
	switch kind {
	case "controller":
		path = filepath.Join("app", "controllers", data.Snake+"_controller.go")
		stub = controllerStub
	case "model":
		path = filepath.Join("app", "models", data.Snake+".go")
		stub = modelStub
	case "service":
		path = filepath.Join("app", "services", data.Snake+"_service.go")
		stub = serviceStub
	case "dto":
		path = filepath.Join("app", "dto", data.Snake+".go")
		stub = dtoStub
	case "middleware":
		path = filepath.Join("app", "middleware", data.Snake+".go")
		stub = middlewareStub
	default:
		return fmt.Errorf("unknown generator %q", kind)
	}

	return writeStub(path, stub, data)
}

// makeMigration generates a timestamped migration file, e.g.
// `flux make:migration create_posts_table`.
func makeMigration(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: flux make:migration <name>  (e.g. create_posts_table)")
	}
	name := snake(args[0])
	data := newStubData(pascal(name))
	data.Timestamp = time.Now().Format("20060102150405")
	data.Snake = name

	path := filepath.Join("database", "migrations", data.Timestamp+"_"+name+".go")
	return writeStub(path, migrationStub, data)
}

// makePage generates a Vue page component, e.g. `flux make:page Users/Index`.
func makePage(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: flux make:page <Path>  (e.g. Users/Index)")
	}
	pagePath := strings.ReplaceAll(args[0], "\\", "/")
	data := newStubData(pascal(filepath.Base(pagePath)))
	data.PagePath = pagePath

	path := filepath.Join("resources", "js", "Pages", filepath.FromSlash(pagePath)+".vue")
	return writeStub(path, pageStub, data)
}

func writeStub(path, stub string, data stubData) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmpl, err := template.New("stub").Parse(stub)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	fmt.Printf("✓ created %s\n", path)
	return nil
}

func newStubData(name string) stubData {
	s := snake(name)
	return stubData{
		Name:   name,
		Snake:  s,
		Camel:  strings.ToLower(name[:1]) + name[1:],
		Plural: s + "s",
	}
}

// pascal converts "user_profile" / "userProfile" / "user profile" to "UserProfile".
func pascal(input string) string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var b strings.Builder
	for _, part := range parts {
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}
	return b.String()
}

// snake converts "UserProfile" to "user_profile".
func snake(input string) string {
	var b strings.Builder
	for i, r := range input {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteRune('_')
			}
			b.WriteRune(r - 'A' + 'a')
		} else if r == '-' || r == ' ' {
			b.WriteRune('_')
		} else {
			b.WriteRune(r)
		}
	}
	return strings.ReplaceAll(b.String(), "__", "_")
}
