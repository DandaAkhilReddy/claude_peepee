package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/DandaAkhilReddy/claude_pp/internal/store"
	"github.com/spf13/cobra"
)

var (
	recallTags  []string
	recallLimit int
	recallLocal bool
)

var recallCmd = &cobra.Command{
	Use:   "recall [query]",
	Short: "Search for stored facts",
	Long:  `Search for previously stored facts using keywords or tags.`,
	RunE:  runRecall,
}

func init() {
	recallCmd.Flags().StringSliceVarP(&recallTags, "tags", "t", nil, "Filter by tags")
	recallCmd.Flags().IntVarP(&recallLimit, "limit", "n", 20, "Maximum number of results")
	recallCmd.Flags().BoolVarP(&recallLocal, "local", "l", false, "Only show facts from current directory")
}

func runRecall(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	query := strings.Join(args, " ")

	var sourceDir string
	if recallLocal {
		sourceDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	facts, err := s.QueryFacts(query, recallTags, sourceDir, recallLimit)
	if err != nil {
		return fmt.Errorf("failed to query facts: %w", err)
	}

	if len(facts) == 0 {
		fmt.Println("No matching facts found.")
		return nil
	}

	for _, f := range facts {
		fmt.Printf("\n#%d [%s]\n", f.ID, f.CreatedAt.Format("2006-01-02 15:04"))
		if len(f.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(f.Tags, ", "))
		}
		if f.SourceDir != "" {
			fmt.Printf("Dir: %s\n", f.SourceDir)
		}
		fmt.Printf("%s\n", f.Content)
	}

	return nil
}
