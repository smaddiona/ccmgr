package cmd

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "ccmgr",
	Short: "Claude Code configuration manager",
	Long:  "Manage multiple environment configurations for Claude Code CLI.\nSwitch between providers like Z.AI and default Anthropic with ease.",
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new configuration profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreate()
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to a different configuration profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSwitch()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList()
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [label]",
	Short: "Delete a configuration profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		label := ""
		if len(args) > 0 {
			label = args[0]
		}
		return runDelete(label)
	},
}

func fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}

func printBanner() {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DC4E4")).Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9399B2"))

	banner := `
 ██████╗ ██████╗███╗   ███╗ ██████╗ ██████╗ 
██╔════╝██╔════╝████╗ ████║██╔════╝ ██╔══██╗
██║     ██║     ██╔████╔██║██║  ███╗██████╔╝
██║     ██║     ██║╚██╔╝██║██║   ██║██╔══██╗
╚██████╗╚██████╗██║ ╚═╝ ██║╚██████╔╝██║  ██║
 ╚═════╝ ╚═════╝╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═╝`

	fmt.Println(titleStyle.Render(banner))
	fmt.Println(subtitleStyle.Render("   Developed by Simone Maddiona"))
	fmt.Println()
}
