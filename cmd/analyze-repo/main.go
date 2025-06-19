package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/replyzer/analyze-repo/internal/analyzer"
	"github.com/replyzer/analyze-repo/internal/output"
	"github.com/replyzer/analyze-repo/internal/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	format    string
	outputFile string
	verbose   bool
	component string
	exclude   []string
	version   string = "dev" // Set by build process
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "analyze-repo [path]",
		Short: "Analyze local repositories to identify development environment requirements",
		Long: `A CLI tool that analyzes local repositories to detect languages, frameworks, 
version requirements, and external dependencies across single projects and monorepos.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runAnalysis,
	}

	rootCmd.Flags().StringVar(&format, "format", "yaml", "Output format (yaml|json)")
	rootCmd.Flags().StringVar(&outputFile, "output", "", "Output file path (default: stdout)")
	rootCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.Flags().StringVar(&component, "component", "", "Analyze specific component only")
	rootCmd.Flags().StringSliceVar(&exclude, "exclude", []string{}, "Exclude patterns (glob format)")

	// Add version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("replyzer version %s\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAnalysis(cmd *cobra.Command, args []string) error {
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}

	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	options := &types.AnalysisOptions{
		Format:    format,
		Output:    outputFile,
		Verbose:   verbose,
		Component: component,
		Exclude:   exclude,
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Analyzing repository at: %s\n", absPath)
	}

	result, err := analyzer.AnalyzeRepository(absPath, options)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	var outputData []byte
	switch format {
	case "json":
		outputData, err = json.MarshalIndent(result, "", "  ")
	case "yaml":
		outputData, err = yaml.Marshal(result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if outputFile != "" {
		return output.WriteToFile(outputFile, outputData)
	}

	fmt.Print(string(outputData))
	return nil
}