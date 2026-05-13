package integrationtests_test

import (
	"context"
	"strings"
	"testing"

	"github.com/langgenius/dify-sandbox/internal/core/runner/types"
	"github.com/langgenius/dify-sandbox/internal/service"
	"github.com/langgenius/dify-sandbox/internal/static"
)

// setExposeEnvs sets ExposeEnvs for the duration of the test and restores the original value on cleanup.
func setExposeEnvs(t *testing.T, envs []string) {
	t.Helper()
	original := static.GetDifySandboxGlobalConfigurations().ExposeEnvs
	static.SetExposeEnvs(envs)
	t.Cleanup(func() {
		static.SetExposeEnvs(original)
	})
}

// TestPythonExposeEnvs verifies that an env var listed in ExposeEnvs is visible inside Python code.
func TestPythonExposeEnvs(t *testing.T) {
	t.Setenv("DIFY_TEST_EXPOSE_VAR", "hello_from_expose_envs")
	setExposeEnvs(t, []string{"DIFY_TEST_EXPOSE_VAR"})

	resp := service.RunPython3Code(context.TODO(), `
import os
print(os.environ.get("DIFY_TEST_EXPOSE_VAR", "not_found"))
`, "", &types.RunnerOptions{
		EnableNetwork: false,
	})

	if resp.Code != 0 {
		t.Fatal(resp)
	}
	if resp.Data.(*service.RunCodeResponse).Stderr != "" {
		t.Fatalf("unexpected stderr: %s", resp.Data.(*service.RunCodeResponse).Stderr)
	}
	if !strings.Contains(resp.Data.(*service.RunCodeResponse).Stdout, "hello_from_expose_envs") {
		t.Fatalf("expected env var value in output, got: %s", resp.Data.(*service.RunCodeResponse).Stdout)
	}
}

// TestPythonExposeEnvs_NotExposed verifies that an env var NOT listed in ExposeEnvs is not visible inside Python code.
func TestPythonExposeEnvs_NotExposed(t *testing.T) {
	t.Setenv("DIFY_TEST_UNEXPOSED_VAR", "should_not_appear")
	setExposeEnvs(t, []string{}) // not including DIFY_TEST_UNEXPOSED_VAR

	resp := service.RunPython3Code(context.TODO(), `
import os
print(os.environ.get("DIFY_TEST_UNEXPOSED_VAR", "not_found"))
`, "", &types.RunnerOptions{
		EnableNetwork: false,
	})

	if resp.Code != 0 {
		t.Fatal(resp)
	}
	if strings.Contains(resp.Data.(*service.RunCodeResponse).Stdout, "should_not_appear") {
		t.Fatalf("env var should not be exposed, but got: %s", resp.Data.(*service.RunCodeResponse).Stdout)
	}
}

// TestPythonExposeEnvs_Missing verifies that a var listed in ExposeEnvs but absent from the container is simply not set.
func TestPythonExposeEnvs_Missing(t *testing.T) {
	setExposeEnvs(t, []string{"DIFY_TEST_NONEXISTENT_VAR_XYZ"})

	resp := service.RunPython3Code(context.TODO(), `
import os
print(os.environ.get("DIFY_TEST_NONEXISTENT_VAR_XYZ", "not_found"))
`, "", &types.RunnerOptions{
		EnableNetwork: false,
	})

	if resp.Code != 0 {
		t.Fatal(resp)
	}
	if !strings.Contains(resp.Data.(*service.RunCodeResponse).Stdout, "not_found") {
		t.Fatalf("missing env var should fall back to default, got: %s", resp.Data.(*service.RunCodeResponse).Stdout)
	}
}

// TestNodejsExposeEnvs verifies that an env var listed in ExposeEnvs is visible inside NodeJS code.
func TestNodejsExposeEnvs(t *testing.T) {
	t.Setenv("DIFY_TEST_EXPOSE_VAR", "hello_from_expose_envs")
	setExposeEnvs(t, []string{"DIFY_TEST_EXPOSE_VAR"})

	resp := service.RunNodeJsCode(context.TODO(), `
console.log(process.env.DIFY_TEST_EXPOSE_VAR || "not_found");
`, "", &types.RunnerOptions{
		EnableNetwork: false,
	})

	if resp.Code != 0 {
		t.Fatal(resp)
	}
	if resp.Data.(*service.RunCodeResponse).Stderr != "" {
		t.Fatalf("unexpected stderr: %s", resp.Data.(*service.RunCodeResponse).Stderr)
	}
	if !strings.Contains(resp.Data.(*service.RunCodeResponse).Stdout, "hello_from_expose_envs") {
		t.Fatalf("expected env var value in output, got: %s", resp.Data.(*service.RunCodeResponse).Stdout)
	}
}

// TestNodejsExposeEnvs_NotExposed verifies that an env var NOT listed in ExposeEnvs is not visible inside NodeJS code.
func TestNodejsExposeEnvs_NotExposed(t *testing.T) {
	t.Setenv("DIFY_TEST_UNEXPOSED_VAR", "should_not_appear")
	setExposeEnvs(t, []string{}) // not including DIFY_TEST_UNEXPOSED_VAR

	resp := service.RunNodeJsCode(context.TODO(), `
console.log(process.env.DIFY_TEST_UNEXPOSED_VAR || "not_found");
`, "", &types.RunnerOptions{
		EnableNetwork: false,
	})

	if resp.Code != 0 {
		t.Fatal(resp)
	}
	if strings.Contains(resp.Data.(*service.RunCodeResponse).Stdout, "should_not_appear") {
		t.Fatalf("env var should not be exposed, but got: %s", resp.Data.(*service.RunCodeResponse).Stdout)
	}
}
