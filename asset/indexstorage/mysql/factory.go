package mysql

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/janoszen/openshiftci_inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci_inspector/common/mysql"
)

const (
	// language=MySQL
	createTableSQL = `
CREATE TABLE IF NOT EXISTS job_assets (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	job_id VARCHAR(255) NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
	UNIQUE u_asset(job_id, asset_name),
	INDEX i_job_id(job_id)
)
`
)

// NewMySQLAssetIndex creates a MySQL storage for asset indexes.
func NewMySQLAssetIndex(config mysql.Config, logger *log.Logger) (indexstorage.AssetIndex, error) {
	db, err := config.CreateMySQLDB()
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create assets table (%w)", err)
	}

	return &mysqlAssetIndex{
		db:     db,
		logger: logger,
	}, nil
}
