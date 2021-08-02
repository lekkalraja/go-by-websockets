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

	if err != nil {
		return 0, err
	}

	// add host services and set to inactive
	stmt := `
		insert into host_services (host_id, service_id, active, schedule_number, schedule_unit,
		status, created_at, updated_at) values ($1, 1, 0, 3, 'm', 'pending', $2, $3)
`

	_, err = repo.DB.ExecContext(ctx, stmt, id, time.Now(), time.Now())
	if err != nil {
		return id, err
	}

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

func (repo *postgresDBRepo) UpdateHostById(pctx context.Context, host models.Host) error {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	stmt := ` update hosts
			  set host_name = $1, canonical_name = $2,
			  url = $3, ip = $4, ipv6 = $5, location = $6, active = $7, os = $8,
			  updated_at = $9 where id = $10
	`

	_, err := repo.DB.ExecContext(ctx, stmt,
		host.HostName,
		host.CanonicalName,
		host.URL,
		host.IP,
		host.IPV6,
		host.Location,
		host.Active,
		host.OS,
		time.Now(),
		host.ID,
	)
	return err
}

// GetHostById returns a Host by id
func (m *postgresDBRepo) GetAllHosts(pctx context.Context) ([]models.Host, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, host_name, canonical_name,  url, ip, ipv6, location, active, os,
			 created_at, updated_at
			 FROM hosts`
	rows, err := m.DB.QueryContext(ctx, stmt)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []models.Host

	for rows.Next() {
		var host models.Host
		err = rows.Scan(
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
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return hosts, nil
}

func (repo *postgresDBRepo) UpdateHostServiceStatus(pctx context.Context, hostId, serviceId, active int) error {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	stmt := `UPDATE host_services
			 SET  active = $1
			 WHERE host_id = $2 and service_id = $3
	`
	_, err := repo.DB.ExecContext(ctx, stmt, active, hostId, serviceId)
	return err
}
