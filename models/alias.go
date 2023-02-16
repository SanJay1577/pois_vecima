package models

// DB model for channel_alias
type Alias struct {
	Channelname string `db:"channelname"`
	AliasName   string `db:"aliasname"`
}
