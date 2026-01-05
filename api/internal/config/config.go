// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

type Config struct {
	Name   string
	Host   string
	Port   int
	DB     DBConfig
}

type DBConfig struct {
	Default DefaultDB
}

type DefaultDB struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
	LogLevel        string
	SlowThreshold   int
	SkipDefaultTxn  bool
	PrepareStmt     bool
	SingularTable   bool
	DisableForeignKey bool
}
