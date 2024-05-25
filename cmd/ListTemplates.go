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
	"fmt"
	"github.com/spf13/cobra"
)

// ListTemplatesCmd represents the ListTemplates command
var ListTemplatesCmd = &cobra.Command{
	Use:   "ListTemplates",
	Short: "List templates in SES",
	Long:  `Returns a list of email templates in SES that can be used to send emails`,
	Run: func(cmd *cobra.Command, args []string) {
		m, err := ses_mailer.NewSesMailer(region)
		if err != nil {
			panic(err)
		}
		o, err := m.ListTemplates()
		if err != nil {
			panic(err)
		}
		for _, t := range o {
			fmt.Printf("Name: %s\t\tCreated: %s\n", *t.Name, t.CreatedTimestamp.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(ListTemplatesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ListTemplatesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ListTemplatesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
