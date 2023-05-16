package main

import (
	"errors"
	"fmt"

	"log"
	"net/http"
	"net/smtp"

	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/HipolitoBautista/internal/models"
)

// shared data store
var dataStore = struct {
	sync.RWMutex
	data map[string]int64
}{data: make(map[string]int64)}

func (app *application) UserUI(w http.ResponseWriter, r *http.Request) {

	UserIDTemp := app.sessionsManager.GetInt(r.Context(), "authenticatedUserID")
	FD, err := app.form.IndividualLoadData(UserIDTemp)
	Username, _, _ := app.publicuser.UserData(UserIDTemp)
	accepted, unverified, denied, total, err := app.form.NormalStats(UserIDTemp)

	if err != nil {
		fmt.Println(err)
	}

	flash := app.sessionsManager.PopString(r.Context(), "flash")
	fmt.Println(flash)
	data := &templateData{
		Form:       FD,
		Flash:      flash,
		Accepted:   accepted,
		Unverified: unverified,
		Denied:     denied,
		Username:   Username,
		Total:      total,
	}

	ts, err := template.ParseFiles("./ui/html/User.UI.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)

	fmt.Println(err)
}

func (app *application) TableView(w http.ResponseWriter, r *http.Request) {

	FD, err := app.form.LoadData()

	if err != nil {
		fmt.Println(err)
	}

	AdminIDTemp := app.sessionsManager.GetInt(r.Context(), "authenticatedAdminID")
	// getting the admin that's currently logged in
	Username, _, err := app.admin.AdminData(AdminIDTemp)
	//Getting the form stats
	accepted, unverified, denied, total, err := app.form.Stats()

	if err != nil {
		fmt.Println(err)
	}

	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Form:       FD,
		Flash:      flash,
		Accepted:   accepted,
		Unverified: unverified,
		Denied:     denied,
		Username:   Username,
		Total:      total,
	}

	ts, err := template.ParseFiles("./ui/html/Table.View.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	fmt.Println(err)
}

func (app *application) ArchiveView(w http.ResponseWriter, r *http.Request) {
	FD, err := app.form.LoadArchiveData()

	if err != nil {
		fmt.Println(err)
	}
	total, err := app.archive.Stats()
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Form:  FD,
		Flash: flash,
		Total: total,
	}

	ts, err := template.ParseFiles("./ui/html/Archive.View.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, data)
	fmt.Println(err)

}

func (app *application) LoadForm(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "./ui/html/Form.UI.tmpl", nil)

}

// Renders the Sign Up page template
func (app *application) SignUpPage(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	RenderTemplate(w, "./ui/html/signup.UI.tmpl", data)

}

// Handles the data given to us in the signup page
func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Lname := r.PostForm.Get("Lname")
	Fname := r.PostForm.Get("Fname")
	name := strings.ReplaceAll(Fname, " ", "") + " " + strings.ReplaceAll(Lname, " ", "")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")

	//write the data to the table
	err := app.publicuser.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			app.sessionsManager.Put(r.Context(), "flash", "Email Already Registered")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
		} else {
			app.sessionsManager.Put(r.Context(), "flash", "Oops! Try again later")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
		}
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Signup was successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Renders the Login page with it's data
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	RenderTemplate(w, "./ui/html/Login.UI.tmpl", data)
}

// Handles the data given to us in the login page
func (app *application) loginSubmitted(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")
	id, err := app.publicuser.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.sessionsManager.Put(r.Context(), "flash", "Invalid Credentials")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	// add user to the session cookie
	err = app.sessionsManager.RenewToken(r.Context())
	if err != nil {
		return
	}

	//add an authenticate
	app.sessionsManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/userui", http.StatusSeeOther)
}

// Loads admin login page
func (app *application) AdminLogin(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	RenderTemplate(w, "./ui/html/AdminLogin.UI.tmpl", data)
}

// Handles the data submitted
func (app *application) AdminLoginSubmitted(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")

	id, err := app.admin.AuthenticateAdmin(email, password)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.sessionsManager.Put(r.Context(), "flash", "Invalid Credentials")
			http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
		}
		return
	}

	// add user to the session cookie
	err = app.sessionsManager.RenewToken(r.Context())
	if err != nil {
		fmt.Println(err)
		return
	}

	//add an authenticate
	app.sessionsManager.Put(r.Context(), "authenticatedAdminID", id)
	app.sessionsManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/TableView", http.StatusSeeOther)
}

// Renders the Sign Up page template for admins
func (app *application) AdminSignUpPage(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	RenderTemplate(w, "./ui/html/adminsignup.UI.tmpl", data)

}

