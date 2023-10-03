package repository

import (
	_interface "github.com/xvbnm48/go-api-catatan/repository/interface"
	"github.com/xvbnm48/go-api-catatan/repository/mysql"
	"io"
)

type DbConf struct {
	URL, Port, Schema, User, Password string
}

type RepoConfigs struct {
	DbConf DbConf
}

type Repo struct {
	DbReadWriter _interface.ReadWriter
	io.Closer
}

func NewNoteServiceRepository(rc RepoConfigs) (*Repo, error) {
	readWriter, err := mysql.NewDbReadWriter(rc.DbConf.URL, rc.DbConf.Port, rc.DbConf.Schema, rc.DbConf.User, rc.DbConf.Password)
	if err != nil {
		return nil, err
	}

	return &Repo{
		DbReadWriter: readWriter,
	}, nil
}
