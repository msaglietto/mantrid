package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

var (
	ErrEmptyAliasName    = errors.New("alias name cannot be empty")
	ErrEmptyAliasCommand = errors.New("alias command cannot be empty")
	ErrAliasNotFound     = errors.New("alias not found")
	ErrAliasExists       = errors.New("alias already exists")
	ErrInvalidAliasName  = errors.New("alias name must contain only alphanumeric characters, hyphens, and underscores")
	ErrNameTooLong       = errors.New("alias name must be 64 characters or fewer")
	ErrCommandTooLong    = errors.New("alias command must be 4096 characters or fewer")
)

const (
	maxNameLength    = 64
	maxCommandLength = 4096
)

var aliasNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type Alias struct {
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewAlias(name, command string) (*Alias, error) {
	if err := validateAlias(name, command); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Alias{
		Name:      name,
		Command:   command,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func validateAlias(name, command string) error {
	if name == "" {
		return ErrEmptyAliasName
	}
	if len(name) > maxNameLength {
		return ErrNameTooLong
	}
	if !aliasNamePattern.MatchString(name) {
		return fmt.Errorf("%w: %q", ErrInvalidAliasName, name)
	}
	if command == "" {
		return ErrEmptyAliasCommand
	}
	if len(command) > maxCommandLength {
		return ErrCommandTooLong
	}
	return nil
}

// UpdateCommand updates the command and timestamp of an alias
func (a *Alias) UpdateCommand(newCommand string) error {
	if newCommand == "" {
		return ErrEmptyAliasCommand
	}
	if len(newCommand) > maxCommandLength {
		return ErrCommandTooLong
	}
	a.Command = newCommand
	a.UpdatedAt = time.Now()
	return nil
}