// Handles the data given to us in the signup page for admins
func (app *application) Adminsignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Lname := r.PostForm.Get("Lname")
	Fname := r.PostForm.Get("Fname")
	name := strings.ReplaceAll(Fname, " ", "") + " " + strings.ReplaceAll(Lname, " ", "")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")
	auth := r.PostForm.Get("AuthenticationCode")

	//write the data to the table

	err := app.admin.InsertAdmin(name, email, password, auth)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			app.sessionsManager.Put(r.Context(), "flash", "Email Already Registered")
			http.Redirect(w, r, "/AdminSignup", http.StatusSeeOther)
		} else {
			app.sessionsManager.Put(r.Context(), "flash", "Oops! Try again later")
			http.Redirect(w, r, "/AdminSignup", http.StatusSeeOther)
		}
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Signup was successful")
		http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
	}

}
func (app *application) SubmittedAffidavit(w http.ResponseWriter, r *http.Request) {

	//Code for INSERT
	w.Write([]byte("submitted"))

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//get the form data

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return

	}
	//collecting values from form
	name := r.PostForm.Get("full-name")
	other_name := r.PostForm.Get("other-names")
	name_change, _ := strconv.ParseBool(r.PostForm.Get("name-changed"))
	previous_name := r.PostForm.Get("previous-names")
	reason_for_change := r.PostForm.Get("reason-for-change")
	social_num, _ := strconv.ParseInt(r.PostForm.Get("SS-number"), 10, 64)
	social_date, _ := time.Parse("2006-01-02", r.PostForm.Get("SS-issue-date"))
	social_country := r.PostForm.Get("SS-country-value")
	passport_num, _ := strconv.ParseInt(r.PostForm.Get("PP-number"), 10, 64)
	passport_date, _ := time.Parse("2006-01-02", r.PostForm.Get("PP-issue-date"))
	passport_country := r.PostForm.Get("PP-country-value")
	NHI_num, _ := strconv.ParseInt(r.PostForm.Get("NHI-number"), 10, 64)
	NHI_date, _ := time.Parse("2006-01-02", r.PostForm.Get("NHI-issue-date"))
	NHI_country := r.PostForm.Get("NHI-country-value")
	dob, _ := time.Parse("2006-01-02", r.PostForm.Get("dateOfBirth"))
	pob := r.PostForm.Get("placeOfBirth")
	nationality := r.PostForm.Get("nationality")
	acq_nationality := r.PostForm.Get("acq-nationality")
	spouse := r.PostForm.Get("spouse-name")
	address := r.PostForm.Get("AF-address")
	phone, _ := strconv.ParseInt(r.PostForm.Get("AF-number"), 10, 64)
	fax, _ := strconv.ParseInt(r.PostForm.Get("Fax-Number"), 10, 64)
	email := r.PostForm.Get("email-address")

	UserID := app.sessionsManager.GetInt(r.Context(), "authenticatedUserID")
	//create instance of form model
	data := &models.Form{
		User_id:                  int32(UserID),
		Form_status:              "unverified",
		Archive_status:           false,
		Affiant_full_name:        name,
		Other_names:              other_name,
		Name_change_status:       name_change,
		Previous_name:            previous_name,
		Reason_for_Change:        reason_for_change,
		Social_security_num:      social_num,
		Social_security_date:     social_date,
		Social_security_country:  social_country,
		Passport_number:          passport_num,
		Passport_date:            passport_date,
		Passport_country:         passport_country,
		NHI_number:               NHI_num,
		NHI_date:                 NHI_date,
		NHI_country:              NHI_country,
		Dob:                      dob,
		Place_of_birth:           pob,
		Nationality:              nationality,
		Acquired_nationality:     acq_nationality,
		Spouse_name:              spouse,
		Affiants_address:         address,
		Residencial_phone_number: phone,
		Residenceial_fax_num:     fax,
		Residencial_email:        email,
	}
	//inserting data into DB using the insert method
	form_id, err := app.form.Insert(data)
	fmt.Println(form_id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/userui", http.StatusSeeOther)

}

func (app *application) updateForm(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("edit"), 10, 64)
	//form ID to load when wanting to edit the form (READING)
	form_id := ID
	FD, err := app.form.Read(form_id)
	CD, err := app.comments.GetComments(form_id)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	var AdminPerms string
	if app.sessionsManager.Exists(r.Context(), "authenticatedAdminID") {
		AdminPerms = "true"

	} else {
		AdminPerms = "false"
	}

	//an instance of templateData
	data := &templateData{
		Form:           FD,
		UserPermStatus: AdminPerms,
		Comments:       CD,
	} //this allows us to pass in mutliple data into the template

	ts, err := template.ParseFiles("./ui/html/form.show.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	//getting the form ID
	dataStore.Lock()
	dataStore.data["key"] = int64(form_id)
	dataStore.Unlock()

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

}

