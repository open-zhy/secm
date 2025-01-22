package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"plugin"
)

type PluginRegistry struct {
	Plugins map[string]string `json:"plugins"` // map of plugin name to .so file path
}

type Manager struct {
	registryPath string
	pluginsDir   string
	registry     *PluginRegistry
	loaded       map[string]*plugin.Plugin
}

func NewManager(pluginsDir string) *Manager {
	return &Manager{
		registryPath: filepath.Join(pluginsDir, "registry.json"),
		pluginsDir:   pluginsDir,
		registry:     &PluginRegistry{Plugins: make(map[string]string)},
		loaded:       make(map[string]*plugin.Plugin),
	}
}

func (m *Manager) loadRegistry() error {
	data, err := os.ReadFile(m.registryPath)
	if os.IsNotExist(err) {
		// Registry doesn't exist yet, that's fine
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read plugin registry: %w", err)
	}

	return json.Unmarshal(data, m.registry)
}

func (m *Manager) saveRegistry() error {
	data, err := json.MarshalIndent(m.registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plugin registry: %w", err)
	}

	return os.WriteFile(m.registryPath, data, 0644)
}

func (m *Manager) Install(name, pluginPath string) error {
	if err := m.loadRegistry(); err != nil {
		return err
	}

	// Verify the plugin can be loaded
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("invalid plugin binary: %w", err)
	}

	// Verify it has an Init function
	if _, err := p.Lookup("Init"); err != nil {
		return fmt.Errorf("plugin does not export Init symbol: %w", err)
	}

	// Copy plugin to plugins directory
	destPath := filepath.Join(m.pluginsDir, filepath.Base(pluginPath))
	if err := os.MkdirAll(m.pluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	data, err := os.ReadFile(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin file: %w", err)
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write plugin file: %w", err)
	}

	// Register the plugin
	m.registry.Plugins[name] = filepath.Base(pluginPath)
	return m.saveRegistry()
}

func (m *Manager) Uninstall(name string) error {
	if err := m.loadRegistry(); err != nil {
		return err
	}

	pluginPath, exists := m.registry.Plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not installed", name)
	}

	// Remove plugin file
	fullPath := filepath.Join(m.pluginsDir, pluginPath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plugin file: %w", err)
	}

	// Remove from registry
	delete(m.registry.Plugins, name)
	return m.saveRegistry()
}

func (m *Manager) List() ([]string, error) {
	if err := m.loadRegistry(); err != nil {
		return nil, err
	}

	plugins := make([]string, 0, len(m.registry.Plugins))
	for name := range m.registry.Plugins {
		plugins = append(plugins, name)
	}
	return plugins, nil
}

func (m *Manager) LoadAll() error {
	if err := m.loadRegistry(); err != nil {
		return err
	}

	for name, pluginPath := range m.registry.Plugins {
		fullPath := path.Join(m.pluginsDir, pluginPath)
		plg, err := plugin.Open(fullPath)
		if err != nil {
			fmt.Printf("Warning: failed to load plugin %s: %v\n", name, err)
			continue
		}

		initFn, err := plg.Lookup("Init")
		if err != nil {
			fmt.Printf("Warning: plugin %s does not export Init symbol\n", name)
			continue
		}

		if init, ok := initFn.(func()); ok {
			init()
			m.loaded[name] = plg
		} else {
			fmt.Printf("Warning: plugin %s Init symbol is not a function\n", name)
		}
	}

	return nil
}
