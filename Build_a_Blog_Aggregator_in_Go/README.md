# Gator CLI

Gator 是一個用 Go 語言編寫的命令列工具，用來幫助你管理並抓取 RSS 訂閱源。

## 先決條件 (Prerequisites)

在執行此程式之前，請確保你的系統中已安裝以下軟體：
* [Go](https://go.dev/)
* [PostgreSQL](https://www.postgresql.org/)

## 安裝方式 (Installation)

你可以使用 `go install` 將 Gator 安裝到系統環境中。這樣你就可以在任何地方直接使用 `gator` 指令，而不需要每次都執行 `go run .`：

```bash
go install github.com/Alan011337/Build_a_Blog_Aggregator_in_Go@latest
```

## 設定檔 (Configuration)

請在你的家目錄（Home directory）建立一個名為 `.gatorconfig.json` 的檔案，並放入以下內容來設定資料庫與當前使用者：

```json
{
  "db_url": "postgres://your_username:your_password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "kahya"
}
```

## 使用方式 (Usage / Commands)

安裝完成並設定好 `.gatorconfig.json` 後，你可以使用以下指令：

| 指令 (Command) | 參數 (Arguments) | 說明 (Description) |
| :--- | :--- | :--- |
| `register` | `<name>` | 註冊一個新使用者，並自動將其設為當前使用者。 |
| `login` | `<name>` | 切換當前登入的使用者。 |
| `addfeed` | `<name>`, `<url>` | 新增一個 feed。 |
| `agg` | `<duration>` | 每隔指定的持續時間（例如 1m, 1h）自動抓取所有訂閱源的最新文章。 |
|`browse`| `[limit]` | 瀏覽最新的幾篇文章。可選填限制數量（預設顯示 2 篇）。 |