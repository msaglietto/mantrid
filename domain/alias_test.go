package domain_test

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, err := domain.NewAlias(tt.aliasName, tt.command)
			assert.Equal(t, tt.expectedErr, err)
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
