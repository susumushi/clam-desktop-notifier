# clam-desktop-notifier
clamavのVirusEventに渡すといい感じにデスクトップ通知してくれるヤツ

## つかいかた

- バイナリの配置
```
curl `curl  https://api.github.com/repos/susumushi/clam-desktop-notifier/releases/latest | jq -r .assets[].browser_download_url` -L -o /etc/clamav/virusevent.d/clamnotify
```

- (必要に応じて)権限の変更

```
chown root:clamav /etc/clamav/virusevent.d/clamnotify 
chmod 4755 /etc/clamav/virusevent.d/clamnotify 
```

- Clamavの設定
```
VirusEvent /etc/clamav/virusevent.d/clamnotify
```