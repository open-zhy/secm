package cmd

import (
	"github.com/open-zhy/secm/pkg/plugin"
	"github.com/open-zhy/secm/pkg/screen"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage secm plugins",
	Long:  `Plugin management commands for installing, uninstalling, and listing plugins.`,
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [name] [path]",
	Short: "Install a plugin",
	Long:  `Install a plugin from a .so file.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugin.NewManager(pluginsDir)
		return manager.Install(args[0], args[1])
	},
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall [name]",
	Short: "Uninstall a plugin",
	Long:  `Uninstall an installed plugin.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugin.NewManager(pluginsDir)
		return manager.Uninstall(args[0])
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Long:  `List all installed plugins.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugin.NewManager(pluginsDir)
		plugins, err := manager.List()
		if err != nil {
			return err
		}

		if len(plugins) == 0 {
			screen.Println("No plugins installed")
			return nil
		}

		screen.Println("Installed plugins:")
		for _, name := range plugins {
			screen.Printf("  - %s\n", name)
		}
		return nil
	},
}

func init() {
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
	pluginCmd.AddCommand(pluginListCmd)
	rootCmd.AddCommand(pluginCmd)
}
