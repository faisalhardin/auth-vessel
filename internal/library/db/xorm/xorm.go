package xormlib

import (
	"fmt"

	"github.com/faisalhardin/auth-vessel/internal/config"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"xorm.io/core"
)

type DBConnect struct {
	MasterDB *xorm.Engine
	SlaveDB  *xorm.Engine
}

func NewDBConnection(cfg *config.Config) (dbConnection *DBConnect, err error) {

	masterDB, err := generateXormEngineInstance(cfg.Vault.DBMaster)
	if err != nil {
		return nil, errors.New("failed to make connection to master db")
	}

	slaveDB, err := generateXormEngineInstance(cfg.Vault.DBSlave)
	if err != nil {
		return nil, errors.New("failed to make connection to slave db")
	}

	return &DBConnect{
		SlaveDB:  slaveDB,
		MasterDB: masterDB,
	}, nil
}

func (conn *DBConnect) CloseDBConnection() error {
	if conn.MasterDB != nil {
		err := conn.MasterDB.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close master db engine")
		}
	}

	if conn.SlaveDB != nil {
		err := conn.SlaveDB.Close()
		if err != nil {
			return errors.Wrap(err, "failed to close slave db engine")
		}
	}
	return nil
}

func generateXormEngineInstance(DBConfig map[string]string) (*xorm.Engine, error) {
	var dsn string
	for k, v := range DBConfig {
		if len(v) == 0 {
			continue
		}
		dsn = fmt.Sprintf("%s %s=%s", dsn, k, v)
	}

	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %v", err)
	}

	// Ping the database to verify the connection
	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	engine.SetTableMapper(core.GonicMapper{})
	engine.SetColumnMapper(core.GonicMapper{})

	return engine, nil

}
