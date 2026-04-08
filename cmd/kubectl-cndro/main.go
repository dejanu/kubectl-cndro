package main

import (
	"fmt"
	"os"

	"github.com/dejanu/cndro/internal/cndro"
	"github.com/spf13/cobra"
)

func ticketsCmd() *cobra.Command {
	var openBrowser bool
	c := &cobra.Command{
		Use:   "tickets",
		Short: "Print the tickets page link and pricing table",
		Long: "Print the tickets URL (as a clickable link when your terminal supports it) and the pricing table.\n" +
			"Use --open to launch the page in your default browser.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if openBrowser {
				if err := cndro.OpenTicketsURL(); err != nil {
					return fmt.Errorf("open browser: %w", err)
				}
			}
			w := cmd.OutOrStdout()
			if err := cndro.WriteTicketsURLLine(w); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
			return cndro.WritePricing(w)
		},
	}
	c.Flags().BoolVarP(&openBrowser, "open", "o", false, "Open the tickets page in the default web browser")
	return c
}

func main() {
	root := &cobra.Command{
		Use:   "cndro",
		Short: "Cloud Native Days Romania conference helper",
		Long:  "kubectl plugin for Cloud Native Days Romania (CNDRO): schedules and tickets.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	root.AddCommand(
		&cobra.Command{
			Use:   "day1",
			Short: "Print the Day 1 schedule",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cndro.WriteDay1(cmd.OutOrStdout())
			},
		},
		&cobra.Command{
			Use:   "day2",
			Short: "Print the Day 2 schedule",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cndro.WriteDay2(cmd.OutOrStdout())
			},
		},
		ticketsCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
