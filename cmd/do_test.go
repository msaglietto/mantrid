package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDoArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedAlias  string
		expectedParams []string
	}{
		{
			name:           "no separator - no params",
			args:           []string{"hello"},
			expectedAlias:  "hello",
			expectedParams: []string{},
		},
		{
			name:           "no separator - with params",
			args:           []string{"greet", "Alice"},
			expectedAlias:  "greet",
			expectedParams: []string{"Alice"},
		},
		{
			name:           "no separator - multiple params",
			args:           []string{"greet", "Alice", "Bob"},
			expectedAlias:  "greet",
			expectedParams: []string{"Alice", "Bob"},
		},
		{
			name:           "with separator - no params",
			args:           []string{"hello", "--"},
			expectedAlias:  "hello",
			expectedParams: []string{},
		},
		{
			name:           "with separator - single param",
			args:           []string{"ls", "--", "-la"},
			expectedAlias:  "ls",
			expectedParams: []string{"-la"},
		},
		{
			name:           "with separator - multiple params with flags",
			args:           []string{"docker", "--", "run", "--rm", "-it", "ubuntu"},
			expectedAlias:  "docker",
			expectedParams: []string{"run", "--rm", "-it", "ubuntu"},
		},
		{
			name:           "with separator - params look like flags",
			args:           []string{"grep", "--", "-r", "pattern", "."},
			expectedAlias:  "grep",
			expectedParams: []string{"-r", "pattern", "."},
		},
		{
			name:           "with separator - single flag param",
			args:           []string{"ls", "--", "--help"},
			expectedAlias:  "ls",
			expectedParams: []string{"--help"},
		},
		{
			name:           "with separator - complex docker command",
			args:           []string{"docker", "--", "run", "-p", "8080:80", "--name", "webserver", "nginx"},
			expectedAlias:  "docker",
			expectedParams: []string{"run", "-p", "8080:80", "--name", "webserver", "nginx"},
		},
		{
			name:           "with separator - kubectl command",
			args:           []string{"k", "--", "get", "pods", "-n", "production", "--watch"},
			expectedAlias:  "k",
			expectedParams: []string{"get", "pods", "-n", "production", "--watch"},
		},
		{
			name:           "empty args",
			args:           []string{},
			expectedAlias:  "",
			expectedParams: []string{},
		},
		{
			name:           "with separator - git command",
			args:           []string{"git", "--", "commit", "-m", "feat: new feature"},
			expectedAlias:  "git",
			expectedParams: []string{"commit", "-m", "feat: new feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliasName, params := parseDoArgs(tt.args)
			assert.Equal(t, tt.expectedAlias, aliasName)
			assert.Equal(t, tt.expectedParams, params)
		})
	}
}

