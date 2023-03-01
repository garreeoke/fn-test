package main

import (
	"fmt"
	fnv1alpha1 "github.com/crossplane/crossplane/apis/apiextensions/fn/io/v1alpha1"
	"github.com/pkg/errors"
	"io"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

type Datasource struct {
	Composite fnv1alpha1.ObservedComposite           `json:"Composite"`
	Resources map[string]fnv1alpha1.ObservedResource `json:resources,omitempty`
}

func main() {

	stdin, err := GetStdIn()
	if err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.New("Cannot read stdin"))
		return
	}

	// Unmarshal the function IO.
	in := fnv1alpha1.FunctionIO{}
	if err := yaml.Unmarshal([]byte(strings.TrimSpace(string(stdin))), &in); err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.Wrap(err, "cannot unmarshal as FunctionIO"))
		return
	}

	out := *(in.DeepCopy())
	output, err := yaml.Marshal(out)
	if err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.Wrap(err, "cannot marshal FunctionIO"))
		return
	}

	fmt.Println(string(output))
}

func BuildFunctionIo(stdin []byte) error {

	// Unmarshal the function IO.
	in := fnv1alpha1.FunctionIO{}
	if err := yaml.Unmarshal([]byte(strings.TrimSpace(string(stdin))), &in); err != nil {
		failFatal(fnv1alpha1.FunctionIO{}, errors.Wrap(err, "cannot unmarshal as FunctionIO"))
		return err
	}
	return nil
}

func GetStdIn() ([]byte, error) {
	// Get stdin
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	if len(stdin) < 1 {
		return stdin, errors.New("Standard in is zero")
	}
	return stdin, nil
}

func failFatal(io fnv1alpha1.FunctionIO, err error) {
	io.Results = append(io.Results, fnv1alpha1.Result{
		Severity: fnv1alpha1.SeverityFatal,
		Message:  err.Error(),
	})
	b, _ := yaml.Marshal(io)
	fmt.Println(string(b))
}
