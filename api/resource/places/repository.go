package places

import (
	"context"
	txUtil "cuide/util/db-tx"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(page uint8) (Places, error) {
	places := make([]*Place, 0)

	rows, err := r.db.Query(`SELECT * FROM get_servicos() LIMIT 20 OFFSET $1;`, (page-1)*20)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			place                                       Place
			serviceTypeJson, segmentJson, regionalsJson string
		)

		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Address,
			&place.PhoneNumber,
			&place.Website,
			&place.Observations,
			&place.GoogleMapsLink,
			&place.GoogleMapsEmbedLink,
			&place.AdmissionCriteria,
			&place.AttendanceType,
			&place.ReferenceWay,
			&serviceTypeJson,
			&segmentJson,
			&regionalsJson,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(serviceTypeJson), &place.ServiceType)
		json.Unmarshal([]byte(segmentJson), &place.Segment)
		json.Unmarshal([]byte(regionalsJson), &place.Regionals)

		places = append(places, &place)
	}

	return places, nil
}

func (r *Repository) Create(ctx context.Context, place *Place) (*Place, error) {
	var placeID uint8

	err := r.db.QueryRow(
		`INSERT INTO
		public.servico (
			tipo_servico_id,
			nome,
			endereco,
			contato,
			site,
			observacoes,
			eixo_id,
			maps_link,
			google_maps_embed_link,
			criterios_admissao,
			tipo_atendimento,
			forma_encaminhamento
		)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;`,
		place.ServiceType.ID,
		place.Name,
		place.Address,
		place.PhoneNumber,
		place.Website,
		place.Observations,
		place.Segment.ID,
		place.GoogleMapsLink,
		place.GoogleMapsEmbedLink,
		place.AdmissionCriteria,
		place.AttendanceType,
		place.ReferenceWay,
	).Scan(&placeID)
	if err != nil {
		return nil, err
	}

	stmt, err := r.db.Prepare(
		`INSERT INTO public.regionais_servico (servico_id, regional_id) VALUES ($1, $2)`,
	)
	if err != nil {
		return nil, err
	}
	for _, rs := range place.Regionals {
		_, err := stmt.Exec(placeID, rs.ID)
		if err != nil {
			return nil, err
		}
	}

	return place, nil
}

func (r *Repository) Read(id uint8) (*Place, error) {
	var (
		place                                       Place
		serviceTypeJson, segmentJson, regionalsJson string
	)

	err := r.db.QueryRow(`
	select
			s.id as servico_id,
			s.nome::text as servico_nome,
			s.endereco::text as servico_endereco,
			s.contato::text as servico_contato,
			s.site::text as servico_site,
			s.observacoes::text as servico_observacoes,
			s.maps_link::text as servico_maps_link,
			s.google_maps_embed_link::text as servico_maps_embed_link,
			s.criterios_admissao::text as servico_criterios_admissao,
			s.tipo_atendimento::text as servico_tipo_atendimento,
			s.forma_encaminhamento::text as servico_forma_encaminhamento,
			jsonb_build_object('id',
			ts.id,
			'name',
			ts.nome) as tipo_servico,
			jsonb_build_object('id',
			e.id,
			'name',
			e.nome) as eixo,
			jsonb_agg(distinct jsonb_build_object('id',
			r.id,
			'name',
			r.nome)) as regionais
		from
			public.servico s
		left join public.tipo_servico ts
						on
			s.tipo_servico_id = ts.id
		left join public.eixo e
						on
			s.eixo_id = e.id
		left join public.regionais_servico rs
						on
			s.id = rs.servico_id
		left join public.regionais r
						on
			rs.regional_id = r.id
	WHERE 
			s.id = $1
	GROUP BY
			s.id, ts.id, ts.nome, e.id, e.nome
	ORDER BY
			s.id`, id).
		Scan(
			&place.ID,
			&place.Name,
			&place.Address,
			&place.PhoneNumber,
			&place.Website,
			&place.Observations,
			&place.GoogleMapsLink,
			&place.GoogleMapsEmbedLink,
			&place.AdmissionCriteria,
			&place.AttendanceType,
			&place.ReferenceWay,
			&serviceTypeJson,
			&segmentJson,
			&regionalsJson,
		)
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(serviceTypeJson), &place.ServiceType)
	json.Unmarshal([]byte(segmentJson), &place.Segment)
	json.Unmarshal([]byte(regionalsJson), &place.Regionals)

	return &place, nil
}

func (r *Repository) Delete(id uint8) (int64, error) {
	var rows int64

	err := txUtil.CallTx(context.Background(), r.db, func(tx *sql.Tx) error {
		result, err := r.db.Exec("DELETE FROM public.regionais_servico WHERE servico_id = $1;", id)
		if err != nil {
			return err
		}

		rws, err := result.RowsAffected()
		if err != nil {
			return err
		}
		rows = rws

		result, err = r.db.Exec("DELETE FROM public.servico WHERE id = $1;", id)
		if err != nil {
			return err
		}

		rws, err = result.RowsAffected()
		if err != nil {
			return err
		}
		rows = rws

		return nil
	})

	return rows, err
}

