package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/KiloProjects/kilonova/internal/util"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// MustBeVisitor is middleware to make sure the user creating the request is not authenticated
func (s *API) MustBeVisitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.base.IsAuthed(util.UserBrief(r)) {
			errorData(w, "You must not be logged in to do this", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MustBeAdmin is middleware to make sure the user creating the request is an admin
func (s *API) MustBeAdmin(next http.Handler) http.Handler {
	return s.MustBeAuthed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.IsAdmin(util.UserBrief(r)) {
			errorData(w, "You must be an admin to do this", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// MustBeAuthed is middleware to make sure the user creating the request is authenticated
func (s *API) MustBeAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.IsAuthed(util.UserBrief(r)) {
			errorData(w, "You must be authenticated to do this", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MustBeProposer is middleware to make sure the user creating the request is a proposer
func (s *API) MustBeProposer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.IsProposer(util.UserBrief(r)) {
			errorData(w, "You must be a proposer to do this", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SetupSession adds the user with the specified user ID to context
func (s *API) SetupSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.base.SessionUser(r.Context(), getAuthHeader(r))
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				zap.S().Warn(err)
			}
			next.ServeHTTP(w, r)
			return
		}
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.UserKey, user)))
	})
}

func (s *API) validateProblemEditor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.IsProblemEditor(util.UserBrief(r), util.Problem(r)) {
			errorData(w, "You must be authorized to update the problem", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *API) validateContestParticipant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.CanSubmitInContest(util.UserBrief(r), util.Contest(r)) {
			errorData(w, "You must be registered and during a contest to do this", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func (s *API) validateContestEditor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.IsContestEditor(util.UserBrief(r), util.Contest(r)) {
			errorData(w, "You must be authorized to update the contest", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func (s *API) validateContestVisible(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.base.IsContestVisible(util.UserBrief(r), util.Contest(r)) {
			errorData(w, "You are not allowed to access this contest", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// validateTestID pre-emptively returns if there isnt a valid test ID in the URL params
// Also, it fetches the test from the DB and makes sure it exists
// NOTE: This does not fetch the test data from disk
func (s *API) validateTestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testID, err := strconv.Atoi(chi.URLParam(r, "tID"))
		if err != nil {
			errorData(w, "invalid test ID", http.StatusBadRequest)
			return
		}
		test, err1 := s.base.Test(r.Context(), util.Problem(r).ID, testID)
		if err1 != nil {
			errorData(w, "test does not exist", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.TestKey, test)))
	})
}

func (s *API) validateAttachmentID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attID, err := strconv.Atoi(chi.URLParam(r, "aID"))
		if err != nil {
			errorData(w, "invalid attachment ID", http.StatusBadRequest)
			return
		}
		if util.Problem(r) == nil {
			zap.S().Fatal("Problem is not available")
			return
		}
		att, err1 := s.base.ProblemAttachment(r.Context(), util.Problem(r).ID, attID)
		if err1 != nil {
			errorData(w, "attachment does not exist", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.AttachmentKey, att)))
	})
}

func (s *API) validateSubmissionID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subID, err := strconv.Atoi(chi.URLParam(r, "submissionID"))
		if err != nil {
			errorData(w, "invalid attachment ID", http.StatusBadRequest)
			return
		}
		att, err1 := s.base.Submission(r.Context(), subID, util.UserBrief(r))
		if err1 != nil {
			errorData(w, "attachment does not exist", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.SubKey, att)))
	})
}

// validateProblemID pre-emptively returns if there isnt a valid problem ID in the URL params
// Also, it fetches the problem from the DB and makes sure it exists
func (s *API) validateProblemID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		problemID, err := strconv.Atoi(chi.URLParam(r, "problemID"))
		if err != nil {
			errorData(w, "invalid problem ID", http.StatusBadRequest)
			return
		}
		problem, err1 := s.base.Problem(r.Context(), problemID)
		if err1 != nil {
			errorData(w, "problem does not exist", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.ProblemKey, problem)))
	})
}

func (s *API) validateContestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contestID, err := strconv.Atoi(chi.URLParam(r, "contestID"))
		if err != nil {
			errorData(w, "invalid contest ID", http.StatusBadRequest)
			return
		}
		contest, err1 := s.base.Contest(r.Context(), contestID)
		if err1 != nil {
			errorData(w, "contest does not exist", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), util.ContestKey, contest)))
	})
}

func getAuthHeader(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "guest" {
		header = ""
	}
	return header
}
