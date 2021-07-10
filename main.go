package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/coreos/go-systemd/journal"
	"github.com/godbus/dbus/v5"
)

func main() {
	// 現在ログインしているユーザー一覧を取得する
	// 端末がスクリーンになっているユーザーのみ抽出する
	output, err := exec.Command("/bin/bash", "-c", "who | grep -e ' :[0-9]' | awk '{print $1}' | sort | uniq").Output()
	if err != nil {
		panic(fmt.Sprintf("get logined user list error: %s", err))
	}
	userNames := strings.Split(string(output), "\n")
	userNames = userNames[:len(userNames)-1]

	// CLAMAVは環境変数に検知したウイルス名を入れてくれるので取得する
	virusName := os.Getenv("CLAM_VIRUSEVENT_VIRUSNAME")
	virusFile := os.Getenv("CLAM_VIRUSEVENT_FILENAME")
	alertMsg := fmt.Sprintf("Virus Detected: VirusName: %s, FileName: %s", virusName, virusFile)

	if err := journal.Print(journal.PriCrit, alertMsg); err != nil {
		panic(fmt.Sprintf("send logging faluer: %s", err))
	}
	for _, userName := range userNames {
		// ログインしているユーザーのUIDを取得し、スイッチする。
		u, _ := user.Lookup(userName)
		intUID, _ := strconv.Atoi(u.Uid)
		syscall.Setuid(intUID)

		// 各ユーザーのDBUSにアクセスして通知を要求する
		conn, err := dbus.Connect(fmt.Sprintf("unix:path=/run/user/%d/bus", intUID))
		if err != nil {
			panic(fmt.Sprintf("dbus connect error: %s", err))
		}
		defer conn.Close()
		obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
		notifyAppName := "clamav"
		notifyIcon := ""
		notifyReplaceID := uint32(0)
		notifySummary := "Clamav Virus detected."
		notifyDesc := alertMsg
		notifyActions := []string{}
		dbusHint := map[string]dbus.Variant{} //多分通知では利用しなんじゃないかな？
		notifyExpire := int32(5000)

		call := obj.Call("org.freedesktop.Notifications.Notify", 0,
			notifyAppName, notifyReplaceID,
			notifyIcon, notifySummary, notifyDesc, notifyActions,
			dbusHint, notifyExpire)
		if call.Err != nil {
			panic(fmt.Sprintf("send notify error: %s", call.Err))
		}
	}
}
