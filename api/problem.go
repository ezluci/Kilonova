package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"

	_ "embed"

	"github.com/KiloProjects/kilonova"
	"github.com/KiloProjects/kilonova/internal/util"
	"github.com/KiloProjects/kilonova/sudoapi"
	"go.uber.org/zap"
)

func (s *API) maxScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args struct {
		UserID int
	}
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	if args.UserID <= 0 {
		if util.UserBrief(r) == nil {
			errorData(w, "No user specified", 400)
			return
		}
		args.UserID = util.UserBrief(r).ID
	}

	returnData(w, s.base.MaxScore(r.Context(), args.UserID, util.Problem(r).ID))
}

func (s *API) problemStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := s.base.ProblemStatistics(r.Context(), util.Problem(r), util.UserBrief(r))
	if err != nil {
		err.WriteError(w)
		return
	}
	returnData(w, stats)
}

type scoreBreakdownRet struct {
	MaxScore int                           `json:"max_score"`
	Problem  *kilonova.Problem             `json:"problem"`
	Subtasks []*kilonova.SubmissionSubTask `json:"subtasks"`

	// ProblemEditor is true only if the request author is public
	// It does not take into consideration if the supplied user is the problem editor
	ProblemEditor bool `json:"problem_editor"`
	// Subtests are arranged from submission subtasks so something legible can be rebuilt to show more information on the subtasks
	Subtests []*kilonova.SubTest `json:"subtests"`
}

func (s *API) maxScoreBreakdown(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args struct {
		UserID int

		ContestID *int
	}
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	// This endpoint may leak stuff that shouldn't be generally seen (like in contests), so restrict this option to editors only
	// It isn't used anywhere right now, but it might be useful in the future
	if !s.base.IsProblemEditor(util.UserBrief(r), util.Problem(r)) {
		args.UserID = -1
	}
	if args.UserID <= 0 {
		if util.UserBrief(r) == nil {
			errorData(w, "No user specified", 400)
			return
		}
		args.UserID = util.UserBrief(r).ID
	}

	maxScore := -1
	if args.ContestID == nil {
		maxScore = s.base.MaxScore(r.Context(), args.UserID, util.Problem(r).ID)
	} else {
		maxScore = s.base.ContestMaxScore(r.Context(), args.UserID, util.Problem(r).ID, *args.ContestID)
	}

	switch util.Problem(r).ScoringStrategy {
	case kilonova.ScoringTypeMaxSub:
		id, err := s.base.MaxScoreSubID(r.Context(), args.UserID, util.Problem(r).ID)
		if err != nil {
			err.WriteError(w)
			return
		}
		if id <= 0 {
			returnData(w, scoreBreakdownRet{
				MaxScore: maxScore,
				Problem:  util.Problem(r),
				Subtasks: []*kilonova.SubmissionSubTask{},
				Subtests: []*kilonova.SubTest{},

				ProblemEditor: s.base.IsProblemEditor(util.UserBrief(r), util.Problem(r)),
			})
			return
		}
		sub, err := s.base.Submission(r.Context(), id, util.UserBrief(r))
		if err != nil {
			err.WriteError(w)
			return
		}

		returnData(w, scoreBreakdownRet{
			MaxScore: maxScore,
			Problem:  util.Problem(r),
			Subtasks: sub.SubTasks,
			Subtests: sub.SubTests,

			ProblemEditor: s.base.IsProblemEditor(util.UserBrief(r), util.Problem(r)),
		})
	case kilonova.ScoringTypeSumSubtasks:
		stks, err := s.base.MaximumScoreSubTasks(r.Context(), util.Problem(r).ID, args.UserID, args.ContestID)
		if err != nil {
			err.WriteError(w)
			return
		}

		tests, err := s.base.MaximumScoreSubTaskTests(r.Context(), util.Problem(r).ID, args.UserID, args.ContestID)
		if err != nil {
			err.WriteError(w)
			return
		}

		returnData(w, scoreBreakdownRet{
			MaxScore: maxScore,
			Problem:  util.Problem(r),
			Subtasks: stks,
			Subtests: tests,

			ProblemEditor: s.base.IsProblemEditor(util.UserBrief(r), util.Problem(r)),
		})
	default:
		zap.S().Warn("Unknown problem scoring type")
		errorData(w, "Unknown problem scoring type", 500)
	}

}

func (s *API) deleteProblem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if err := s.base.DeleteProblem(context.WithoutCancel(r.Context()), util.Problem(r)); err != nil {
		errorData(w, err, 500)
		return
	}
	returnData(w, "Deleted problem")
}

var (
	//go:embed templData/default_en_statement.md
	enPbStatementStr string
	//go:embed templData/default_ro_statement.md
	roPbStatementStr string

	defaultEnProblemStatement = template.Must(template.New("enStmt").Parse(enPbStatementStr))
	defaultRoProblemStatement = template.Must(template.New("enStmt").Parse(roPbStatementStr))
)

