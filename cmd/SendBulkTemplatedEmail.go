/*
Copyright © 2024 Morgan MacKechnie morgan@kernelwit.ch

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	ses_mailer "AlbumCodeMailer/lib/ses"
	"AlbumCodeMailer/mailer"
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	codeFile string
	dryRun   bool
)

// SendBulkTemplatedEmailCmd represents the SendBulkTemplatedEmail command
var SendBulkTemplatedEmailCmd = &cobra.Command{
	Use:   "SendBulkTemplatedEmail",
	Short: "Send a bulk templated email",
	Long: `After uploading and testing your template, this command renders
the template and sends in bulk to all recipients in the supplied input file`,
	Run: func(cmd *cobra.Command, args []string) {
		recipsFile, err := os.Open(RecipientFile)
		if err != nil {
			panic(err)
		}
		defer recipsFile.Close()
		codesFile, err := os.Open(codeFile)
		if err != nil {
			panic(err)
		}
		defer codesFile.Close()
		m, err := ses_mailer.NewSesMailer(region)
		if err != nil {
			panic(err)
		}
		cs := bufio.NewScanner(codesFile)
		scanner := bufio.NewScanner(recipsFile)
		// skip the first line as it's almost certainly a header
		scanner.Scan()
		var recipients []mailer.Recipient
		for scanner.Scan() {
			address := strings.Split(scanner.Text(), ",")[1]
			_, err := mail.ParseAddress(address)
			if err != nil {
				log.Fatalf("Invalid email address: %s", address)
			}
			if ok := cs.Scan(); !ok {
				log.Fatalf("you seem to have more recipients than codes?")
			}
			recipients = append(recipients, mailer.Recipient{
				Email:     address,
				AlbumCode: cs.Text(),
			})
		}
		if dryRun {
			for _, recipient := range recipients {
				fmt.Printf("Sending code %s to %s\n", recipient.AlbumCode, recipient.Email)
			}
			return
		}
		for i := 0; i < len(recipients); i += 50 {
			end := i + 50
			if end > len(recipients) {
				end = len(recipients)
			}
			err := m.SendBulkTemplatedEmail(TemplateName, FromAddress, ReplyToAddress, recipients[i:end])
			if err != nil {
				panic(err)
			}
			time.Sleep(30 * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(SendBulkTemplatedEmailCmd)
	SendBulkTemplatedEmailCmd.Flags().StringVarP(&RecipientFile, "recipients", "i", "",
		"path to file containing list of recipients")
	SendBulkTemplatedEmailCmd.MarkFlagRequired("recipients")
	SendBulkTemplatedEmailCmd.Flags().StringVarP(&TemplateName, "templateName", "t", "", "Name of the template")
	SendBulkTemplatedEmailCmd.MarkFlagRequired("templateName")
	SendBulkTemplatedEmailCmd.Flags().StringVarP(&FromAddress, "from", "f", "", "from address")
	SendBulkTemplatedEmailCmd.MarkFlagRequired("from")
	SendBulkTemplatedEmailCmd.Flags().StringVarP(&ReplyToAddress, "replyTo", "s", "", "reply to address")
	SendBulkTemplatedEmailCmd.MarkFlagRequired("replyTo")
	SendBulkTemplatedEmailCmd.Flags().StringVarP(&codeFile, "codes", "c", "", "path to file containing one redemption code per line")
	SendBulkTemplatedEmailCmd.MarkFlagRequired("codes")
	SendBulkTemplatedEmailCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Dry run - will print each recipient email and code, but not send emails")
}
