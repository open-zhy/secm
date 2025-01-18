package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DirName     = ".secm"
	SecretsDir  = "secrets"
	IdentityKey = "identity.key"
)

// Workspace represents the secm workspace configuration
type Workspace struct {
	RootDir    string
	SecretsDir string
	KeyPath    string
}

// Initialize creates the workspace directory structure
func Initialize(profile string) (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	ws := &Workspace{
		RootDir:    filepath.Join(homeDir, DirName, profile),
		SecretsDir: filepath.Join(homeDir, DirName, profile, SecretsDir),
		KeyPath:    filepath.Join(homeDir, DirName, profile, IdentityKey),
	}

	// Check if workspace exists, if so we don't override. instead we throw error
	if _, err := os.Stat(ws.KeyPath); err == nil {
		return nil, fmt.Errorf("workspace already initialized, use \"--profile <new-profile>\" to initialize a new workspace")
	}

	// Create root directory
	if err := os.MkdirAll(ws.RootDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Create secrets directory
	if err := os.MkdirAll(ws.SecretsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create secrets directory: %w", err)
	}

	return ws, nil
}

// Load returns an existing workspace configuration
func Load(profile string) (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	ws := &Workspace{
		RootDir:    filepath.Join(homeDir, DirName, profile),
		SecretsDir: filepath.Join(homeDir, DirName, profile, SecretsDir),
		KeyPath:    filepath.Join(homeDir, DirName, profile, IdentityKey),
	}

	// Check if workspace exists
	if _, err := os.Stat(ws.RootDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace not initialized, run 'secm init' first")
	}

	return ws, nil
}

// SecretPath returns the full path for a secret with the given ID
func (w *Workspace) SecretPath(id string) string {
	return filepath.Join(w.SecretsDir, id)
}
