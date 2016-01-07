package ituneslib

import (
    . "ml/trace"
    "unsafe"
    "ml/os2"
    "path/filepath"
)

type KeybagSyncType int

const (
    Keybag_Buy          = KeybagSyncType(1)
    Keybag_Refetch      = KeybagSyncType(1)
    Keybag_Default      = KeybagSyncType(2)
    Keybag_Upgrade      = KeybagSyncType(5)
    Keybag_Authorize    = KeybagSyncType(8)
    Keybag_Update       = KeybagSyncType(11)
    Keybag_LoginiOS     = KeybagSyncType(0x135)
)

type KeybagSession struct {
    session         uintptr
    uniqueDeviceID  []byte
}

func NewKeybagSession(uniqueDeviceID []byte) *KeybagSession {
    session := &KeybagSession{}
    session.initialize(uniqueDeviceID)

    return session
}

func (self *KeybagSession) initialize(uniqueDeviceID []byte) {
    self.uniqueDeviceID = uniqueDeviceID

    // scinfo := String(`C:\ProgramData\Apple Computer\iTunes\SC Info`).Encode(CP_UTF8)
    // scinfo := []byte(`C:\ProgramData\Apple Computer\iTunes\SC Info`)

    scinfo := []byte(filepath.Join(os2.ExecutablePath(), "SC Info"))

    var udid *FairPlayHWInfo = nil

    if uniqueDeviceID != nil {
        udid = &FairPlayHWInfo{
            Length: int32(len(uniqueDeviceID)),
        }

        copy(udid.Id[:], uniqueDeviceID)
        udid.Length = 6
    }

    self.Close()

    st, _, _ := itunes.KbsyncCreateSession.Call(
                    uintptr(unsafe.Pointer(&self.session)),
                    uintptr(unsafe.Pointer(udid)),
                    0,
                    uintptr(unsafe.Pointer(&scinfo[0])),
                )

    if int32(st) != 0 {
        Raise(newiTunesHelperErrorf("KeybagSession.initialize failed: %X", uint32(st)))
    }
}

func (self *KeybagSession) Close() {
    if self.session == 0 {
        return
    }

    itunes.KbsyncCloseSession.Call(self.session)
    self.session = 0
}

func (self *KeybagSession) validate() int {
    status, _, _ := itunes.KbsyncValidate.Call(self.session)
    return getStatus(status)
}

func (self *KeybagSession) GetData(dsid int64, syncType KeybagSyncType) []byte {
    var buf *byte
    var size int

    var status int

    switch unsafe.Sizeof(uintptr(0)) {
        case 4:
            st, _, _ := itunes.KbsyncGetData.Call(
                            self.session,
                            uintptr(dsid & 0xFFFFFFFF),
                            uintptr((dsid >> 32) & 0xFFFFFFFF),
                            0,
                            uintptr(syncType),
                            uintptr(unsafe.Pointer(&buf)),
                            uintptr(unsafe.Pointer(&size)),
                        )
            status = getStatus(st)

        case 8:
            st, _, _ := itunes.KbsyncGetData.Call(
                            self.session,
                            uintptr(dsid),
                            0,
                            uintptr(syncType),
                            uintptr(unsafe.Pointer(&buf)),
                            uintptr(unsafe.Pointer(&size)),
                        )
            status = getStatus(st)
    }

    if status != 0 {
        Raise(newiTunesHelperErrorf("Keybag.GetData failed: %X", uint(status)))
    }

    defer FreeSessionData(buf)

    return toBytes(buf, size)
}

func (self *KeybagSession) Import(keybag []byte) {
    status := self.importKeybag(keybag)

    switch status {
        case -42001, -42003:
            status = self.validate()
            if status == -42153 {
                self.Close()
                self.initialize(self.uniqueDeviceID)
            }

            status = self.importKeybag(keybag)
    }

    if status != 0 {
        Raise(newiTunesHelperErrorf("Keybag.Import failed: %d", status))
    }
}

func (self *KeybagSession) SaveDsid(dsid int64) {
    var status int

    switch unsafe.Sizeof(uintptr(0)) {
        case 4:
            st, _, _ := itunes.KbsyncSaveDsid.Call(
                            self.session,
                            uintptr(dsid & 0xFFFFFFFF),
                            uintptr((dsid >> 32) & 0xFFFFFFFF),
                        )
            status = getStatus(st)

        case 8:
            st, _, _ := itunes.KbsyncSaveDsid.Call(
                            self.session,
                            uintptr(dsid),
                        )
            status = getStatus(st)
    }

    if status != 0 {
        Raise(newiTunesHelperErrorf("Keybag.SaveDsid failed: %d", status))
    }
}

func (self *KeybagSession) importKeybag(keybag []byte) int {
    st, _, _ := itunes.KbsyncImport.Call(self.session, uintptr(unsafe.Pointer(&keybag[0])), uintptr(len(keybag)))
    return getStatus(st)
}
