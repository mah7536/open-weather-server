## 說明
此為抓取open weather 的server 
需要申請
1. telegram的bot token
2. 中央氣象局api用token
即可透過telegram機器人 直接抓取當地天氣
或 也可以透過telegram的send location功能 取得所在地的天氣資訊

## 使用方法
go run main.go -conf {{您的設定檔}}

# 設定檔案
使用setting.conf 中的格式 填入telegram token及 中央氣象局的token 即可
