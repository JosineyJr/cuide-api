package service_types

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (ServiceTypes, error) {
	serviceTypes := make([]*ServiceType, 0)

	rows, err := r.db.Query("SELECT id, nome FROM public.tipo_servico;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var serviceType ServiceType
		rows.Scan(&serviceType.ID, &serviceType.Name)

		serviceTypes = append(serviceTypes, &serviceType)
	}

	return serviceTypes, nil
}

func (r *Repository) Create(serviceType *ServiceType) (*ServiceType, error) {
	err := r.db.QueryRow("INSERT INTO public.tipo_servico (nome) VALUES ($1) RETURNING id;", serviceType.Name).
		Scan(&serviceType.ID)
	if err != nil {
		return nil, err
	}

	return serviceType, nil
}

func (r *Repository) Read(id uint8) (*ServiceType, error) {
	var serviceType ServiceType
	err := r.db.QueryRow("SELECT id, nome FROM public.tipo_servico r WHERE r.id = $1;", id).
		Scan(&serviceType.ID, &serviceType.Name)
	if err != nil {
		return nil, err
	}

	return &serviceType, nil
}

func (r *Repository) Update(serviceType *ServiceType) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.tipo_servico SET nome = $1 WHERE id= $2;",
		serviceType.Name,
		serviceType.ID,
	)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}

func (r *Repository) Delete(id uint8) (int64, error) {
	result, err := r.db.Exec("DELETE FROM public.tipo_servico WHERE id = $1;", id)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}
