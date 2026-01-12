package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fe80/go-repoflow/internal/factory"
	"github.com/fe80/go-repoflow/pkg/repoflow"
)

// WorkspaceManager handles the state and configuration for workspace commands
type WorkspaceManager struct {
	*factory.Utils
	packageLimit   *int
	bandwidthLimit *int
	storageLimit   *int
	comments       *string
}

// WorkspaceCmd initializes the parent command and its subcommands
func WorkspaceCmd(u *factory.Utils) *cobra.Command {
	m := &WorkspaceManager{Utils: u}

	// Main workspace command
	var workspaceCmd = &cobra.Command{
		Use:   "workspace",
		Short: "Manage RepoFlow workspaces",
	}

	// List sub-command
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all workspaces",
		RunE:         m.workspaceList,
		SilenceUsage: true,
	}

	// Get sub-command
	var getCmd = &cobra.Command{
		Use:          "get [name]",
		Short:        "Get workspace metadata (workspace ID or name)",
		Args:         cobra.ExactArgs(1),
		RunE:         m.workspaceGet,
		SilenceUsage: true,
	}

	// Delete sub-command
	var deleteCmd = &cobra.Command{
		Use:          "delete [name]",
		Short:        "Delete a workspace (workspace ID or name)",
		Args:         cobra.ExactArgs(1),
		RunE:         m.workspaceDelete,
		SilenceUsage: true,
	}

	// Create sub-command with flags
	var (
		pkgLim   int
		bwLim    int
		stLim    int
		comments string
	)
	var createCmd = &cobra.Command{
		Use:          "create [name]",
		Short:        "Create a new workspace",
		Args:         cobra.ExactArgs(1),
		RunE:         m.workspaceCreate,
		SilenceUsage: true,
	}
	createCmd.Flags().IntVarP(&pkgLim, "package-limit", "p", 0, "Maximum packages allowed")
	createCmd.Flags().IntVarP(&bwLim, "bandwidth-limit", "b", 0, "Bandwidth limit in bytes")
	createCmd.Flags().IntVarP(&stLim, "storage-limit", "s", 0, "Storage limit in bytes")
	createCmd.Flags().StringVarP(&comments, "comments", "c", "", "Optional notes about the workspace")
	createCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("package-limit") {
			m.packageLimit = &pkgLim
		}
		if cmd.Flags().Changed("bandwidth-limit") {
			m.bandwidthLimit = &bwLim
		}
		if cmd.Flags().Changed("storage-limit") {
			m.storageLimit = &stLim
		}
		if cmd.Flags().Changed("comments") {
			m.comments = &comments
		}
		return nil
	}

	// Register sub-commands
	workspaceCmd.AddCommand(listCmd, createCmd, getCmd, deleteCmd)

	return workspaceCmd
}

// --- Runners Implementation ---

func (m *WorkspaceManager) workspaceList(cmd *cobra.Command, args []string) error {
	data, err := m.GetAPIClient().ListWorkspaces()
	if err != nil {
		return err
	}
	return factory.HandleOutput(m.Utils, data)
}

func (m *WorkspaceManager) workspaceGet(cmd *cobra.Command, args []string) error {
	data, err := m.GetAPIClient().GetWorkspace(args[0])
	if err != nil {
		return err
	}

	return factory.HandleOutput(m.Utils, data)
}

func (m *WorkspaceManager) workspaceDelete(cmd *cobra.Command, args []string) error {
	data, err := m.GetAPIClient().DeleteWorkspace(args[0])
	if err != nil {
		return err
	}

	if m.Output == "text" || m.Output == "" {
		fmt.Printf("Successfully deleted workspace '%s'\n", args[0])
		return nil
	}
	return factory.HandleOutput(m.Utils, data)
}

func (m *WorkspaceManager) workspaceCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	opts := repoflow.WorkspaceOptions{
		Name:           name,
		PackageLimit:   m.packageLimit,
		BandwidthLimit: m.bandwidthLimit,
		StorageLimit:   m.storageLimit,
		Comments:       m.comments,
	}

	data, err := m.GetAPIClient().CreateWorkspace(opts)

	if err != nil {
		return err
	}

	if m.Output == "text" || m.Output == "" {
		fmt.Printf("Successfully created workspace '%s'.\n", name)
		return nil
	}

	return factory.HandleOutput(m.Utils, data)
}
