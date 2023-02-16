package store_test

import (
	"pois/models"
	"pois/store"
	"testing"

	"git.eng.vecima.com/cloud/golib/v4/zaplogger"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type AliasSuite struct {
	suite.Suite
	aliasCount uint
}

var AliasStore *store.SQLAlias

func TestAliasStore(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(AliasSuite))
}

func (suite *AliasSuite) SetupSuite() {
	logger := zaplogger.MainLog
	AliasStore = &store.SQLAlias{DB: DB, Log: logger}
}

func (suite *AliasSuite) SetupTest() {
	AliasStore.DB.Exec("DELETE FROM channel_alias")
	suite.aliasCount = 0
}

func (suite *AliasSuite) insertAlias(channelname string, aliasname string) {
	AliasStore.DB.MustExec(
		`INSERT INTO channel_alias(channelname, aliasname) VALUES ($1, $2)`,
		channelname, aliasname,
	)
}

func (suite *AliasSuite) TestDeleteAlias() {
	suite.insertAlias("testChannel", "testAlias")

	err := AliasStore.DeleteAlias("testChannel", "testAlias")
	suite.NoError(err)

	_, err = AliasStore.FindAlias("testChannel")
	suite.Nil(err)

	err = AliasStore.DeleteAlias("testChannel", "testAlias")
	suite.Error(err, "shouldn't be able to delete a deleted user")
}

func (suite *AliasSuite) TestFindAlias() {
	suite.insertAlias("testChannelx", "testAliasx")

	alias, err := AliasStore.FindAlias("testChannelx")

	if suite.NoError(err) {
		for _, element := range alias {
			suite.Equal("testChannelx", element.Channelname)
			suite.Equal("testAliasx", element.AliasName)

			suite.NotEqual("testChannely", element.Channelname)
			suite.NotEqual("testAliasy", element.AliasName)
		}
	}

	alias, err = AliasStore.FindAllAlias()

	if suite.NoError(err) {
		for _, element := range alias {
			suite.Equal("testChannelx", element.Channelname)
			suite.Equal("testAliasx", element.AliasName)

			suite.NotEqual("testChannely", element.Channelname)
			suite.NotEqual("testAliasy", element.AliasName)
		}
	}

	_, err = AliasStore.FindAlias("testChannely")
	suite.Nil(err)

	err = AliasStore.DeleteAlias("testChannelx", "testAliasx")
	suite.NoError(err)
}

func (suite *AliasSuite) TestCreateAlias() {
	request := models.Alias{Channelname: "foo", AliasName: "zzz"}
	alias, err := AliasStore.CreateAlias(request)

	if suite.NoError(err) {
		for _, element := range alias {
			suite.Equal("foo", element.Channelname)
			suite.Equal("zzz", element.AliasName)

			suite.NotEqual("testChannely", element.Channelname)
			suite.NotEqual("testAliasy", element.AliasName)
		}
	}

	err = AliasStore.DeleteAlias("foo", "zzz")
	suite.NoError(err)
}
