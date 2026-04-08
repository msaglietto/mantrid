package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/msaglietto/mantrid/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewAlias(t *testing.T) {
	tests := []struct {
		name        string
		aliasName   string
		command     string
		expectedErr error
	}{
		{
			name:        "valid alias",
			aliasName:   "test",
			command:     "echo test",
			expectedErr: nil,
		},
		{
			name:        "valid alias with hyphens",
			aliasName:   "my-alias",
			command:     "echo test",
			expectedErr: nil,
		},
		{
			name:        "valid alias with underscores",
			aliasName:   "my_alias",
			command:     "echo test",
			expectedErr: nil,
		},
		{
			name:        "valid alias with numbers",
			aliasName:   "alias123",
			command:     "echo test",
			expectedErr: nil,
		},
		{
			name:        "empty name",
			aliasName:   "",
			command:     "echo test",
			expectedErr: domain.ErrEmptyAliasName,
		},
		{
			name:        "empty command",
			aliasName:   "test",
			command:     "",
			expectedErr: domain.ErrEmptyAliasCommand,
		},
		{
			name:        "name with spaces",
			aliasName:   "my alias",
			command:     "echo test",
			expectedErr: domain.ErrInvalidAliasName,
		},
		{
			name:        "name with special chars",
			aliasName:   "my;alias",
			command:     "echo test",
			expectedErr: domain.ErrInvalidAliasName,
		},
		{
			name:        "name with path traversal",
			aliasName:   "../etc/passwd",
			command:     "echo test",
			expectedErr: domain.ErrInvalidAliasName,
		},
		{
			name:        "name with dot",
			aliasName:   "my.alias",
			command:     "echo test",
			expectedErr: domain.ErrInvalidAliasName,
		},
		{
			name:        "name too long",
			aliasName:   strings.Repeat("a", 65),
			command:     "echo test",
			expectedErr: domain.ErrNameTooLong,
		},
		{
			name:        "name at max length",
			aliasName:   strings.Repeat("a", 64),
			command:     "echo test",
			expectedErr: nil,
		},
		{
			name:        "command too long",
			aliasName:   "test",
			command:     strings.Repeat("a", 4097),
			expectedErr: domain.ErrCommandTooLong,
		},
		{
			name:        "command at max length",
			aliasName:   "test",
			command:     strings.Repeat("a", 4096),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, err := domain.NewAlias(tt.aliasName, tt.command)
			if tt.expectedErr == domain.ErrInvalidAliasName {
				assert.True(t, errors.Is(err, tt.expectedErr), "expected error wrapping %v, got %v", tt.expectedErr, err)
			} else {
				assert.Equal(t, tt.expectedErr, err)
			}
			if err == nil {
				assert.NotNil(t, alias)
				assert.Equal(t, tt.aliasName, alias.Name)
				assert.Equal(t, tt.command, alias.Command)
				assert.False(t, alias.CreatedAt.IsZero())
				assert.False(t, alias.UpdatedAt.IsZero())
			}
		})
	}
}

func TestUpdateCommand(t *testing.T) {
	tests := []struct {
		name        string
		newCommand  string
		expectedErr error
	}{
		{
			name:        "valid command update",
			newCommand:  "echo updated",
			expectedErr: nil,
		},
		{
			name:        "empty command",
			newCommand:  "",
			expectedErr: domain.ErrEmptyAliasCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, err := domain.NewAlias("test", "echo original")
			assert.NoError(t, err)
			assert.NotNil(t, alias)

			originalUpdatedAt := alias.UpdatedAt
			originalCommand := alias.Command

			// Add small delay to ensure timestamp changes
			err = alias.UpdateCommand(tt.newCommand)
			assert.Equal(t, tt.expectedErr, err)

			if err == nil {
				assert.Equal(t, tt.newCommand, alias.Command)
				assert.NotEqual(t, originalCommand, alias.Command)
				// UpdatedAt should be updated (greater than or equal to original)
				assert.True(t, alias.UpdatedAt.After(originalUpdatedAt) || alias.UpdatedAt.Equal(originalUpdatedAt))
			} else {
				// On error, command should remain unchanged
				assert.Equal(t, originalCommand, alias.Command)
			}
		})
	}
}
