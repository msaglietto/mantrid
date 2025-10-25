package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyAliasName    = errors.New("alias name cannot be empty")
	ErrEmptyAliasCommand = errors.New("alias command cannot be empty")
	ErrAliasNotFound     = errors.New("alias not found")
	ErrAliasExists       = errors.New("alias already exists")
)

type Alias struct {
	ID        string    `json:"id"`
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
	if command == "" {
		return ErrEmptyAliasCommand
	}
	return nil
}

// UpdateCommand updates the command and timestamp of an alias
func (a *Alias) UpdateCommand(newCommand string) error {
	if newCommand == "" {
		return ErrEmptyAliasCommand
	}
	a.Command = newCommand
	a.UpdatedAt = time.Now()
	return nil
}
