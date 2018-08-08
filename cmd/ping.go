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
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/canghai908/zabbix-mymon/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"log"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Connected line checker",
	Long: `Connected line checker. For example:

mymon ping --config=./mymon.json`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Using config file:", ConfigFile, " successfully!")
		//Aes解密密码
		decodeBytes_password, err := base64.StdEncoding.DecodeString(Mysql_password)
		if err != nil {
			fmt.Println(err)
			return
		}
		dec_password, err := utils.AesDecode([]byte(decodeBytes_password))
		if err != nil {
			fmt.Println(err)
			return
		}
		//连接mysql
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", Mysql_username, string(dec_password),
			Mysql_host, Mysql_port) + "?clientFoundRows=false&timeout=5s&charset=utf8&collation=utf8_general_ci"
		db, err := sql.Open("mysql", dsn)
		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(20)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()
		t := db.Ping()
		if t != nil {
			fmt.Println("0")
			return
		}
		fmt.Println("1")
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
