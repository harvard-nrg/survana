package study

import (
	"github.com/vpetrov/perfect"
	"github.com/vpetrov/perfect/auth"
	"github.com/vpetrov/perfect/orm"
	"log"
	"net/http"
	"neuroinformatics.harvard.edu/survana"
	"strings"
	"time"
)

type SubjectIdStrategy struct {
}

func NewSubjectIdStrategy(config *Config) *SubjectIdStrategy {
	return &SubjectIdStrategy{}
}

func NewSubjectIdStrategyFunc(config *Config) auth.Strategy {
	return NewSubjectIdStrategy(config)
}

func (sid *SubjectIdStrategy) Attach(module *perfect.Module) {
	module.Get("/login", perfect.NotLoggedIn(sid.LoginPage))
}

func (sid *SubjectIdStrategy) LoginPage(w http.ResponseWriter, r *perfect.Request) {

	//if we have both the study_id and the subject_id in the URL, perform an internal
	//redirect to .Login
	if strings.Contains(r.URL.RawQuery, "/") {
		sid.Login(w, r)
		return
	}

	data := &struct {
		Study_id string
		Msg      string
	}{}

	data.Study_id = r.URL.RawQuery
	if len(data.Study_id) == 0 {
		data.Msg = "This study does not exist. Please check to make sure you've used the correct link."
	}

	r.Module.RenderTemplate(w, r, "auth/subject_id/login", data)
}

func (sid *SubjectIdStrategy) RegistrationPage(w http.ResponseWriter, r *perfect.Request) {
	perfect.NotFound(w)
}

func (sid *SubjectIdStrategy) Login(w http.ResponseWriter, r *perfect.Request) (profile_id *string, err error) {
	var db = r.Module.Db

	//get the session
	session, err := r.Session()
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//if the user is already authenticated, redirect to home
	if session != nil && *session.Authenticated {
		perfect.Redirect(w, r, "/")
		return
	}

	params := strings.SplitN(r.URL.RawQuery, "/", 2)
	study_id := params[0]
	subject_id := strings.ToUpper(params[1])

	//no study id?
	if len(study_id) == 0 {
		perfect.Redirect(w, r, "/")
		return
	}

	if len(subject_id) == 0 {
		perfect.Redirect(w, r, auth.LOGIN_PATH+"?"+study_id)
		return
	}

	log.Println("study_id=", study_id, "subject_id=", subject_id)

	study := &survana.Study{Id: &study_id}
	err = db.Find(study)
	if err != nil {
		if err == orm.ErrNotFound {
			perfect.NotFound(w)
		} else {
			perfect.Error(w, r, err)
		}
		return
	}

	if study.Subjects == nil || *study.Subjects == nil {
		perfect.NotFound(w)
		return
	}

	data := &struct {
		Study_id string
		Msg      string
	}{
		Study_id: study_id,
	}

	enabled, ok := (*study.Subjects)[subject_id]
	if !ok {
		data.Msg = "We were unable to find this participant ID."
		r.Module.RenderTemplate(w, r, "auth/subject_id/login", data)
		return
	}

	if !enabled {
		data.Msg = "We were unable to find this participant ID."
		r.Module.RenderTemplate(w, r, "auth/subject_id/login", data)
		return
	}

	//mark the session as authenticated
	session.Authenticated = orm.Bool(true)

	//regenerate the session Id
	session.Id = orm.String(r.Module.Db.UniqueId())
	(*session.Values)["study_id"] = study_id
	(*session.Values)["subject_id"] = subject_id

	// update the session
	err = db.Save(session)
	if err != nil {
		perfect.Error(w, r, err)
		return
	}

	//set the cookie and make it valid for a month
	http.SetCookie(w, &http.Cookie{
		Name:     perfect.SESSION_ID,
		Value:    *session.Id,
		Path:     r.Module.MountPoint,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   true,
		HttpOnly: true,
	})

	perfect.Redirect(w, r, "?"+r.URL.RawQuery)

	return nil, nil
}

func (sid *SubjectIdStrategy) Register(w http.ResponseWriter, r *perfect.Request) {
	perfect.NotFound(w)
}

func (sid *SubjectIdStrategy) Logout(w http.ResponseWriter, r *perfect.Request) {
	perfect.NotFound(w)
}
