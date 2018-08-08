// Copyright © 2018 canghai908 <lovecanghai@gmail.com.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/canghai908/zabbix-mymon/utils"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt passwords in AES mode",
	Long: `Encrypt passwords in AES mode. For example:

mymon enc password`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			fmt.Println("add password")
			return
		}
		passwd := []byte(args[0])
		en_password, err := utils.AesEncode(passwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		encodeString := base64.StdEncoding.EncodeToString(en_password)
		fmt.Println(encodeString)
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
}
