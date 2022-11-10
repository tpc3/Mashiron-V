package config

type Strings struct {
	Lang    string
	Help    string
	Done    string
	Warning warningstr
	Error   errorstr
}

type warningstr struct {
	Title     string
	Overwrite string
}

type errorstr struct {
	Title      string
	Unknown    string
	SubCmd     string
	Permission string
	ZeroEntry  string
	NoEntry    string
	Invalid    string
	NoEmoji    string
}

var (
	Lang map[string]Strings
)

func loadLang() {
	Lang = map[string]Strings{}
	Lang["japanese"] = Strings{
		Lang: "japanese",
		Help: "Botの使い方に関しては、Wikiをご覧ください。",
		Done: "操作は正常に完了しました。",
		Warning: warningstr{
			Title:     "注意",
			Overwrite: "この操作は既存の以下の項目を上書きします。",
		},
		Error: errorstr{
			Title:      "エラー",
			Unknown:    "不明なエラーが発生しました。\nこの問題は管理者に報告されます。",
			SubCmd:     "サブコマンドが不正です。",
			Permission: "あなたにはこのコマンドを実行する権限がありません。",
			ZeroEntry:  "登録されているコマンドが無いため、結果を表示できません。",
			NoEntry:    "そのようなコマンドはありません。",
			Invalid:    "コマンドの形式が不正です。",
			NoEmoji:    "絵文字が見つかりません。",
		},
	}
	Lang["english"] = Strings{
		Lang: "english",
		Help: "Usage is available on the Wiki.",
		Done: "Operation completed.",
		Warning: warningstr{
			Title:     "Warning",
			Overwrite: "This operation will overwrite this definition.",
		},
		Error: errorstr{
			Title:      "Error",
			Unknown:    "Unknown Error!\nThis error will be reported to the admin.",
			SubCmd:     "Invalid subcommand.",
			Permission: "You don't have permission to do that.",
			ZeroEntry:  "Cannot show results since there's no registered command.",
			NoEntry:    "There's no such command.",
			Invalid:    "Invalid command.",
			NoEmoji:    "Emoji not found.",
		},
	}
}
