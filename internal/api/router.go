// TODO: MAKE THIS FILE MORE READABLE AND ORGANIZED..:((((((
package api

import (
	"burakozkan138/questionanswerapi/config"
	"burakozkan138/questionanswerapi/internal/api/answers"
	"burakozkan138/questionanswerapi/internal/api/auth"
	"burakozkan138/questionanswerapi/internal/api/questions"
	"burakozkan138/questionanswerapi/internal/api/users"
	"burakozkan138/questionanswerapi/internal/middleware"
	"burakozkan138/questionanswerapi/internal/models"
	"log"
	"net/http"

	_ "burakozkan138/questionanswerapi/docs"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

var routes *http.ServeMux = http.NewServeMux()

func Rewrite() {
	routes = http.NewServeMux()
}

func InitializeRoutes() http.Handler {
	if config.SvConfig.SwaggerUI {
		routes.Handle("/swagger/", initilazeSwagger())
	}
	routes.Handle("/api/v1/", http.StripPrefix("/api/v1", groupRoutes())) // TODO: Add versioning to config

	stack := middleware.CreateStack(
		cors.AllowAll().Handler,
		middleware.ErrorHandler,
		middleware.Logging,
	)

	return stack(routes)
}

func groupRoutes() http.Handler {
	router := http.NewServeMux()

	router.Handle("/auth/", http.StripPrefix("/auth", initAuthRoutes()))
	router.Handle("/user/", http.StripPrefix("/user", initUserRoutes()))
	router.Handle("/question/", http.StripPrefix("/question", initQuestionRoutes()))
	router.Handle("/answer/", http.StripPrefix("/answer", initAnswerRoutes()))

	return router
}

func initAuthRoutes() http.Handler {
	router := http.NewServeMux()
	handler := auth.Handler{}

	router.Handle("POST /register", middleware.Validation(http.HandlerFunc(handler.Register), &models.RegisterValidation{}))
	router.Handle("POST /login", middleware.Validation(http.HandlerFunc(handler.Login), &models.LoginValidation{}))

	router.Handle("POST /forgotpassword", middleware.Validation(http.HandlerFunc(handler.ForgotPassword), &models.ForgotPasswordValidation{}))
	router.Handle("POST /resetpassword", middleware.Validation(http.HandlerFunc(handler.ResetPassword), &models.ResetPasswordValidation{}))
	return router
}

func initUserRoutes() http.Handler {
	router := http.NewServeMux()
	handler := users.Handler{}

	router.Handle("GET /profile", middleware.IsAuthenticated(http.HandlerFunc(handler.Profile)))
	router.Handle("PUT /edit", middleware.IsAuthenticated(middleware.Validation(http.HandlerFunc(handler.Edit), &models.EditUserValidation{})))
	router.Handle("POST /upload", middleware.IsAuthenticated(http.HandlerFunc(handler.UploadImage)))

	return router
}

func initQuestionRoutes() http.Handler {
	router := http.NewServeMux()
	handler := questions.Handler{}

	router.Handle("GET /", http.HandlerFunc(handler.GetAll))
	router.Handle("GET /{question_id}", http.HandlerFunc(handler.GetById))
	router.Handle("POST /ask", middleware.IsAuthenticated(middleware.Validation(http.HandlerFunc(handler.Ask), &models.CreateQuestionValidation{})))
	router.Handle("PUT /{question_id}/edit", middleware.IsAuthenticated(middleware.Validation(http.HandlerFunc(handler.Edit), &models.EditQuestionValidation{})))
	router.Handle("DELETE /{question_id}/delete", middleware.IsAuthenticated(http.HandlerFunc(handler.Delete)))
	router.Handle("PUT /{question_id}/like", middleware.IsAuthenticated(http.HandlerFunc(handler.Like)))

	return router
}

func initAnswerRoutes() http.Handler {
	router := http.NewServeMux()
	handler := answers.Handler{}

	router.Handle("POST /{question_id}", middleware.IsAuthenticated(middleware.Validation(http.HandlerFunc(handler.Create), &models.CreateAnswerValidation{})))
	router.Handle("GET /{question_id}", http.HandlerFunc(handler.GetAnswers))
	router.Handle("PUT /{answer_id}/like", middleware.IsAuthenticated(http.HandlerFunc(handler.Like)))
	router.Handle("PUT /{answe_id}/edit", middleware.IsAuthenticated(middleware.Validation(http.HandlerFunc(handler.Edit), &models.EditAnswerValidation{})))
	router.Handle("DELETE /{answer_id}/delete", middleware.IsAuthenticated(http.HandlerFunc(handler.Delete)))

	return router
}

func initilazeSwagger() http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	log.Println("Swagger UI is available at http://localhost:8080/swagger/")
	return router
}
