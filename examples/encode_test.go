package base64_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-base64"
	"github.com/gloo-foo/testable"
)

func ExampleBase64_encode() {
	// echo "Hello World" | base64
	output, _ := testable.Test(command.Base64(), "Hello World\n")
	fmt.Print(output)
	// Output:
	// SGVsbG8gV29ybGQ=
}
