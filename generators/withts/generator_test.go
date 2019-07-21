package withts

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"go.uber.org/zap"
)

func TestGenerator_Generate(t *testing.T) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	logger, _ := config.Build()

	generator := New(logger)

	generator.options.Def()
	generator.options.URL = `postgres://genna:genna@localhost:5432/genna?sslmode=disable`
	generator.options.Output = path.Join(os.TempDir(), "model_test.go")
	generator.options.CreatedAt = "createdAt"
	generator.options.UpdatedAt = "updatedAt"
	generator.options.FollowFKs = true

	if err := generator.Generate(); err != nil {
		t.Errorf("generate error = %v", err)
		return
	}

	generated, err := ioutil.ReadFile(generator.options.Output)
	if err != nil {
		t.Errorf("file not generated = %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	check, err := ioutil.ReadFile(path.Join(path.Dir(filename), "generator_test.output"))
	if err != nil {
		t.Errorf("check file not found = %v", err)
	}

	if string(generated) != string(check) {
		t.Errorf("generated not mathed with check")
		return
	}
}
