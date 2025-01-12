package secret

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Secret represents a stored secret with metadata
type Secret struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description,omitempty"`
	Data        string    `yaml:"data"` // base64 encoded encrypted data
	CreatedAt   time.Time `yaml:"created_at"`
	Tags        []string  `yaml:"tags,omitempty"`
	Type        string    `yaml:"type,omitempty"`   // optional type of secret (e.g., "api-key", "certificate")
	Format      string    `yaml:"format,omitempty"` // original format of the secret (e.g., "text", "json", "binary")
}

// New creates a new Secret with the given name and encrypted data
func New(name string, encryptedData []byte) *Secret {
	now := time.Now()
	return &Secret{
		Name:      name,
		Data:      base64.StdEncoding.EncodeToString(encryptedData),
		CreatedAt: now,
	}
}

// Save writes the secret to a YAML file
func (s *Secret) Save(path string) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

// Load reads a secret from a YAML file
func Load(path string) (*Secret, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret file: %w", err)
	}

	var secret Secret
	if err := yaml.Unmarshal(data, &secret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return &secret, nil
}

// GetData returns the decoded encrypted data
func (s *Secret) GetData() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.Data)
}
