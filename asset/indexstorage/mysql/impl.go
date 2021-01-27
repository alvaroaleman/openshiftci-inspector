package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlAssetIndex struct {
	db     *sql.DB
	logger *log.Logger
}

func (m *mysqlAssetIndex) AddAsset(jobID string, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := m.db.QueryContext(
		ctx,
		`INSERT INTO job_assets (job_id, asset_name) VALUES (?, ?)`,
		jobID,
		name,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Close()
	}()
	return nil
}

func (m *mysqlAssetIndex) HasAsset(jobID string, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := m.db.QueryContext(
		ctx,
		`SELECT COUNT(*) cnt FROM job_assets WHERE job_id=? AND asset_name=?`,
		jobID,
		name,
	)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = res.Close()
	}()
	if !res.Next() {
		return false, errors.New("no rows returned from HasAsset query")
	}
	var count int
	if err := res.Scan(&count); err != nil {
		return false, err
	}
	return count == 1, nil
}

func (m *mysqlAssetIndex) ListAssets(jobID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := m.db.QueryContext(
		ctx,
		`SELECT asset_name FROM job_assets WHERE job_id=?`,
		jobID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()
	assets := []string{}
	for {
		if !res.Next() {
			break
		}
		asset := ""
		if err := res.Scan(&asset); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func (m *mysqlAssetIndex) Shutdown(_ context.Context) {
	_ = m.db.Close()
}
