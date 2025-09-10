package routes

import (
	app "Salehaskarzadeh/internal/appl"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(a *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(a.Middleware.Authenticate)

		r.Get("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", a.Middleware.RequireUser(a.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", a.Middleware.RequireUser(a.WorkoutHandler.HandleDeleteWorkoutByID))

	})

	r.Get("/health", a.HealthCheck)

	r.Post("/users", a.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", a.TokenHandler.HandleCreateToken)

	return r
}
