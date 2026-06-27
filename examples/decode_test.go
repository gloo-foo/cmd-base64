package base64_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-base64"
	"github.com/gloo-foo/testable"
)

func ExampleBase64_decode() {
	// echo "SGVsbG8gV29ybGQ=" | base64 -d
	output, _ := testable.Test(command.Base64(command.Base64Decode), "SGVsbG8gV29ybGQ=\n")
	fmt.Print(output)
	// Output:
	// Hello World
}
