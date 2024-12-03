package router

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	admission_criteria "cuide/api/resource/admission-criteria"
	attendance_types "cuide/api/resource/attendance-types"
	"cuide/api/resource/health"
	"cuide/api/resource/places"
	reference_ways "cuide/api/resource/reference-ways"
	"cuide/api/resource/regionals"
	"cuide/api/resource/segments"
	service_types "cuide/api/resource/service-types"
	"cuide/api/router/middleware"
	"cuide/api/router/middleware/requestlog"
)

func New(l *zerolog.Logger, v *validator.Validate, db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/health", health.Read)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)

		regionalAPI := regionals.New(l, v, db)
		r.Method(http.MethodGet, "/regionals", requestlog.NewHandler(regionalAPI.List, l))
		r.Method(http.MethodPost, "/regionals", requestlog.NewHandler(regionalAPI.Create, l))
		r.Method(http.MethodGet, "/regionals/{id}", requestlog.NewHandler(regionalAPI.Read, l))
		r.Method(http.MethodPut, "/regionals/{id}", requestlog.NewHandler(regionalAPI.Update, l))
		r.Method(http.MethodDelete, "/regionals/{id}", requestlog.NewHandler(regionalAPI.Delete, l))

		attendanceTypeAPI := attendance_types.New(l, v, db)
		r.Method(
			http.MethodGet,
			"/attendance-types",
			requestlog.NewHandler(attendanceTypeAPI.List, l),
		)
		r.Method(
			http.MethodPost,
			"/attendance-types",
			requestlog.NewHandler(attendanceTypeAPI.Create, l),
		)
		r.Method(
			http.MethodGet,
			"/attendance-types/{id}",
			requestlog.NewHandler(attendanceTypeAPI.Read, l),
		)
		r.Method(
			http.MethodPut,
			"/attendance-types/{id}",
			requestlog.NewHandler(attendanceTypeAPI.Update, l),
		)
		r.Method(
			http.MethodDelete,
			"/attendance-types/{id}",
			requestlog.NewHandler(attendanceTypeAPI.Delete, l),
		)

		segmentAPI := segments.New(l, v, db)
		r.Method(http.MethodGet, "/segments", requestlog.NewHandler(segmentAPI.List, l))
		r.Method(http.MethodPost, "/segments", requestlog.NewHandler(segmentAPI.Create, l))
		r.Method(http.MethodGet, "/segments/{id}", requestlog.NewHandler(segmentAPI.Read, l))
		r.Method(http.MethodPut, "/segments/{id}", requestlog.NewHandler(segmentAPI.Update, l))
		r.Method(http.MethodDelete, "/segments/{id}", requestlog.NewHandler(segmentAPI.Delete, l))

		serviceTypeAPI := service_types.New(l, v, db)
		r.Method(http.MethodGet, "/service-types", requestlog.NewHandler(serviceTypeAPI.List, l))
		r.Method(http.MethodPost, "/service-types", requestlog.NewHandler(serviceTypeAPI.Create, l))
		r.Method(
			http.MethodGet,
			"/service-types/{id}",
			requestlog.NewHandler(serviceTypeAPI.Read, l),
		)
		r.Method(
			http.MethodPut,
			"/service-types/{id}",
			requestlog.NewHandler(serviceTypeAPI.Update, l),
		)
		r.Method(
			http.MethodDelete,
			"/service-types/{id}",
			requestlog.NewHandler(serviceTypeAPI.Delete, l),
		)

		referenceWaysAPI := reference_ways.New(l, v, db)
		r.Method(http.MethodGet, "/reference-ways", requestlog.NewHandler(referenceWaysAPI.List, l))
		r.Method(
			http.MethodPost,
			"/reference-ways",
			requestlog.NewHandler(referenceWaysAPI.Create, l),
		)
		r.Method(
			http.MethodGet,
			"/reference-ways/{id}",
			requestlog.NewHandler(referenceWaysAPI.Read, l),
		)
		r.Method(
			http.MethodPut,
			"/reference-ways/{id}",
			requestlog.NewHandler(referenceWaysAPI.Update, l),
		)
		r.Method(
			http.MethodDelete,
			"/reference-ways/{id}",
			requestlog.NewHandler(referenceWaysAPI.Delete, l),
		)

		admissionCriteriaAPI := admission_criteria.New(l, v, db)
		r.Method(
			http.MethodGet,
			"/admission-criteria",
			requestlog.NewHandler(admissionCriteriaAPI.List, l),
		)
		r.Method(
			http.MethodPost,
			"/admission-criteria",
			requestlog.NewHandler(admissionCriteriaAPI.Create, l),
		)
		r.Method(
			http.MethodGet,
			"/admission-criteria/{id}",
			requestlog.NewHandler(admissionCriteriaAPI.Read, l),
		)
		r.Method(
			http.MethodPut,
			"/admission-criteria/{id}",
			requestlog.NewHandler(admissionCriteriaAPI.Update, l),
		)
		r.Method(
			http.MethodDelete,
			"/admission-criteria/{id}",
			requestlog.NewHandler(admissionCriteriaAPI.Delete, l),
		)

		placeAPI := places.New(l, v, db)
		r.Method(http.MethodGet, "/places", requestlog.NewHandler(placeAPI.List, l))
		r.Method(http.MethodPost, "/places", requestlog.NewHandler(placeAPI.Create, l))
	})

	return r
}
