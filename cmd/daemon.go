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
	log "github.com/Sirupsen/logrus"
	. "github.com/canghai908/go-zabbix"
	"github.com/canghai908/zabbix-mymon/utils"
	"github.com/spf13/cobra"
	"time"
)

var SlaveStatusToSend = []string{
	"Exec_Master_Log_Pos",
	"Read_Master_Log_Pos",
	"Relay_Log_Pos",
	"Seconds_Behind_Master",
	"Slave_IO_Running",
	"Slave_SQL_Running",
}

var Db *sql.DB

func mysqlState(db *sql.DB, host, strsql string) ([]*Metric, error) {
	rows, err := db.Query(strsql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	met := make([]*Metric, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		var ma *Metric
		ma = NewMetric(host, "mysql."+string(values[0]))
		ma.SetValue(string(values[1]))
		met = append(met, ma)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return met, nil
}

func GlobalStatus(db *sql.DB, host string) ([]*Metric, error) {
	return mysqlState(db, host, "SHOW /*!50001 GLOBAL */ STATUS")
}

func GlobalVariables(db *sql.DB, host string) ([]*Metric, error) {
	return mysqlState(db, host, "SHOW /*!50001 GLOBAL */ VARIABLES")
}

func SlaveStatus(db *sql.DB, host string) ([]*Metric, error) {
	isSlave := NewMetric(Zabbix_hostname, "mysql.Is_slave")
	rows, err := db.Query("SHOW SLAVE STATUS")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cols, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	isSlave.SetValue(0)
	met := make([]*Metric, 0)
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			fmt.Println("Error")
			return nil, err
		}
		isSlave.SetValue(1)
		for j, v := range values {
			key := cols[j]
			var ma *Metric
			ma = NewMetric(Zabbix_hostname, "mysql."+key)
			switch key {
			case "Slave_SQL_Running", "Slave_IO_Running":
				ma.SetValue(0)
				if string(v) == "Yes" {
					ma.SetValue(1)
				}
			default:
				ma.SetValue(string(v))
			}
			met = append(met, ma)
		}
	}
	if isSlave.Value == 0 {
		data := make([]*Metric, len(SlaveStatusToSend))
		for i, s := range SlaveStatusToSend {
			data[i] = NewMetric(Zabbix_hostname, "mysql."+s)
			switch s {
			case "mysql." + "Slave_SQL_Running", "mysql." + "Slave_IO_Running":
				data[i].SetValue(0)
			default:
				data[i].SetValue(0)
			}
		}
		met = append(met, data...)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return append(met, isSlave), nil
}

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Running as a daemon",
	Long: `Running as a daemon. For example:

mymon deamon --config=./mymon.json`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Using config file:", ConfigFile, " successfully!")
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
		Db, err := sql.Open("mysql", dsn)
		Db.SetMaxOpenConns(20)
		Db.SetMaxIdleConns(20)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			data := make([]*Metric, 0)
			//slaveStatus
			met, err := SlaveStatus(Db, Zabbix_hostname)
			if err != nil {
				return
			}
			data = append(data, met...)

			//globalStatus
			globalStatus, err := GlobalStatus(Db, Zabbix_hostname)
			if err != nil {
				return
			}
			data = append(data, globalStatus...)

			//globalVars
			globalVars, err := GlobalVariables(Db, Zabbix_hostname)
			if err != nil {
				return
			}
			data = append(data, globalVars...)
			log.Debugf("Send to %s, size: %d", Zabbix_hostname, len(data))
			// Send packet to zabbix
			packet := NewPacket(data)
			z := NewSender(Zabbix_host, Zabbix_port)
			z.Send(packet)
			time.Sleep(time.Duration(Step) * time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
