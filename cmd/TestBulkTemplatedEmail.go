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
	"github.com/spf13/cobra"
)

var (
	code      string
	recipient string
)

// TestBulkTemplatedEmailCmd represents the TestBulkTemplatedEmail command
var TestBulkTemplatedEmailCmd = &cobra.Command{
	Use:   "TestBulkTemplatedEmail",
	Short: "Sends a test bulk templated email",
	Long:  `Sends a test bulk templated email to a single recipient.`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := ses_mailer.NewSesMailer(region)
		if err != nil {
			panic(err)
		}
		err = m.SendBulkTemplatedEmail(TemplateName, FromAddress, ReplyToAddress, []mailer.Recipient{mailer.Recipient{
			Name:      "TestUser",
			Email:     recipient,
			AlbumCode: code,
		}})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(TestBulkTemplatedEmailCmd)
	TestBulkTemplatedEmailCmd.Flags().StringVarP(&TemplateName, "templateName", "t", "", "Name of the template")
	TestBulkTemplatedEmailCmd.MarkFlagRequired("templateName")
	TestBulkTemplatedEmailCmd.Flags().StringVarP(&FromAddress, "from", "f", "", "from address")
	TestBulkTemplatedEmailCmd.MarkFlagRequired("from")
	TestBulkTemplatedEmailCmd.Flags().StringVarP(&ReplyToAddress, "replyTo", "s", "", "reply to address")
	TestBulkTemplatedEmailCmd.MarkFlagRequired("replyTo")
	TestBulkTemplatedEmailCmd.Flags().StringVarP(&recipient, "destinationAddress", "d", "", "destination address")
	TestBulkTemplatedEmailCmd.MarkFlagRequired("destinationAddress")
	code = "TestCodePleaseIgnore"
}
