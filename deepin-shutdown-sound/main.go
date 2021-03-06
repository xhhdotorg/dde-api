/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"os/signal"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/sound"
)

var logger = log.NewLogger("api/shutdown-sound")

func main() {
	logger.Info("[DEEPIN SHUTDOWN SOUND] play shutdown sound")
	handleSignal()

	canPlay, theme, event, err := soundutils.GetShutdownSound()
	if err != nil {
		logger.Warning("[DEEPIN SHUTDOWN SOUND] get shutdown sound info failed:", err)
		return
	}
	logger.Info("[DEEPIN SHUTDOWN SOUND] can play:", canPlay, theme, event)

	if !canPlay {
		return
	}

	err = doPlayShutdwonSound(theme, event)
	if err != nil {
		logger.Error("[DEEPIN SHUTDOWN SOUND] play shutdown sound failed:", theme, event, err)
	}
}

func handleSignal() {
	var sigs = make(chan os.Signal, 2)
	signal.Notify(sigs, os.Kill, os.Interrupt)
	go func() {
		sig := <-sigs
		switch sig {
		case os.Kill, os.Interrupt:
			// Nothing to do
			logger.Info("[DEEPIN SHUTDOWN SOUND] receive signal:", sig.String())
		}
	}()
}

func doPlayShutdwonSound(theme, event string) error {
	logger.Info("[DEEPIN SHUTDOWN SOUND] do play:", theme, event)
	err := sound.PlayThemeSound(theme, event, "", "alsa")
	if err != nil {
		logger.Error("[DEEPIN SHUTDOWN SOUND] do play failed:", theme, event, err)
		return err
	}
	return nil
}
