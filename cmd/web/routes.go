package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// ROUTES: 10
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	dynamicMiddleware := alice.New(app.sessionsManager.LoadAndSave)

	//Middelware to protect routs from unathenticated access
	protected := dynamicMiddleware.Append(app.requireAuthnticationMiddle)
	protectAdmin := dynamicMiddleware.Append(app.requireAdminAuthnticationMiddle)

	//Login Page Endpoints for public users
	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.login))
	router.Handler(http.MethodPost, "/", dynamicMiddleware.ThenFunc(app.loginSubmitted))

	//Sign up page endpoints for public users
	router.Handler(http.MethodGet, "/signup", dynamicMiddleware.ThenFunc(app.SignUpPage))
	router.Handler(http.MethodPost, "/signup", dynamicMiddleware.ThenFunc(app.signup))

	//Login Page Endpoints for admins
	router.Handler(http.MethodGet, "/adminLogin", dynamicMiddleware.ThenFunc(app.AdminLogin))
	router.Handler(http.MethodPost, "/adminLogin", dynamicMiddleware.ThenFunc(app.AdminLoginSubmitted))

	//Login Page Endpoints for admins
	router.Handler(http.MethodGet, "/adminSignup", dynamicMiddleware.ThenFunc(app.AdminSignUpPage))
	router.Handler(http.MethodPost, "/adminSignup", dynamicMiddleware.ThenFunc(app.Adminsignup))

	//UserDashboard
	router.Handler(http.MethodGet, "/userui", protected.ThenFunc(app.UserUI))

	//Submitting and displaying Affidavit form
	router.Handler(http.MethodGet, "/AffidavitForm", protected.ThenFunc(app.LoadForm))
	router.Handler(http.MethodPost, "/AffidavitForm", protected.ThenFunc(app.SubmittedAffidavit))

	//Views
	router.Handler(http.MethodGet, "/TableView", protectAdmin.ThenFunc(app.TableView))
	router.Handler(http.MethodGet, "/ArchiveView", protectAdmin.ThenFunc(app.ArchiveView))

	//Endpoints for viewing form with previous data and for updating data within that form
	router.Handler(http.MethodPost, "/updateForm", protected.ThenFunc(app.updateForm))
	router.Handler(http.MethodPost, "/updateFormSubmitted", protected.ThenFunc(app.updateFormQuery))

	//Endpoint for deleting form
	router.Handler(http.MethodPost, "/deleteForm", protected.ThenFunc(app.deleteForm))


	//Log out
	router.Handler(http.MethodGet, "/logout", protected.ThenFunc(app.Logout))

	//Yet to be implemented
	router.Handler(http.MethodPost, "/affidavitAccept", protected.ThenFunc(app.affidavitAccept))
	router.Handler(http.MethodPost, "/affidavitDeny", protected.ThenFunc(app.affidavitDeny))
	router.Handler(http.MethodPost, "/ArchiveForm", protected.ThenFunc(app.ArchiveForm))
	router.Handler(http.MethodPost, "/UnarchiveForm", protected.ThenFunc(app.UnarchiveForm))
	router.Handler(http.MethodPost, "/DownloadPDF", protected.ThenFunc(app.DownloadPDF))

	standardMiddleware := alice.New(app.recoverPanicMiddlewear, app.logRequestMiddlewear, securityHeadersMiddleware)
	return standardMiddleware.Then(router)
}
