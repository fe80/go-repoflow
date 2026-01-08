package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"repoflow/internal/factory"
	"repoflow/pkg/client"
)

// RepositoryManager handles the state and configuration for workspace commands
type RepositoryManager struct {
	*factory.Utils
	childRepositoryIds                []string
	workspace                         string
	packageType                       string
	uploadLocalRepositoryId           string
	remoteRepositoryUrl               string
	remoteRepositoryUsername          string
	remoteRepositoryPassword          string
	isRemoteCacheEnabled              bool
	fileCacheTimeTillRevalidation     *int
	metadataCacheTimeTillRevalidation *int
}

// RepositoryCmd initializes the parent command and its subcommands
func RepositoryCmd(u *factory.Utils) *cobra.Command {
	m := &RepositoryManager{Utils: u}

	// Main repository command
	var repositoryCmd = &cobra.Command{
		Use:   "repository",
		Short: "Manage RepoFlow repository",
	}

	repositoryCmd.PersistentFlags().StringVarP(
		&m.workspace, "workspace", "w", "", "Repository workspace to work (id or name)",
	)
	repositoryCmd.MarkPersistentFlagRequired("workspace")

	// List sub-command
	var listCmd = &cobra.Command{
		Use:          "list",
		Short:        "List all repositories",
		RunE:         m.repositoryList,
		SilenceUsage: true,
	}

	// Get sub-command
	var getCmd = &cobra.Command{
		Use:          "get [name]",
		Short:        "Get repository metadata (ID or name)",
		Args:         cobra.ExactArgs(1),
		RunE:         m.repositoryGet,
		SilenceUsage: true,
	}

	// Delete sub-command
	var deleteCmd = &cobra.Command{
		Use:          "delete [name]",
		Short:        "Delete a repository (ID or name)",
		Args:         cobra.ExactArgs(1),
		RunE:         m.repositoryDelete,
		SilenceUsage: true,
	}

	// Delete sub-command
	var deleteContentCmd = &cobra.Command{
		Use:          "prune [name]",
		Short:        "Delete a repository content (ID or name)",
		Args:         cobra.ExactArgs(1),
		RunE:         m.repositoryDeleteContent,
		SilenceUsage: true,
	}

	// Create repository
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new repository",
	}

	createCmd.PersistentFlags().StringVarP(&m.packageType, "type", "t", "", "Package type stored by the repository.")
	createCmd.MarkPersistentFlagRequired("type")

	// Create local repository
	var createLocalCmd = &cobra.Command{
		Use:   "local [name]",
		Short: "Create a new local repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.repositoryCreate(cmd, args, "local")
		},
		SilenceUsage: true,
	}

	// Create remote repository
	var (
		fileCacheTimeTillRevalidation     int
		metadataCacheTimeTillRevalidation int
	)
	var createRemoteCmd = &cobra.Command{
		Use:   "remote [name]",
		Short: "Create a new remote repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.repositoryCreate(cmd, args, "remote")
		},
		SilenceUsage: true,
	}
	createRemoteCmd.Flags().StringVarP(
		&m.remoteRepositoryUrl, "remote-url", "r", "", "URL of the remote repository",
	)
	createRemoteCmd.Flags().StringVarP(
		&m.remoteRepositoryUsername, "remote-username", "u", "", "URL of the remote repository.",
	)
	createRemoteCmd.Flags().StringVarP(
		&m.remoteRepositoryPassword, "remote-password", "p", "", "Username for the remote repository.",
	)
	createRemoteCmd.Flags().BoolVarP(
		&m.isRemoteCacheEnabled, "cache", "c", false, "Whether caching is enabled.",
	)
	createRemoteCmd.Flags().IntVar(
		&fileCacheTimeTillRevalidation, "file-cache-ttr", -1,
		"Milliseconds before cached files require revalidation (-1 for indefinite caching).",
	)
	createRemoteCmd.Flags().IntVar(
		&metadataCacheTimeTillRevalidation, "metadata-cache-ttr", -1,
		"Milliseconds before cached metadata requires revalidation (-1 for indefinite caching).",
	)

	createRemoteCmd.MarkFlagRequired("remote-url")
	createRemoteCmd.MarkFlagRequired("remote-username")

	createRemoteCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file-cache-ttr") && fileCacheTimeTillRevalidation >= 0 {
			m.fileCacheTimeTillRevalidation = &fileCacheTimeTillRevalidation
		} else {
			m.fileCacheTimeTillRevalidation = nil
		}
		if cmd.Flags().Changed("metadata-cache-ttr") && metadataCacheTimeTillRevalidation >= 0 {
			m.metadataCacheTimeTillRevalidation = &metadataCacheTimeTillRevalidation
		} else {
			m.metadataCacheTimeTillRevalidation = nil
		}
		return nil
	}

	var createVirtualCmd = &cobra.Command{
		Use:   "virtual [name]",
		Short: "Create a new virtual repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.repositoryCreate(cmd, args, "virtual")
		},
		SilenceUsage: true,
	}

	createVirtualCmd.Flags().StringSliceVarP(
		&m.childRepositoryIds, "child-repository", "r", []string{},
		"IDs of repositories included in the virtual repository.",
	)
	createVirtualCmd.Flags().StringVar(
		&m.uploadLocalRepositoryId, "local-repository", "",
		"ID of a local repository where uploads will be stored (must also be in childRepositoryIds).",
	)

	createVirtualCmd.MarkFlagRequired("child-repository")

	createCmd.AddCommand(createLocalCmd, createRemoteCmd, createVirtualCmd)

	// Register sub-commands
	repositoryCmd.AddCommand(
		listCmd, createCmd, getCmd, deleteCmd, deleteContentCmd,
	)

	return repositoryCmd
}

