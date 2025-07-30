package templates

import (
	_ "embed"
)

//go:embed files/readme/readme-with-db.tmpl
var helloWorldReadmeWithDb []byte

func HelloWorldReadmeWithDbTemplate() ([]byte, string) {
	return helloWorldReadmeWithDb, "README.md"
}

//go:embed files/cmd/main.tmpl
var helloWorldMain []byte

func HelloWorldMainTemplate() ([]byte, string) {
	return helloWorldMain, "cmd/main.go"
}

//go:embed files/cmd/main_test.tmpl
var helloWorldMainTest []byte

func HelloWorldMainTestTemplate() ([]byte, string) {
	return helloWorldMainTest, "cmd/main_test.go"
}

//go:embed files/cmd/main-with-db.tmpl
var helloWorldMainWithDb []byte

func HelloWorldMainWithDbTemplate() ([]byte, string) {
	return helloWorldMainWithDb, "cmd/main.go"
}