func (s *API) initProblem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args struct {
		Title        string `json:"title"`
		ConsoleInput bool   `json:"consoleInput"`

		StatementLang *string `json:"statementLang"`
	}
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	// Do the check before problem creation because it'd be awkward to create the problem and then show the error
	if args.StatementLang != nil && !(*args.StatementLang == "en" || *args.StatementLang == "ro" || *args.StatementLang == "") {
		errorData(w, "Invalid initial statement language", 400)
		return
	}

	pb, err := s.base.CreateProblem(r.Context(), args.Title, util.UserBrief(r), args.ConsoleInput)
	if err != nil {
		err.WriteError(w)
		return
	}

	if args.StatementLang != nil && *args.StatementLang != "" {
		var attTempl *template.Template
		if *args.StatementLang == "en" {
			attTempl = defaultEnProblemStatement
		} else if *args.StatementLang == "ro" {
			attTempl = defaultRoProblemStatement
		} else {
			zap.S().Warn("How did we get here? %q", *args.StatementLang)
			returnData(w, pb.ID)
			return
		}
		inFile := "stdin"
		outFile := "stdout"
		if !args.ConsoleInput {
			inFile = pb.TestName + ".in"
			outFile = pb.TestName + ".out"
		}
		var buf bytes.Buffer
		if err := attTempl.Execute(&buf, struct {
			InputFile  string
			OutputFile string
		}{InputFile: inFile, OutputFile: outFile}); err != nil {
			zap.S().Warnf("Template rendering error: %v", err)
		}
		if err := s.base.CreateAttachment(r.Context(), &kilonova.Attachment{
			Visible: false,
			Private: false,
			Exec:    false,
			Name:    fmt.Sprintf("statement-%s.md", *args.StatementLang),
		}, pb.ID, &buf, &util.UserBrief(r).ID); err != nil {
			zap.S().Warn(err)
		}
	}

	returnData(w, pb.ID)
}

func (s *API) getProblems(w http.ResponseWriter, r *http.Request) {
	var args kilonova.ProblemFilter
	if err := parseJsonBody(r, &args); err != nil {
		err.WriteError(w)
		return
	}

	args.Look = true
	args.LookingUser = util.UserBrief(r)

	problems, err := s.base.Problems(r.Context(), args)
	if err != nil {
		err.WriteError(w)
		return
	}
	returnData(w, problems)
}

func (s *API) searchProblems(w http.ResponseWriter, r *http.Request) {
	var args kilonova.ProblemFilter
	if err := parseJsonBody(r, &args); err != nil {
		err.WriteError(w)
		return
	}

	args.Look = true
	args.LookingUser = util.UserBrief(r)

	if args.Limit == 0 || args.Limit > 50 {
		args.Limit = 50
	}

	problems, cnt, err := s.base.SearchProblems(r.Context(), args, util.UserBrief(r))
	if err != nil {
		err.WriteError(w)
		return
	}
	returnData(w, struct {
		Problems []*sudoapi.FullProblem `json:"problems"`

		Count int `json:"count"`
	}{Problems: problems, Count: cnt})
}

func (s *API) updateProblem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args kilonova.ProblemUpdate
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	if err := s.base.UpdateProblem(r.Context(), util.Problem(r).ID, args, util.UserBrief(r)); err != nil {
		err.WriteError(w)
		return
	}

	returnData(w, "Updated problem")
}

func (s *API) addProblemEditor(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args struct {
		Username string `json:"username"`
	}
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	user, err := s.base.UserBriefByName(r.Context(), args.Username)
	if err != nil {
		err.WriteError(w)
		return
	}

	if err := s.base.AddProblemEditor(r.Context(), util.Problem(r).ID, user.ID); err != nil {
		err.WriteError(w)
		return
	}

	returnData(w, "Added problem editor")
}

func (s *API) addProblemViewer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args struct {
		Username string `json:"username"`
	}
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	user, err := s.base.UserBriefByName(r.Context(), args.Username)
	if err != nil {
		err.WriteError(w)
		return
	}

	if user.ID == util.UserBrief(r).ID {
		errorData(w, "You can't demote yourself to viewer rank!", 400)
		return
	}

	if err := s.base.AddProblemViewer(r.Context(), util.Problem(r).ID, user.ID); err != nil {
		err.WriteError(w)
		return
	}

	returnData(w, "Added problem viewer")
}

func (s *API) stripProblemAccess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var args struct {
		UserID int `json:"user_id"`
	}
	if err := decoder.Decode(&args, r.Form); err != nil {
		errorData(w, err, 400)
		return
	}

	if args.UserID == util.UserBrief(r).ID {
		errorData(w, "You can't strip your own access!", 400)
		return
	}

	if err := s.base.StripProblemAccess(r.Context(), util.Problem(r).ID, args.UserID); err != nil {
		err.WriteError(w)
		return
	}

	returnData(w, "Stripped problem access")
}

func (s *API) getProblemAccessControl(w http.ResponseWriter, r *http.Request) {
	editors, err := s.base.ProblemEditors(r.Context(), util.Problem(r).ID)
	if err != nil {
		err.WriteError(w)
		return
	}

	viewers, err := s.base.ProblemViewers(r.Context(), util.Problem(r).ID)
	if err != nil {
		err.WriteError(w)
		return
	}

	returnData(w, struct {
		Editors []*kilonova.UserBrief `json:"editors"`
		Viewers []*kilonova.UserBrief `json:"viewers"`
	}{
		Editors: editors,
		Viewers: viewers,
	})
}
