package main

import (
	"errors"
	"fmt"

	"github.com/thisissoon/bucket-boss/internal/storage"
	"github.com/thisissoon/bucket-boss/internal/storage/s3"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	noBucketProvider = errors.New("Please configure a bucket provider")
	force            *bool
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

func confirm(msg string) bool {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err == nil
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
		return noBucketProvider
	}
	keys, err := store.List(ext)
	if err != nil {
		fmt.Println(err)
		return err
	}
	total := len(keys)
	if total == 0 {
		fmt.Printf("No files to delete matching extension: %v\n", ext)
		return nil
	}
	confirmText := fmt.Sprintf("Are you sure you want to delete %d %v files", total, ext)
	ok := confirm(confirmText)
	if !ok {
		fmt.Println("exiting")
		return nil
	}

	err = store.DeleteMulti(keys)
	if err != nil {
		return err
	}
	return nil
}