func TestSubstituteParams(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		params   []string
		expected string
	}{
		{
			name:     "no parameters",
			command:  "echo hello",
			params:   []string{},
			expected: "echo hello",
		},
		{
			name:     "single positional parameter",
			command:  "echo Hello, $1!",
			params:   []string{"World"},
			expected: "echo Hello, World!",
		},
		{
			name:     "multiple positional parameters",
			command:  "deploy $1 to $2",
			params:   []string{"app", "production"},
			expected: "deploy app to production",
		},
		{
			name:     "all parameters with $@",
			command:  "ls -la $@",
			params:   []string{"/home", "/tmp"},
			expected: "ls -la /home /tmp",
		},
		{
			name:     "all parameters with $*",
			command:  "echo $*",
			params:   []string{"hello", "world"},
			expected: "echo hello world",
		},
		{
			name:     "mixed positional and all with $@",
			command:  "kubectl apply -f $1 $@",
			params:   []string{"app.yaml", "--namespace", "prod"},
			expected: "kubectl apply -f app.yaml app.yaml --namespace prod",
		},
		{
			name:     "parameter used multiple times",
			command:  "echo $1 and $1 again",
			params:   []string{"test"},
			expected: "echo test and test again",
		},
		{
			name:     "parameters with special characters",
			command:  "grep $1 $2",
			params:   []string{"pattern-test", "file.txt"},
			expected: "grep pattern-test file.txt",
		},
		{
			name:     "command with $@ and no params",
			command:  "ls -la $@",
			params:   []string{},
			expected: "ls -la $@",
		},
		{
			name:     "command with $* and no params",
			command:  "echo $*",
			params:   []string{},
			expected: "echo $*",
		},
		{
			name:     "complex kubectl command",
			command:  "kubectl apply -f $1 --namespace $2",
			params:   []string{"deployment.yaml", "production"},
			expected: "kubectl apply -f deployment.yaml --namespace production",
		},
		{
			name:     "git command with multiple params",
			command:  "git commit -m $1",
			params:   []string{"feat: add new feature"},
			expected: "git commit -m feat: add new feature",
		},
		{
			name:     "docker command with positional and all params",
			command:  "docker run $1 $@",
			params:   []string{"nginx", "-p", "8080:80", "-d"},
			expected: "docker run nginx nginx -p 8080:80 -d",
		},
		{
			name:     "ssh command with parameters",
			command:  "ssh $1@$2",
			params:   []string{"user", "example.com"},
			expected: "ssh user@example.com",
		},
		{
			name:     "empty params list",
			command:  "echo $1",
			params:   []string{},
			expected: "echo $1",
		},
		{
			name:     "parameter numbers beyond provided params",
			command:  "echo $1 $2 $3",
			params:   []string{"first"},
			expected: "echo first $2 $3",
		},
		{
			name:     "single $@ in middle of command",
			command:  "find $@ -name '*.go'",
			params:   []string{"/src", "/tests"},
			expected: "find /src /tests -name '*.go'",
		},
		{
			name:     "multiple occurrences of $@",
			command:  "echo $@ and again $@",
			params:   []string{"a", "b"},
			expected: "echo a b and again a b",
		},
		{
			name:     "both $@ and $* in same command",
			command:  "echo $@ then $*",
			params:   []string{"x", "y"},
			expected: "echo x y then x y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteParams(tt.command, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSubstituteParams_AutoAppend(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		params   []string
		expected string
	}{
		{
			name:     "no placeholders - auto append single param",
			command:  "ls",
			params:   []string{"-la"},
			expected: "ls -la",
		},
		{
			name:     "no placeholders - auto append multiple params",
			command:  "docker",
			params:   []string{"ps", "-a"},
			expected: "docker ps -a",
		},
		{
			name:     "no placeholders - auto append flags and args",
			command:  "kubectl get pods",
			params:   []string{"-n", "production", "--watch"},
			expected: "kubectl get pods -n production --watch",
		},
		{
			name:     "has $1 placeholder - substitute not append",
			command:  "echo Hello, $1!",
			params:   []string{"World"},
			expected: "echo Hello, World!",
		},
		{
			name:     "has $@ placeholder - substitute not append",
			command:  "grep -r $@ .",
			params:   []string{"TODO"},
			expected: "grep -r TODO .",
		},
		{
			name:     "has $* placeholder - substitute not append",
			command:  "echo $*",
			params:   []string{"hello", "world"},
			expected: "echo hello world",
		},
		{
			name:     "no placeholders - no params - no change",
			command:  "ls",
			params:   []string{},
			expected: "ls",
		},
		{
			name:     "complex command no placeholders - auto append",
			command:  "git log --oneline",
			params:   []string{"-10", "--author=me"},
			expected: "git log --oneline -10 --author=me",
		},
		{
			name:     "simple alias with dash flag",
			command:  "ls",
			params:   []string{"-la", "/tmp"},
			expected: "ls -la /tmp",
		},
		{
			name:     "docker alias with subcommand and flags",
			command:  "docker",
			params:   []string{"run", "--rm", "-it", "ubuntu"},
			expected: "docker run --rm -it ubuntu",
		},
		{
			name:     "kubectl alias with complex flags",
			command:  "kubectl",
			params:   []string{"get", "pods", "-n", "production"},
			expected: "kubectl get pods -n production",
		},
		{
			name:     "git alias with single param",
			command:  "git status",
			params:   []string{"--short"},
			expected: "git status --short",
		},
		{
			name:     "grep alias with pattern and path",
			command:  "grep -r",
			params:   []string{"TODO", "."},
			expected: "grep -r TODO .",
		},
		{
			name:     "npm alias with subcommand",
			command:  "npm",
			params:   []string{"install", "--save-dev", "typescript"},
			expected: "npm install --save-dev typescript",
		},
		{
			name:     "alias with equals sign in params",
			command:  "git log",
			params:   []string{"--author=John Doe", "-10"},
			expected: "git log --author=John Doe -10",
		},
		{
			name:     "mixed $2 present - should substitute not append",
			command:  "echo $2",
			params:   []string{"first", "second"},
			expected: "echo second",
		},
		{
			name:     "has $3 - should substitute not append",
			command:  "deploy $1 $2 $3",
			params:   []string{"app", "env", "region"},
			expected: "deploy app env region",
		},
		{
			name:     "no placeholders - single word param",
			command:  "cat",
			params:   []string{"file.txt"},
			expected: "cat file.txt",
		},
		{
			name:     "no placeholders - path param",
			command:  "cd",
			params:   []string{"/home/user/projects"},
			expected: "cd /home/user/projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := substituteParams(tt.command, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}
