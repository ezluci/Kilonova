// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.adminsStmt, err = db.PrepareContext(ctx, admins); err != nil {
		return nil, fmt.Errorf("error preparing query Admins: %w", err)
	}
	if q.biggestVIDStmt, err = db.PrepareContext(ctx, biggestVID); err != nil {
		return nil, fmt.Errorf("error preparing query BiggestVID: %w", err)
	}
	if q.countProblemsStmt, err = db.PrepareContext(ctx, countProblems); err != nil {
		return nil, fmt.Errorf("error preparing query CountProblems: %w", err)
	}
	if q.countUsersStmt, err = db.PrepareContext(ctx, countUsers); err != nil {
		return nil, fmt.Errorf("error preparing query CountUsers: %w", err)
	}
	if q.createProblemStmt, err = db.PrepareContext(ctx, createProblem); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProblem: %w", err)
	}
	if q.createSubTestStmt, err = db.PrepareContext(ctx, createSubTest); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSubTest: %w", err)
	}
	if q.createSubmissionStmt, err = db.PrepareContext(ctx, createSubmission); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSubmission: %w", err)
	}
	if q.createTestStmt, err = db.PrepareContext(ctx, createTest); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTest: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.maxScoreStmt, err = db.PrepareContext(ctx, maxScore); err != nil {
		return nil, fmt.Errorf("error preparing query MaxScore: %w", err)
	}
	if q.problemStmt, err = db.PrepareContext(ctx, problem); err != nil {
		return nil, fmt.Errorf("error preparing query Problem: %w", err)
	}
	if q.problemTestsStmt, err = db.PrepareContext(ctx, problemTests); err != nil {
		return nil, fmt.Errorf("error preparing query ProblemTests: %w", err)
	}
	if q.problemsStmt, err = db.PrepareContext(ctx, problems); err != nil {
		return nil, fmt.Errorf("error preparing query Problems: %w", err)
	}
	if q.proposersStmt, err = db.PrepareContext(ctx, proposers); err != nil {
		return nil, fmt.Errorf("error preparing query Proposers: %w", err)
	}
	if q.purgePbTestsStmt, err = db.PrepareContext(ctx, purgePbTests); err != nil {
		return nil, fmt.Errorf("error preparing query PurgePbTests: %w", err)
	}
	if q.setAdminStmt, err = db.PrepareContext(ctx, setAdmin); err != nil {
		return nil, fmt.Errorf("error preparing query SetAdmin: %w", err)
	}
	if q.setBioStmt, err = db.PrepareContext(ctx, setBio); err != nil {
		return nil, fmt.Errorf("error preparing query SetBio: %w", err)
	}
	if q.setCompilationStmt, err = db.PrepareContext(ctx, setCompilation); err != nil {
		return nil, fmt.Errorf("error preparing query SetCompilation: %w", err)
	}
	if q.setConsoleInputStmt, err = db.PrepareContext(ctx, setConsoleInput); err != nil {
		return nil, fmt.Errorf("error preparing query SetConsoleInput: %w", err)
	}
	if q.setEmailStmt, err = db.PrepareContext(ctx, setEmail); err != nil {
		return nil, fmt.Errorf("error preparing query SetEmail: %w", err)
	}
	if q.setLimitsStmt, err = db.PrepareContext(ctx, setLimits); err != nil {
		return nil, fmt.Errorf("error preparing query SetLimits: %w", err)
	}
	if q.setMemoryLimitStmt, err = db.PrepareContext(ctx, setMemoryLimit); err != nil {
		return nil, fmt.Errorf("error preparing query SetMemoryLimit: %w", err)
	}
	if q.setPbTestScoreStmt, err = db.PrepareContext(ctx, setPbTestScore); err != nil {
		return nil, fmt.Errorf("error preparing query SetPbTestScore: %w", err)
	}
	if q.setPbTestVisibleIDStmt, err = db.PrepareContext(ctx, setPbTestVisibleID); err != nil {
		return nil, fmt.Errorf("error preparing query SetPbTestVisibleID: %w", err)
	}
	if q.setProblemDescriptionStmt, err = db.PrepareContext(ctx, setProblemDescription); err != nil {
		return nil, fmt.Errorf("error preparing query SetProblemDescription: %w", err)
	}
	if q.setProblemNameStmt, err = db.PrepareContext(ctx, setProblemName); err != nil {
		return nil, fmt.Errorf("error preparing query SetProblemName: %w", err)
	}
	if q.setProblemVisibilityStmt, err = db.PrepareContext(ctx, setProblemVisibility); err != nil {
		return nil, fmt.Errorf("error preparing query SetProblemVisibility: %w", err)
	}
	if q.setProposerStmt, err = db.PrepareContext(ctx, setProposer); err != nil {
		return nil, fmt.Errorf("error preparing query SetProposer: %w", err)
	}
	if q.setStackLimitStmt, err = db.PrepareContext(ctx, setStackLimit); err != nil {
		return nil, fmt.Errorf("error preparing query SetStackLimit: %w", err)
	}
	if q.setSubmissionStatusStmt, err = db.PrepareContext(ctx, setSubmissionStatus); err != nil {
		return nil, fmt.Errorf("error preparing query SetSubmissionStatus: %w", err)
	}
	if q.setSubmissionTestStmt, err = db.PrepareContext(ctx, setSubmissionTest); err != nil {
		return nil, fmt.Errorf("error preparing query SetSubmissionTest: %w", err)
	}
	if q.setSubmissionVisibilityStmt, err = db.PrepareContext(ctx, setSubmissionVisibility); err != nil {
		return nil, fmt.Errorf("error preparing query SetSubmissionVisibility: %w", err)
	}
	if q.setTestNameStmt, err = db.PrepareContext(ctx, setTestName); err != nil {
		return nil, fmt.Errorf("error preparing query SetTestName: %w", err)
	}
	if q.setTimeLimitStmt, err = db.PrepareContext(ctx, setTimeLimit); err != nil {
		return nil, fmt.Errorf("error preparing query SetTimeLimit: %w", err)
	}
	if q.setVisibleIDStmt, err = db.PrepareContext(ctx, setVisibleID); err != nil {
		return nil, fmt.Errorf("error preparing query SetVisibleID: %w", err)
	}
	if q.subTestsStmt, err = db.PrepareContext(ctx, subTests); err != nil {
		return nil, fmt.Errorf("error preparing query SubTests: %w", err)
	}
	if q.submissionStmt, err = db.PrepareContext(ctx, submission); err != nil {
		return nil, fmt.Errorf("error preparing query Submission: %w", err)
	}
	if q.submissionsStmt, err = db.PrepareContext(ctx, submissions); err != nil {
		return nil, fmt.Errorf("error preparing query Submissions: %w", err)
	}
	if q.testStmt, err = db.PrepareContext(ctx, test); err != nil {
		return nil, fmt.Errorf("error preparing query Test: %w", err)
	}
	if q.testVisibleIDStmt, err = db.PrepareContext(ctx, testVisibleID); err != nil {
		return nil, fmt.Errorf("error preparing query TestVisibleID: %w", err)
	}
	if q.userStmt, err = db.PrepareContext(ctx, user); err != nil {
		return nil, fmt.Errorf("error preparing query User: %w", err)
	}
	if q.userByEmailStmt, err = db.PrepareContext(ctx, userByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query UserByEmail: %w", err)
	}
	if q.userByNameStmt, err = db.PrepareContext(ctx, userByName); err != nil {
		return nil, fmt.Errorf("error preparing query UserByName: %w", err)
	}
	if q.userProblemSubmissionsStmt, err = db.PrepareContext(ctx, userProblemSubmissions); err != nil {
		return nil, fmt.Errorf("error preparing query UserProblemSubmissions: %w", err)
	}
	if q.usersStmt, err = db.PrepareContext(ctx, users); err != nil {
		return nil, fmt.Errorf("error preparing query Users: %w", err)
	}
	if q.visibleProblemsStmt, err = db.PrepareContext(ctx, visibleProblems); err != nil {
		return nil, fmt.Errorf("error preparing query VisibleProblems: %w", err)
	}
	if q.waitingSubmissionsStmt, err = db.PrepareContext(ctx, waitingSubmissions); err != nil {
		return nil, fmt.Errorf("error preparing query WaitingSubmissions: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.adminsStmt != nil {
		if cerr := q.adminsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing adminsStmt: %w", cerr)
		}
	}
	if q.biggestVIDStmt != nil {
		if cerr := q.biggestVIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing biggestVIDStmt: %w", cerr)
		}
	}
	if q.countProblemsStmt != nil {
		if cerr := q.countProblemsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countProblemsStmt: %w", cerr)
		}
	}
	if q.countUsersStmt != nil {
		if cerr := q.countUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countUsersStmt: %w", cerr)
		}
	}
	if q.createProblemStmt != nil {
		if cerr := q.createProblemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createProblemStmt: %w", cerr)
		}
	}
	if q.createSubTestStmt != nil {
		if cerr := q.createSubTestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSubTestStmt: %w", cerr)
		}
	}
	if q.createSubmissionStmt != nil {
		if cerr := q.createSubmissionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSubmissionStmt: %w", cerr)
		}
	}
	if q.createTestStmt != nil {
		if cerr := q.createTestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTestStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.maxScoreStmt != nil {
		if cerr := q.maxScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing maxScoreStmt: %w", cerr)
		}
	}
	if q.problemStmt != nil {
		if cerr := q.problemStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing problemStmt: %w", cerr)
		}
	}
	if q.problemTestsStmt != nil {
		if cerr := q.problemTestsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing problemTestsStmt: %w", cerr)
		}
	}
	if q.problemsStmt != nil {
		if cerr := q.problemsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing problemsStmt: %w", cerr)
		}
	}
	if q.proposersStmt != nil {
		if cerr := q.proposersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing proposersStmt: %w", cerr)
		}
	}
	if q.purgePbTestsStmt != nil {
		if cerr := q.purgePbTestsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing purgePbTestsStmt: %w", cerr)
		}
	}
	if q.setAdminStmt != nil {
		if cerr := q.setAdminStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setAdminStmt: %w", cerr)
		}
	}
	if q.setBioStmt != nil {
		if cerr := q.setBioStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setBioStmt: %w", cerr)
		}
	}
	if q.setCompilationStmt != nil {
		if cerr := q.setCompilationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setCompilationStmt: %w", cerr)
		}
	}
	if q.setConsoleInputStmt != nil {
		if cerr := q.setConsoleInputStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setConsoleInputStmt: %w", cerr)
		}
	}
	if q.setEmailStmt != nil {
		if cerr := q.setEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setEmailStmt: %w", cerr)
		}
	}
	if q.setLimitsStmt != nil {
		if cerr := q.setLimitsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setLimitsStmt: %w", cerr)
		}
	}
	if q.setMemoryLimitStmt != nil {
		if cerr := q.setMemoryLimitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setMemoryLimitStmt: %w", cerr)
		}
	}
	if q.setPbTestScoreStmt != nil {
		if cerr := q.setPbTestScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setPbTestScoreStmt: %w", cerr)
		}
	}
	if q.setPbTestVisibleIDStmt != nil {
		if cerr := q.setPbTestVisibleIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setPbTestVisibleIDStmt: %w", cerr)
		}
	}
	if q.setProblemDescriptionStmt != nil {
		if cerr := q.setProblemDescriptionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setProblemDescriptionStmt: %w", cerr)
		}
	}
	if q.setProblemNameStmt != nil {
		if cerr := q.setProblemNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setProblemNameStmt: %w", cerr)
		}
	}
	if q.setProblemVisibilityStmt != nil {
		if cerr := q.setProblemVisibilityStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setProblemVisibilityStmt: %w", cerr)
		}
	}
	if q.setProposerStmt != nil {
		if cerr := q.setProposerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setProposerStmt: %w", cerr)
		}
	}
	if q.setStackLimitStmt != nil {
		if cerr := q.setStackLimitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setStackLimitStmt: %w", cerr)
		}
	}
	if q.setSubmissionStatusStmt != nil {
		if cerr := q.setSubmissionStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setSubmissionStatusStmt: %w", cerr)
		}
	}
	if q.setSubmissionTestStmt != nil {
		if cerr := q.setSubmissionTestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setSubmissionTestStmt: %w", cerr)
		}
	}
	if q.setSubmissionVisibilityStmt != nil {
		if cerr := q.setSubmissionVisibilityStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setSubmissionVisibilityStmt: %w", cerr)
		}
	}
	if q.setTestNameStmt != nil {
		if cerr := q.setTestNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setTestNameStmt: %w", cerr)
		}
	}
	if q.setTimeLimitStmt != nil {
		if cerr := q.setTimeLimitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setTimeLimitStmt: %w", cerr)
		}
	}
	if q.setVisibleIDStmt != nil {
		if cerr := q.setVisibleIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setVisibleIDStmt: %w", cerr)
		}
	}
	if q.subTestsStmt != nil {
		if cerr := q.subTestsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing subTestsStmt: %w", cerr)
		}
	}
	if q.submissionStmt != nil {
		if cerr := q.submissionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing submissionStmt: %w", cerr)
		}
	}
	if q.submissionsStmt != nil {
		if cerr := q.submissionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing submissionsStmt: %w", cerr)
		}
	}
	if q.testStmt != nil {
		if cerr := q.testStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing testStmt: %w", cerr)
		}
	}
	if q.testVisibleIDStmt != nil {
		if cerr := q.testVisibleIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing testVisibleIDStmt: %w", cerr)
		}
	}
	if q.userStmt != nil {
		if cerr := q.userStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userStmt: %w", cerr)
		}
	}
	if q.userByEmailStmt != nil {
		if cerr := q.userByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userByEmailStmt: %w", cerr)
		}
	}
	if q.userByNameStmt != nil {
		if cerr := q.userByNameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userByNameStmt: %w", cerr)
		}
	}
	if q.userProblemSubmissionsStmt != nil {
		if cerr := q.userProblemSubmissionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing userProblemSubmissionsStmt: %w", cerr)
		}
	}
	if q.usersStmt != nil {
		if cerr := q.usersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing usersStmt: %w", cerr)
		}
	}
	if q.visibleProblemsStmt != nil {
		if cerr := q.visibleProblemsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing visibleProblemsStmt: %w", cerr)
		}
	}
	if q.waitingSubmissionsStmt != nil {
		if cerr := q.waitingSubmissionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing waitingSubmissionsStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                          DBTX
	tx                          *sql.Tx
	adminsStmt                  *sql.Stmt
	biggestVIDStmt              *sql.Stmt
	countProblemsStmt           *sql.Stmt
	countUsersStmt              *sql.Stmt
	createProblemStmt           *sql.Stmt
	createSubTestStmt           *sql.Stmt
	createSubmissionStmt        *sql.Stmt
	createTestStmt              *sql.Stmt
	createUserStmt              *sql.Stmt
	maxScoreStmt                *sql.Stmt
	problemStmt                 *sql.Stmt
	problemTestsStmt            *sql.Stmt
	problemsStmt                *sql.Stmt
	proposersStmt               *sql.Stmt
	purgePbTestsStmt            *sql.Stmt
	setAdminStmt                *sql.Stmt
	setBioStmt                  *sql.Stmt
	setCompilationStmt          *sql.Stmt
	setConsoleInputStmt         *sql.Stmt
	setEmailStmt                *sql.Stmt
	setLimitsStmt               *sql.Stmt
	setMemoryLimitStmt          *sql.Stmt
	setPbTestScoreStmt          *sql.Stmt
	setPbTestVisibleIDStmt      *sql.Stmt
	setProblemDescriptionStmt   *sql.Stmt
	setProblemNameStmt          *sql.Stmt
	setProblemVisibilityStmt    *sql.Stmt
	setProposerStmt             *sql.Stmt
	setStackLimitStmt           *sql.Stmt
	setSubmissionStatusStmt     *sql.Stmt
	setSubmissionTestStmt       *sql.Stmt
	setSubmissionVisibilityStmt *sql.Stmt
	setTestNameStmt             *sql.Stmt
	setTimeLimitStmt            *sql.Stmt
	setVisibleIDStmt            *sql.Stmt
	subTestsStmt                *sql.Stmt
	submissionStmt              *sql.Stmt
	submissionsStmt             *sql.Stmt
	testStmt                    *sql.Stmt
	testVisibleIDStmt           *sql.Stmt
	userStmt                    *sql.Stmt
	userByEmailStmt             *sql.Stmt
	userByNameStmt              *sql.Stmt
	userProblemSubmissionsStmt  *sql.Stmt
	usersStmt                   *sql.Stmt
	visibleProblemsStmt         *sql.Stmt
	waitingSubmissionsStmt      *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                          tx,
		tx:                          tx,
		adminsStmt:                  q.adminsStmt,
		biggestVIDStmt:              q.biggestVIDStmt,
		countProblemsStmt:           q.countProblemsStmt,
		countUsersStmt:              q.countUsersStmt,
		createProblemStmt:           q.createProblemStmt,
		createSubTestStmt:           q.createSubTestStmt,
		createSubmissionStmt:        q.createSubmissionStmt,
		createTestStmt:              q.createTestStmt,
		createUserStmt:              q.createUserStmt,
		maxScoreStmt:                q.maxScoreStmt,
		problemStmt:                 q.problemStmt,
		problemTestsStmt:            q.problemTestsStmt,
		problemsStmt:                q.problemsStmt,
		proposersStmt:               q.proposersStmt,
		purgePbTestsStmt:            q.purgePbTestsStmt,
		setAdminStmt:                q.setAdminStmt,
		setBioStmt:                  q.setBioStmt,
		setCompilationStmt:          q.setCompilationStmt,
		setConsoleInputStmt:         q.setConsoleInputStmt,
		setEmailStmt:                q.setEmailStmt,
		setLimitsStmt:               q.setLimitsStmt,
		setMemoryLimitStmt:          q.setMemoryLimitStmt,
		setPbTestScoreStmt:          q.setPbTestScoreStmt,
		setPbTestVisibleIDStmt:      q.setPbTestVisibleIDStmt,
		setProblemDescriptionStmt:   q.setProblemDescriptionStmt,
		setProblemNameStmt:          q.setProblemNameStmt,
		setProblemVisibilityStmt:    q.setProblemVisibilityStmt,
		setProposerStmt:             q.setProposerStmt,
		setStackLimitStmt:           q.setStackLimitStmt,
		setSubmissionStatusStmt:     q.setSubmissionStatusStmt,
		setSubmissionTestStmt:       q.setSubmissionTestStmt,
		setSubmissionVisibilityStmt: q.setSubmissionVisibilityStmt,
		setTestNameStmt:             q.setTestNameStmt,
		setTimeLimitStmt:            q.setTimeLimitStmt,
		setVisibleIDStmt:            q.setVisibleIDStmt,
		subTestsStmt:                q.subTestsStmt,
		submissionStmt:              q.submissionStmt,
		submissionsStmt:             q.submissionsStmt,
		testStmt:                    q.testStmt,
		testVisibleIDStmt:           q.testVisibleIDStmt,
		userStmt:                    q.userStmt,
		userByEmailStmt:             q.userByEmailStmt,
		userByNameStmt:              q.userByNameStmt,
		userProblemSubmissionsStmt:  q.userProblemSubmissionsStmt,
		usersStmt:                   q.usersStmt,
		visibleProblemsStmt:         q.visibleProblemsStmt,
		waitingSubmissionsStmt:      q.waitingSubmissionsStmt,
	}
}
