package main

import (
	"fmt"
	"github.com/thisissoon/bucket-boss/internal/storage"
	"github.com/thisissoon/bucket-boss/internal/storage/s3"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func purgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "command to delete files from bucket",
		RunE:  purgeRun,
	}
	cmd.AddCommand(purgeByExtentionCmd())
	return cmd
}

func purgeRun(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func purgeByExtentionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ext",
		Short:   "delete files by their extension",
		Example: "ext csv",
		RunE:    purgeByExtentionRun,
		Args:    cobra.MaximumNArgs(1),
	}
}

func purgeByExtentionRun(_ *cobra.Command, args []string) error {
	var ext string
	if len(args) == 1 {
		ext = args[0]
	}
	var store storage.Storer
	if cfg.AWS.Enabled {
		s3, err := s3.NewS3(cfg.AWS)
		if err != nil {
			return err
		}
		store = s3
	} else {
		return fmt.Errorf("Please configure a bucket provider")
	}
	keys, err := store.List(ext)
	if err != nil {
		return err
	}
	total := len(keys)
	if total == 0 {
		fmt.Printf("No files to delete matching extension: %v\n", ext)
		return nil
	}
	confirmText := fmt.Sprintf("Are you sure you want to delete %d %v files", len(keys), ext)
	prompt := promptui.Prompt{
		Label:     confirmText,
		IsConfirm: true,
	}
	_, err = prompt.Run()
	if err != nil {
		fmt.Printf("exiting")
		return nil
	}

	fmt.Printf("deleting files")
	return nil
}
