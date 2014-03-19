/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
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

// #cgo pkg-config: glib-2.0 libcanberra
// #include <stdlib.h>
// #include "canberra_wrapper.h"
import "C"
import "unsafe"

import (
	"dlib/dbus"
	"dlib/gio-2.0"
)

const (
	personalizationID     = "com.deepin.dde.personalization"
	gkeyCurrentSoundTheme = "current-sound-theme"
)

var personSettings = gio.NewSettings(personalizationID)

type Sound struct{}

func (s *Sound) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Sound",
		"/com/deepin/api/Sound",
		"com.deepin.api.Sound",
	}
}

// PlaySystemSound play a target event sound, such as "bell".
func (s *Sound) PlaySystemSound(event string) (err error) {
	return s.PlayThemeSound(s.getCurrentSoundTheme(), event)
}

func (s *Sound) getCurrentSoundTheme() string {
	return personSettings.GetString(gkeyCurrentSoundTheme)
}

// PlayThemeSound play a target theme's event sound.
func (s *Sound) PlayThemeSound(theme, event string) (err error) {
	go func() {
		ctheme := C.CString(theme)
		defer C.free(unsafe.Pointer(ctheme))
		cevent := C.CString(event)
		defer C.free(unsafe.Pointer(cevent))
		ret := C.canberra_play_system_sound(ctheme, cevent)
		if ret != 0 {
			logger.Errorf("play system sound failed: theme=%s, event=%s, %s",
				theme, event, C.GoString(C.ca_strerror(ret)))
		}
	}()
	return
}

// PlaySoundFile play a target sound file.
func (s *Sound) PlaySoundFile(file string) (err error) {
	go func() {
		cfile := C.CString(file)
		defer C.free(unsafe.Pointer(cfile))
		ret := C.canberra_play_sound_file(cfile)
		if ret != 0 {
			logger.Errorf("play sound file failed: %s, %s", file, C.GoString(C.ca_strerror(ret)))
		}
	}()
	return
}