func (r *Repository) Filter(filters Filters, page uint8) (Places, string, error) {
	conditionals := ``

	st := make([]string, len(filters.ServiceTypes))
	for i, service_type := range filters.ServiceTypes {
		st[i] = fmt.Sprintf("s.tipo_servico_id = %d", service_type)
	}
	if len(st) > 0 {
		conditionals += " AND "
		conditionals += fmt.Sprintf("(%s)", strings.Join(st, " OR "))
	}

	sg := make([]string, len(filters.Segments))
	for i, segment := range filters.Segments {
		sg[i] = fmt.Sprintf("s.eixo_id = %d", segment)
	}
	if len(sg) > 0 {
		conditionals += " AND "
		conditionals += fmt.Sprintf("(%s)", strings.Join(sg, " OR "))
	}

	rg := make([]string, len(filters.Regionals))
	for i, regional := range filters.Regionals {
		rg[i] = fmt.Sprintf("r.id = %d", regional)
	}
	if len(rg) > 0 {
		conditionals += " AND "
		conditionals += fmt.Sprintf("(%s)", strings.Join(rg, " OR "))
	}

	if filters.Name != "" {
		conditionals += " AND "
		conditionals += fmt.Sprintf(
			" (lower(s.nome) LIKE lower('%%%s%%') OR lower(s.tipo_atendimento) LIKE lower('%%%s%%'))",
			filters.Name,
			filters.Name,
		)
	}

	query := fmt.Sprintf(`
	select
		gs.*
	from
		get_servicos() gs,
		(
			select
				s.id as servico_id
			from
				public.servico s
				left join public.tipo_servico ts on s.tipo_servico_id = ts.id
				left join public.eixo e on s.eixo_id = e.id
				left join public.regionais_servico rs on s.id = rs.servico_id
				left join public.regionais r on rs.regional_id = r.id
			where
				1 = 1 %s
			group by
				s.id,
				ts.id,
				ts.nome,
				e.id,
				e.nome
			order by
				s.id
		) ss
	where
		ss.servico_id = gs.servico_id
	limit
		20 offset $1;`, conditionals)

	fmt.Println(query, filters)

	rows, err := r.db.Query(query, (page-1)*10)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	places := make([]*Place, 0)
	for rows.Next() {
		var (
			place                                       Place
			serviceTypeJson, segmentJson, regionalsJson string
		)

		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Address,
			&place.PhoneNumber,
			&place.Website,
			&place.Observations,
			&place.GoogleMapsLink,
			&place.GoogleMapsEmbedLink,
			&place.AdmissionCriteria,
			&place.AttendanceType,
			&place.ReferenceWay,
			&serviceTypeJson,
			&segmentJson,
			&regionalsJson,
		)
		if err != nil {
			return nil, "", err
		}

		json.Unmarshal([]byte(serviceTypeJson), &place.ServiceType)
		json.Unmarshal([]byte(segmentJson), &place.Segment)
		json.Unmarshal([]byte(regionalsJson), &place.Regionals)

		places = append(places, &place)
	}

	return places, conditionals, nil
}

func (r *Repository) PaginationMetadata() (pm PaginationMetadata, err error) {
	err = r.db.QueryRow(`
	SELECT 
			COUNT(s.id) AS "total", 
			CEIL(COUNT(s.id)::FLOAT / 20) AS "pages" 
	FROM 
			public.servico s;`,
	).Scan(&pm.Metadata.TotalPlaces, &pm.Metadata.Pages)

	return
}

func (r *Repository) FilterPaginationMetadata(filters string) (pm PaginationMetadata, err error) {
	query := fmt.Sprintf(`
	select
		COUNT(f.servico_id) as "total", 
		ceil(COUNT(f.servico_id)::FLOAT / 20) as "pages"
	from
		(
		select
			gs.*
		from
			get_servicos() gs,
			(
			select
				s.id as servico_id
			from
				public.servico s
			left join public.tipo_servico ts on
				s.tipo_servico_id = ts.id
			left join public.eixo e on
				s.eixo_id = e.id
			left join public.regionais_servico rs on
				s.id = rs.servico_id
			left join public.regionais r on
				rs.regional_id = r.id
			where
				1 = 1 %s
			group by
				s.id,
				ts.id,
				ts.nome,
				e.id,
				e.nome
			order by
				s.id
									) ss
		where
			ss.servico_id = gs.servico_id) f`, filters,
	)

	err = r.db.QueryRow(query).Scan(&pm.Metadata.TotalPlaces, &pm.Metadata.Pages)

	return
}

func (r *Repository) Update(ctx context.Context, place *Place) (int64, error) {
	var rowsAffected int64

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(
		ctx,
		`UPDATE public.servico
		SET
			tipo_servico_id = $1,
			nome = $2,
			endereco = $3,
			contato = $4,
			site = $5,
			observacoes = $6,
			eixo_id = $7,
			maps_link = $8,
			google_maps_embed_link = $9,
			criterios_admissao = $10,
			tipo_atendimento = $11,
			forma_encaminhamento = $12
		WHERE id = $13`,
		place.ServiceType.ID,
		place.Name,
		place.Address,
		place.PhoneNumber,
		place.Website,
		place.Observations,
		place.Segment.ID,
		place.GoogleMapsLink,
		place.GoogleMapsEmbedLink,
		place.AdmissionCriteria,
		place.AttendanceType,
		place.ReferenceWay,
		place.ID,
	)
	if err != nil {
		return 0, err
	}

	rw, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	rowsAffected += rw

	result, err = tx.ExecContext(
		ctx,
		`DELETE FROM public.regionais_servico WHERE servico_id = $1`,
		place.ID,
	)
	if err != nil {
		return 0, err
	}

	rw, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}
	rowsAffected += rw

	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO public.regionais_servico (servico_id, regional_id) VALUES ($1, $2)`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, rs := range place.Regionals {
		result, err = stmt.Exec(place.ID, rs.ID)
		if err != nil {
			return 0, err
		}

		rw, err = result.RowsAffected()
		if err != nil {
			return 0, err
		}
		rowsAffected += rw
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
