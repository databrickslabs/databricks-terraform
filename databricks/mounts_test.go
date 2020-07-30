package databricks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var executeMock func(clusterID, language, commandStr string) (string, error)

type commandExecutorMock struct{}

func (a commandExecutorMock) Execute(clusterID, language, commandStr string) (string, error) {
	return executeMock(clusterID, language, commandStr)
}

func TestValidateMountDirectory(t *testing.T) {
	testCases := []struct {
		directory  string
		errorCount int
	}{
		{"", 0},
		{"/directory", 0},
		{"directory", 1},
	}
	for _, tc := range testCases {
		_, errs := ValidateMountDirectory(tc.directory, "key")

		assert.Lenf(t, errs, tc.errorCount, "directory '%s' does not generate the expected error count", tc.directory)
	}
}
