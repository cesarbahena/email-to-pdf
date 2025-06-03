package cmd

import (
	"email-to-pdf/internal/gmail"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
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
			messages, err := gmail.GetMessages(srv, "has:attachment filename:pdf")
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
					filePath := fmt.Sprintf("tmp_attachments/%s", att.Filename)
					err := os.WriteFile(filePath, att.Data, 0644)
					if err != nil {
						log.Printf("Unable to save attachment %s: %v", att.Filename, err)
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