func (app *application) updateFormQuery(w http.ResponseWriter, r *http.Request) {
	//Code to UPDATE a form

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//get the form data

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return

	}
	//collecting values from form
	//form_id, _ := strconv.ParseInt(r.PostForm.Get("form-id"), 10, 64)
	name := r.PostForm.Get("full-name")
	other_name := r.PostForm.Get("other-names")
	name_change, _ := strconv.ParseBool(r.PostForm.Get("name-changed"))
	previous_name := r.PostForm.Get("previous-names")
	reason_for_change := r.PostForm.Get("reason-for-change")
	social_num, _ := strconv.ParseInt(r.PostForm.Get("SS-number"), 10, 64)
	social_date, _ := time.Parse("2006-01-02", r.PostForm.Get("SS-issue-date"))
	social_country := r.PostForm.Get("SS-country-value")
	passport_num, _ := strconv.ParseInt(r.PostForm.Get("PP-number"), 10, 64)
	passport_date, _ := time.Parse("2006-01-02", r.PostForm.Get("PP-issue-date"))
	passport_country := r.PostForm.Get("PP-country-value")
	NHI_num, _ := strconv.ParseInt(r.PostForm.Get("NHI-number"), 10, 64)
	NHI_date, _ := time.Parse("2006-01-02", r.PostForm.Get("NHI-issue-date"))
	NHI_country := r.PostForm.Get("NHI-country-value")
	dob, _ := time.Parse("2006-01-02", r.PostForm.Get("dateOfBirth"))
	pob := r.PostForm.Get("placeOfBirth")
	nationality := r.PostForm.Get("nationality")
	acq_nationality := r.PostForm.Get("acq-nationality")
	spouse := r.PostForm.Get("spouse-name")
	address := r.PostForm.Get("AF-address")
	phone, _ := strconv.ParseInt(r.PostForm.Get("AF-number"), 10, 64)
	fax, _ := strconv.ParseInt(r.PostForm.Get("Fax-Number"), 10, 64)
	email := r.PostForm.Get("email-address")
	dataStore.RLock()
	form_id := dataStore.data["key"]
	dataStore.RUnlock()

	data := &models.Form{
		Form_id:                  form_id,
		Form_status:              "unverified",
		Archive_status:           false,
		Affiant_full_name:        name,
		Other_names:              other_name,
		Name_change_status:       name_change,
		Previous_name:            previous_name,
		Reason_for_Change:        reason_for_change,
		Social_security_num:      social_num,
		Social_security_date:     social_date,
		Social_security_country:  social_country,
		Passport_number:          passport_num,
		Passport_date:            passport_date,
		Passport_country:         passport_country,
		NHI_number:               NHI_num,
		NHI_date:                 NHI_date,
		NHI_country:              NHI_country,
		Dob:                      dob,
		Place_of_birth:           pob,
		Nationality:              nationality,
		Acquired_nationality:     acq_nationality,
		Spouse_name:              spouse,
		Affiants_address:         address,
		Residencial_phone_number: phone,
		Residenceial_fax_num:     fax,
		Residencial_email:        email,
	}

	_, err = app.form.Update(data)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.sessionsManager.Put(r.Context(), "flash", "Form edited Successfully")

	if app.isAuthenticatedAdmin(r) == true {

		removeThis := "\\r+"

		comment1 := strings.TrimSpace(strings.Replace(r.PostForm.Get("1"), removeThis, "", -1))
		comment2 := strings.TrimSpace(strings.Replace(r.PostForm.Get("2"), removeThis, "", -1))
		comment3 := strings.TrimSpace(strings.Replace(r.PostForm.Get("3"), removeThis, "", -1))
		comment35 := strings.TrimSpace(strings.Replace(r.PostForm.Get("3.5"), removeThis, "", -1))
		comment4 := strings.TrimSpace(strings.Replace(r.PostForm.Get("4"), removeThis, "", -1))
		comment5 := strings.TrimSpace(strings.Replace(r.PostForm.Get("5"), removeThis, "", -1))
		comment6 := strings.TrimSpace(strings.Replace(r.PostForm.Get("6"), removeThis, "", -1))
		comment7 := strings.TrimSpace(strings.Replace(r.PostForm.Get("7"), removeThis, "", -1))
		comment8 := strings.TrimSpace(strings.Replace(r.PostForm.Get("8"), removeThis, "", -1))
		comment9 := strings.TrimSpace(strings.Replace(r.PostForm.Get("9"), removeThis, "", -1))
		comment10 := strings.TrimSpace(strings.Replace(r.PostForm.Get("10"), removeThis, "", -1))

		var ShouldInsert bool
		if comment1 != "" || comment2 != "" || comment3 != "" || comment35 != "" || comment4 != "" || comment5 != "" || comment6 != "" || comment7 != "" || comment8 != "" || comment9 != "" || comment10 != "" {
			ShouldInsert = true
		}

		AdminIDTemp := app.sessionsManager.GetInt(r.Context(), "authenticatedAdminID")
		AdminID := int64(AdminIDTemp)

		if ShouldInsert == true {
			commentData := &models.Comments{
				Admin_id:  AdminID,
				Form_id:   form_id,
				Comment1:  comment1,
				Comment2:  comment2,
				Comment3:  comment3,
				Comment35: comment35,
				Comment4:  comment4,
				Comment5:  comment5,
				Comment6:  comment6,
				Comment7:  comment7,
				Comment8:  comment8,
				Comment9:  comment9,
				Comment10: comment10,
			}

			fmt.Println(comment1)
			if app.comments.Exists(form_id) == true {

				app.comments.UpdateComments(commentData)
			} else {

				app.comments.InsertComment(commentData)
			}
		}
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/userui", http.StatusSeeOther)

	}

}

