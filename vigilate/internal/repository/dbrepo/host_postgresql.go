package dbrepo

import (
	"context"
	"log"
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

	if err != nil {
		return host, err
	}

	hostServicesQuery := `SELECT
		hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit,
		hs.last_check, hs.status, hs.created_at, hs.updated_at,
		s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at
		FROM host_services hs left join services s on (s.id = hs.service_id)
		WHERE host_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, hostServicesQuery, id)
	if err != nil {
		log.Println(err)
		return host, err
	}
	defer rows.Close()

	var hostServices []models.HostService
	for rows.Next() {
		var hs models.HostService
		err := rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.LastCheck,
			&hs.Status,
			&hs.CreatedAt,
			&hs.UpdatedAt,
			&hs.Service.ID,
			&hs.Service.ServiceName,
			&hs.Service.Active,
			&hs.Service.Icon,
			&hs.Service.CreatedAt,
			&hs.Service.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return host, err
		}
		hostServices = append(hostServices, hs)
	}

	host.HostServices = hostServices
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

		hst, _ := m.GetHostById(pctx, host.ID)
		hosts = append(hosts, hst)
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

// GetAllActiveServicesCountByStatus Returns all services counts w.r.t status
func (repo *postgresDBRepo) GetAllActiveServicesCountByStatus(pctx context.Context) (models.StatusCount, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	query := ` select
		(select count(id) from host_services where active = 1 and status = 'pending') as pending,
		(select count(id) from host_services where active = 1 and status = 'healthy') as healthy,
		(select count(id) from host_services where active = 1 and status = 'warning') as warning,
		(select count(id) from host_services where active = 1 and status = 'problem') as problem
	`
	var counts models.StatusCount

	row := repo.DB.QueryRowContext(ctx, query)

	err := row.Scan(
		&counts.Pending,
		&counts.Healthy,
		&counts.Warning,
		&counts.Problem,
	)

	return counts, err
}

// GetServicesByStatus returns all active services with a given status
func (m *postgresDBRepo) GetServicesByStatus(pctx context.Context, status string) ([]models.HostService, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	query := `
		select
			hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit,
			hs.last_check, hs.status, hs.created_at, hs.updated_at,
			h.host_name, s.service_name
		from
			host_services hs
			left join hosts h on (hs.host_id = h.id)
			left join services s on (hs.service_id = s.id)
		where
			status = $1
			and hs.active = 1
		order by host_name, service_name`

	var services []models.HostService

	rows, err := m.DB.QueryContext(ctx, query, status)
	if err != nil {
		return services, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.HostService

		err := rows.Scan(
			&h.ID,
			&h.HostID,
			&h.ServiceID,
			&h.Active,
			&h.ScheduleNumber,
			&h.ScheduleUnit,
			&h.LastCheck,
			&h.Status,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.HostName,
			&h.Service.ServiceName,
		)
		if err != nil {
			return nil, err
		}

		services = append(services, h)
	}

	return services, nil
}
