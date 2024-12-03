package attendance_types

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

func (r *Repository) List() (AttendanceTypes, error) {
	attendanceTypes := make([]*AttendanceType, 0)

	rows, err := r.db.Query("SELECT id, nome FROM public.tipo_atendimento;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var regional AttendanceType
		rows.Scan(&regional.ID, &regional.Name)

		attendanceTypes = append(attendanceTypes, &regional)
	}

	return attendanceTypes, nil
}

func (r *Repository) Create(attendanceType *AttendanceType) (*AttendanceType, error) {
	err := r.db.QueryRow("INSERT INTO public.tipo_atendimento (nome) VALUES ($1) RETURNING id;", attendanceType.Name).
		Scan(&attendanceType.ID)
	if err != nil {
		return nil, err
	}

	return attendanceType, nil
}

func (r *Repository) Read(id uint8) (*AttendanceType, error) {
	var attendanceType AttendanceType
	err := r.db.QueryRow("SELECT id, nome FROM public.tipo_atendimento r WHERE r.id = $1;", id).
		Scan(&attendanceType.ID, &attendanceType.Name)
	if err != nil {
		return nil, err
	}

	return &attendanceType, nil
}

func (r *Repository) Update(attendanceType *AttendanceType) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.tipo_atendimento SET nome = $1 WHERE id= $2;",
		attendanceType.Name,
		attendanceType.ID,
	)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}

func (r *Repository) Delete(id uint8) (int64, error) {
	result, err := r.db.Exec("DELETE FROM public.tipo_atendimento WHERE id = $1;", id)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}
