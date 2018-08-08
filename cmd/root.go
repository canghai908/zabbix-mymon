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
	"fmt"
	"github.com/canghai908/zabbix-mymon/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	cfgFile         string
	Debug           bool
	Step            int
	Mysql_username  string
	Mysql_password  string
	Mysql_host      string
	Mysql_port      int
	Zabbix_host     string
	Zabbix_port     int
	Zabbix_hostname string
)
var ConfigFile string

var rootCmd = &cobra.Command{
	Use:   "mymon",
	Short: "Zabbix mysql monitoring tool",
	Long: `Zabbix mysql database monitoring tool. For example:

mymon daemon --config=./mymon.json`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /.mymon.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("mymon")

	}
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}
	ConfigFile = viper.ConfigFileUsed()
	Debug = viper.GetBool("debug")
	Step = viper.GetInt("step")
	Mysql_username = viper.GetString("mysql.username")
	Mysql_password = viper.GetString("mysql.password")
	Mysql_host = viper.GetString("mysql.host")
	Mysql_port = viper.GetInt("mysql.port")
	Zabbix_host = viper.GetString("zabbix.server")
	Zabbix_port = viper.GetInt("zabbix.port")
	Zabbix_hostname = viper.GetString("zabbix.hostname")

	if Debug {
		utils.InitLog("debug")
	} else {
		utils.InitLog("info")
	}
}
