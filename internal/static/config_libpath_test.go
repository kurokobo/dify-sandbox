//go:build linux

package static

import (
	"os"
	"reflect"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp config: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	f.Close()
	return f.Name()
}

// extra をYAML側で指定した場合に python_lib_path に追記されること
func TestPythonLibPathExtraInYaml(t *testing.T) {
	t.Setenv("PYTHON_LIB_PATH", "")
	t.Setenv("PYTHON_LIB_PATH_EXTRA", "")

	cfg := writeTempConfig(t, `
python_lib_path:
  - /lib/a
  - /lib/b
python_lib_path_extra:
  - /lib/c
`)
	if err := InitConfig(cfg); err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	got := GetDifySandboxGlobalConfigurations().PythonLibPaths
	want := []string{"/lib/a", "/lib/b", "/lib/c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// extra を環境変数で指定した場合に python_lib_path に追記されること
func TestPythonLibPathExtraViaEnv(t *testing.T) {
	t.Setenv("PYTHON_LIB_PATH", "/lib/a,/lib/b")
	t.Setenv("PYTHON_LIB_PATH_EXTRA", "/lib/c,/lib/d")

	cfg := writeTempConfig(t, `{}`)
	if err := InitConfig(cfg); err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	got := GetDifySandboxGlobalConfigurations().PythonLibPaths
	want := []string{"/lib/a", "/lib/b", "/lib/c", "/lib/d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// extra と base に重複があっても重複排除されること
func TestPythonLibPathDeduplication(t *testing.T) {
	t.Setenv("PYTHON_LIB_PATH", "")
	t.Setenv("PYTHON_LIB_PATH_EXTRA", "")

	cfg := writeTempConfig(t, `
python_lib_path:
  - /lib/a
  - /lib/b
  - /lib/a
python_lib_path_extra:
  - /lib/b
  - /lib/c
`)
	if err := InitConfig(cfg); err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	got := GetDifySandboxGlobalConfigurations().PythonLibPaths
	want := []string{"/lib/a", "/lib/b", "/lib/c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
