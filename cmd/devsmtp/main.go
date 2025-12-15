package main

import (
	"fmt"
	"os"

	"github.com/lawnchairsociety/devsmtp/internal/config"
	"github.com/lawnchairsociety/devsmtp/internal/database"
	"github.com/lawnchairsociety/devsmtp/internal/smtp"
	"github.com/lawnchairsociety/devsmtp/internal/tui"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "devsmtp",
	Short: "A developer SMTP server for testing email functionality",
	Long: `DevSmtp is a lightweight SMTP server that captures all emails
sent to it and stores them in a SQLite database. It provides a
terminal-based UI for viewing and managing messages.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := database.New(cfg.Database.Path)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		// Create logger for SMTP server
		logger := smtp.NewLogger(1000)

		// Start SMTP server in background
		server := smtp.NewServer(cfg, db, logger)
		go func() {
			if err := server.ListenAndServe(); err != nil {
				logger.Error("SMTP server error: %v", err)
			}
		}()

		// Run TUI in foreground with log channel
		return tui.Run(db, cfg, logger.Channel())
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is ./devsmtp.yaml)")
	rootCmd.Flags().String("host", "0.0.0.0", "SMTP server bind address")
	rootCmd.Flags().Int("port", 587, "SMTP server port")
	rootCmd.Flags().String("db", "./devsmtp.db", "SQLite database path")
	rootCmd.Flags().Bool("auth-required", false, "Require SMTP authentication")
	rootCmd.Flags().String("auth-user", "", "Username for SMTP AUTH")
	rootCmd.Flags().String("auth-pass", "", "Password for SMTP AUTH")
	rootCmd.Flags().String("tls-cert", "", "Path to TLS certificate")
	rootCmd.Flags().String("tls-key", "", "Path to TLS private key")
}

func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile, rootCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
}
