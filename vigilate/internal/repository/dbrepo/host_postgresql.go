package dbrepo

import (
	"context"
	"time"

	"github.com/lekkalraja/go-by-websockets/vigilate/internal/models"
)

// InsertHost inserts Host record into the database
func (repo *postgresDBRepo) InsertHost(pctx context.Context, host models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	query := `insert into hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	var id int
	err := repo.DB.QueryRowContext(ctx, query,
		host.HostName,
		host.CanonicalName,
		host.URL,
		host.IP,
		host.IPV6,
		host.Location,
		host.OS,
		host.Active,
		time.Now(),
		time.Now(),
	).Scan(&id)
	return id, err
}

// GetHostById returns a Host by id
func (m *postgresDBRepo) GetHostById(pctx context.Context, id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, host_name, canonical_name,  url, ip, ipv6, location, active, os,
			created_at, updated_at
			FROM hosts where id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var host models.Host

	err := row.Scan(
		&host.ID,
		&host.HostName,
		&host.CanonicalName,
		&host.URL,
		&host.IP,
		&host.IPV6,
		&host.Location,
		&host.Active,
		&host.OS,
		&host.CreatedAt,
		&host.UpdatedAt,
	)
	return host, err
}
