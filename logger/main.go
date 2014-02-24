/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"dlib/dbus"
	"fmt"
	_log "log"
	"os"
	"strings"
	"time"
)

const (
	_LOG_FILE          = "/var/log/deepin.log"
	_LOG_FILE_MAX_SIZE = 1024 * 1024 * 100 // 100mb
)

var (
	_LOGGER_ID uint64 = 0
	_LOGGER    *_log.Logger
)

type Logger struct {
	names map[uint64]string
}

func NewLogger() *Logger {
	logger := &Logger{}
	logger.names = make(map[uint64]string)
	return logger
}

func (logger *Logger) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Logger",
		"/com/deepin/api/Logger",
		"com.deepin.api.Logger",
	}
}

func (logger *Logger) NewLogger(name string) (id uint64, err error) {
	_LOGGER_ID++
	id = _LOGGER_ID
	logger.names[id] = name
	logger.doLog(id, "NEW", fmt.Sprintf("id=%d", id))
	return
}

func (logger *Logger) DeleteLogger(id uint64) {
	logger.doLog(id, "DELETE", fmt.Sprintf("id=%d", id))
	delete(logger.names, id)
}

func (logger *Logger) getName(id uint64) (name string) {
	if id == 0 {
		name = "<common>"
		return
	}
	name = logger.names[id]
	if len(name) == 0 {
		name = "<unknown>"
	}
	return
}

func (logger *Logger) doLog(id uint64, level, msg string) {
	now := time.Now()
	date := fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	prefix := fmt.Sprintf("%s %s: [%s] ", date, logger.getName(id), level)
	fmtMsg := prefix + msg
	fmtMsg = strings.Replace(fmtMsg, "\n", "\n"+prefix, -1)
	_LOGGER.Println(fmtMsg)
	return
}

func (logger *Logger) Debug(id uint64, msg string) {
	logger.doLog(id, "DEBUG", msg)
}

func (logger *Logger) Info(id uint64, msg string) {
	logger.doLog(id, "INFO", msg)
}

func (logger *Logger) Warning(id uint64, msg string) {
	logger.doLog(id, "WARNING", msg)
}

func (logger *Logger) Error(id uint64, msg string) {
	logger.doLog(id, "ERROR", msg)
}

func (logger *Logger) Fatal(id uint64, msg string) {
	logger.doLog(id, "FATAL", msg)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_log.Fatal(err) // TODO
		}
	}()

	// open log file
	logfile, err := os.OpenFile(_LOG_FILE, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()

	_LOGGER = _log.New(logfile, "", 0)

	logger := NewLogger()
	err = dbus.InstallOnSystem(logger)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}