// --- Runners Implementation ---
func (m *RepositoryManager) repositoryList(cmd *cobra.Command, args []string) error {
	resp, err := m.GetAPIClient().ListRepositories(m.workspace)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := m.HandleResponse(resp); err != nil {
		return err
	}

	var repositories []client.Repositories
	return factory.HandleOutput(m.Utils, resp, &repositories)
}

func (m *RepositoryManager) repositoryGet(cmd *cobra.Command, args []string) error {
	resp, err := m.GetAPIClient().GetRepository(m.workspace, args[0])
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := m.HandleResponse(resp); err != nil {
		return err
	}

	var repository client.Repository
	return factory.HandleOutput(m.Utils, resp, &repository)
}

func (m *RepositoryManager) repositoryDelete(cmd *cobra.Command, args []string) error {
	resp, err := m.GetAPIClient().DeleteRepository(m.workspace, args[0])
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := m.HandleResponse(resp); err != nil {
		return err
	}

	if m.Output == "text" || m.Output == "" {
		fmt.Printf("Successfully deleted repository '%s' workspace '%s'\n", args[0], m.workspace)
		return nil
	}

	var deletedRepository client.RepostotryDelete
	return factory.HandleOutput(m.Utils, resp, &deletedRepository)
}

func (m *RepositoryManager) repositoryCreate(cmd *cobra.Command, args []string, store string) error {
	name := args[0]

	var opts any

	switch store {
	case "local":
		opts = client.RepositoryOptions{
			Name:        name,
			PackageType: m.packageType,
		}

	case "remote":
		opts = client.RepositoryRemoteOptions{
			Name:                              name,
			PackageType:                       m.packageType,
			RemoteRepositoryUrl:               m.remoteRepositoryUrl,
			RemoteRepositoryUsername:          m.remoteRepositoryUsername,
			RemoteRepositoryPassword:          m.remoteRepositoryPassword,
			IsRemoteCacheEnabled:              m.isRemoteCacheEnabled,
			FileCacheTimeTillRevalidation:     m.fileCacheTimeTillRevalidation,
			MetadataCacheTimeTillRevalidation: m.metadataCacheTimeTillRevalidation,
		}

	case "virtual":
		opts = client.RepositoryVirtualOptions{
			Name:                    name,
			PackageType:             m.packageType,
			ChildRepositoryIds:      m.childRepositoryIds,
			UploadLocalRepositoryId: &m.uploadLocalRepositoryId,
		}

	default:
		return fmt.Errorf("Unsuported store store type: %s", store)
	}

	resp, err := m.GetAPIClient().CreateRepository(m.workspace, store, opts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := m.HandleResponse(resp); err != nil {
		return err
	}

	if m.Output == "text" || m.Output == "" {
		fmt.Printf("Successfully created repository '%s' on workspace '%s'.\n", name, m.workspace)
		return nil
	}

	switch store {
	case "local":
		var createdRepository client.RepositoryOptions
		return factory.HandleOutput(m.Utils, resp, createdRepository)

	case "remote":
		var createdRepository client.RepositoryRemoteOptions
		return factory.HandleOutput(m.Utils, resp, createdRepository)

	case "virtual":
		var createdRepository client.RepositoryVirtualOptions
		return factory.HandleOutput(m.Utils, resp, createdRepository)

	}

	return nil
}

func (m *RepositoryManager) repositoryDeleteContent(cmd *cobra.Command, args []string) error {
	name := args[0]

	resp, err := m.GetAPIClient().DeleteRepositoryContent(m.workspace, name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := m.HandleResponse(resp); err != nil {
		return err
	}

	if m.Output == "text" || m.Output == "" {
		fmt.Printf("Successfully delete content of repository '%s' workspace '%s'\n", args[0], m.workspace)
		return nil
	}

	if m.Output == "text" || m.Output == "" {
		fmt.Printf("Successfully deleted repository '%s' workspace '%s'\n", name, m.workspace)
		return nil
	}

	var deletedRepository client.RepostotryDelete
	return factory.HandleOutput(m.Utils, resp, &deletedRepository)
}