type json_data struct {
	formid string
}

func (app *application) deleteForm(w http.ResponseWriter, r *http.Request) {
	//Code to delete a form
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("del"), 10, 64)

	if err != nil {
		panic(err)
	}

	err = app.form.Delete(ID)

	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Deletion Failed, Try Again Later")
		http.Redirect(w, r, "/userui", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Successful Deletion")
		http.Redirect(w, r, "/userui", http.StatusSeeOther)
	}

}

func (app *application) affidavitAccept(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("accept"), 10, 64)

	if err != nil {
		panic(err)
	}

	err = app.form.AcceptForm(ID)

	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Action Failed, Try Again Later")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "form has been accepted")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	}

}

func (app *application) affidavitDeny(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("deny"), 10, 64)
	IDString := strconv.FormatInt(ID, 10)

	if err != nil {
		panic(err)
	}

	Owner, err := app.form.FormOwner(ID)
	if err != nil {
		fmt.Println(err)
	}
	_, email, err := app.publicuser.UserData(Owner)
	if err != nil {
		fmt.Println(err)
	}
	
	err = app.form.DenyForm(ID)
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Action Failed, Try Again Later")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Form has been denied ")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)

	}

	//Automatically sends email to inform the user the form has been denied

	from := "osippproject@gmail.com"
	password := "pzkrxnbaznlannih"
	to := email

	//giving it the server to authenticate details.
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	smtpServer := "smtp.gmail.com:587"
	conn, err := smtp.Dial(smtpServer)
	if err != nil {
		panic(err)
	}
	//closes connection once email is sent
	defer conn.Close()

	// message to be emailed
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Form Has been denied\r\n" +
		"\r\n" +
		"Hello,\r\n" +
		"We regret to inform you that your form (" + IDString + ") has been denied. \r\n" +
		"Your form may have been denied for several reasons.\r\n" +
		"In order to fix this issue log into your OSIPP account and open the denied form.\r\n" +
		"Within the form change the fields according to the comments our OSIPP Officers have left. Once the changes have been made resubmit the form\r\n")

	// Send the email.
	if err := smtp.SendMail(smtpServer, auth, from, []string{to}, msg); err != nil {
		panic(err)
	}

}

func (app *application) ArchiveForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("archive"), 10, 64)

	if err != nil {
		panic(err)
	}

	err = app.archive.ArchiveForm(ID)

	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Archive Failed, Try Again Later")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Form has been Archived ")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	}

}

func (app *application) UnarchiveForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering Endpoint")
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("unarchive"), 10, 64)

	if err != nil {
		panic(err)
	}

	err = app.archive.UnArchiveForm(ID)
	fmt.Println(err)
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Archive Failed, Try Again Later")
		http.Redirect(w, r, "/ArchiveView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Form has been Archived ")
		http.Redirect(w, r, "/ArchiveView", http.StatusSeeOther)
	}
}

func (app *application) DownloadPDF(w http.ResponseWriter, r *http.Request) {

}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {

	if app.isAuthenticatedAdmin(r) == true {
		app.sessionsManager.PopString(r.Context(), "authenticatedUserID")
		app.sessionsManager.PopString(r.Context(), "authenticatedAdminID")
		app.sessionsManager.Put(r.Context(), "flash", "Successfully Logged out ")
		http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
	} else {
		app.sessionsManager.PopString(r.Context(), "authenticatedUserID")
		app.sessionsManager.Put(r.Context(), "flash", "Successfully Logged out ")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}
