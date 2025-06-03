package cmd

import (
	"email-to-pdf/internal/gmail"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	outputDir       string
	namingPattern   string
	month           string
)

var rootCmd = &cobra.Command{
	Use:   "email-to-pdf-organizer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := gmail.GetService()
		if srv != nil {
			fmt.Println("Successfully connected to Gmail API")

			// Ensure output directory exists
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				err = os.MkdirAll(outputDir, 0755)
				if err != nil {
					log.Fatalf("Unable to create output directory %s: %v", outputDir, err)
				}
			}

			query := fmt.Sprintf("has:attachment filename:pdf after:%s-01 before:%s-31", month, month)
			messages, err := gmail.GetMessages(srv, query)
			if err != nil {
				log.Fatalf("Unable to get messages: %v", err)
			}
			for _, msg := range messages {
				for _, header := range msg.Payload.Headers {
					if header.Name == "Subject" {
						fmt.Printf("Subject: %s\n", header.Value)
					}
				}
				attachments, err := gmail.GetAttachments(srv, msg)
				if err != nil {
					log.Printf("Unable to get attachments for message %s: %v", msg.Id, err)
					continue
				}
				for _, att := range attachments {
					newFilename := gmail.FormatFilename(msg, att.Filename, namingPattern)
					filePath := fmt.Sprintf("%s/%s", outputDir, newFilename)
					err := os.WriteFile(filePath, att.Data, 0644)
					if err != nil {
						log.Printf("Unable to save attachment %s: %v", newFilename, err)
					} else {
						fmt.Printf("Saved attachment: %s\n", filePath)
					}
				}
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "output", "Output directory for PDF files")
	rootCmd.PersistentFlags().StringVarP(&namingPattern, "name-pattern", "n", "{date}_{id}_{subject}_{original_filename}", "Naming pattern for PDF files")
	rootCmd.PersistentFlags().StringVarP(&month, "month", "m", time.Now().Format("2006-01"), "Month to filter emails by (YYYY-MM)")
}


