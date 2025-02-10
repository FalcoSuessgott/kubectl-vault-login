package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FalcoSuessgott/kubectl-vault-login/pkg/utils"
)

// write each execcredential to a logfile (per default in /tmp, configurable via env var)
// current entry for
// previous once in previous
// add SA name

const (
	LogFileName       = ".kubectl-vault-login-log.json"
	LogFilePathEnvVar = "KUBECTL_VAULT_LOGIN_LOG_PATH"
)

type AuthLog struct {
	file string
	Log  *Log
}

type Log struct {
	Current  *Entry   `json:"current"`
	Previous []*Entry `json:"previous"`
}

type Entry struct {
	ServiceAccountName      string    `json:"service_account_name"`
	ServiceAccountNamespace string    `json:"service_account_namespace"`
	Uid                     string    `json:"uid"`
	ValidFrom               time.Time `json:"valid_from"`
	ValidUntil              time.Time `json:"valid_until"`
}

func New() (*AuthLog, error) {
	logFilePath := "/tmp/"

	if v, ok := os.LookupEnv(LogFilePathEnvVar); ok {
		logFilePath = v
	}

	// error if dir does not exists
	if err := utils.DirExists(logFilePath); err != nil {
		return nil, fmt.Errorf("log file directory does not exist %s: %w", logFilePath, err)
	}

	// error if dir is writable
	if err := utils.DirIsWritable(logFilePath); err != nil {
		return nil, fmt.Errorf("log file directory is not writable %s: %w", logFilePath, err)
	}

	logFile := path.Join(logFilePath, LogFileName)

	// check if file exists, otherwise create it
	if err := utils.CreateFileIfNotExists(logFile); err != nil {
		return nil, fmt.Errorf("error creating or appending to log file %s: %w", logFile, err)
	}

	return &AuthLog{
		file: logFile,
	}, nil
}

func (a *AuthLog) Read() ([]byte, error) {
	f, err := os.ReadFile(a.file)
	if err != nil {
		return nil, fmt.Errorf("error reading log file %s: %w", a.file, err)
	}

	// return if file empty
	if len(f) == 0 {
		return nil, nil
	}

	// otherwise try to parse it as json
	if err := json.Unmarshal(f, a.Log); err != nil {
		return nil, fmt.Errorf("invalid log file %s: %w", a.file, err)
	}

	return json.Marshal(a.Log)
}

// adds new entry as current, shifts previous current to previous
func (a *AuthLog) Write(data interface{}) error {
	f, err := os.ReadFile(a.file)
	if err != nil {
		return fmt.Errorf("error reading log file %s: %w", a.file, err)
	}

	if _, err := a.Read(); err != nil {
		return fmt.Errorf("error reading log file %s: %w", a.file, err)
	}

	if json.Unmarshal(f, a.Log) != nil {
		return fmt.Errorf("invalid log file %s: %w", a.file, err)
	}

	fmt.Println(a.Log)

	return nil
}
