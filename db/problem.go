package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KiloProjects/kilonova"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
)

type pbAccessType string

const (
	pbAccessEditor pbAccessType = "editor"
	pbAccessViewer pbAccessType = "viewer"
)

func (s *DB) Problem(ctx context.Context, id int) (*kilonova.Problem, error) {
	var pb dbProblem
	err := s.conn.GetContext(ctx, &pb, s.conn.Rebind("SELECT * FROM problems WHERE id = ? LIMIT 1"), id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s.internalToProblem(ctx, &pb)
}

func (s *DB) VisibleProblem(ctx context.Context, id int, user *kilonova.UserBrief) (*kilonova.Problem, error) {
	pbs, err := s.Problems(ctx, kilonova.ProblemFilter{ID: &id, LookingUser: user, Look: true})
	if err != nil || len(pbs) == 0 {
		return nil, nil
	}
	return pbs[0], nil
}

func (s *DB) Problems(ctx context.Context, filter kilonova.ProblemFilter) ([]*kilonova.Problem, error) {
	var pbs []*dbProblem
	where, args := problemFilterQuery(&filter)
	query := s.conn.Rebind("SELECT * FROM problems WHERE " + strings.Join(where, " AND ") + " ORDER BY id ASC " + FormatLimitOffset(filter.Limit, filter.Offset))
	err := s.conn.SelectContext(ctx, &pbs, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return []*kilonova.Problem{}, nil
	}
	return s.internalToProblems(ctx, pbs), err
}

const problemCreateQuery = `INSERT INTO problems (
	name, description, console_input, test_name, memory_limit, source_size, time_limit, visible, source_credits, author_credits, short_description, default_points
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id;`

func (s *DB) CreateProblem(ctx context.Context, p *kilonova.Problem, authorID int) error {
	if p.Name == "" || authorID == 0 {
		return kilonova.ErrMissingRequired
	}
	if p.TestName == "" {
		p.TestName = slug.Make(p.Name)
	}
	if p.MemoryLimit == 0 {
		p.MemoryLimit = 65536 // 64MB
	}
	if p.TimeLimit == 0 {
		p.TimeLimit = 1 // 1s
	}
	if p.SourceSize == 0 {
		p.SourceSize = 10000
	}
	var id int
	err := s.conn.GetContext(ctx, &id, s.conn.Rebind(problemCreateQuery), p.Name, p.Description, p.ConsoleInput, p.TestName, p.MemoryLimit, p.SourceSize, p.TimeLimit, p.Visible, p.SourceCredits, p.AuthorCredits, p.ShortDesc, p.DefaultPoints)
	if err == nil {
		p.ID = id
	}
	if err != nil {
		return err
	}
	return s.AddProblemEditor(ctx, id, authorID)
}

const problemUpdateStatement = `UPDATE problems SET %s WHERE id = ?`

func (s *DB) UpdateProblem(ctx context.Context, id int, upd kilonova.ProblemUpdate) error {
	toUpd, args := problemUpdateQuery(&upd)
	if len(toUpd) == 0 {
		return kilonova.ErrNoUpdates
	}
	args = append(args, id)
	query := s.conn.Rebind(fmt.Sprintf(problemUpdateStatement, strings.Join(toUpd, ", ")))
	_, err := s.conn.ExecContext(ctx, query, args...)
	return err
}

const bulkProblemUpdateStatement = `UPDATE problems SET %s WHERE %s`

func (s *DB) BulkUpdateProblems(ctx context.Context, filter kilonova.ProblemFilter, upd kilonova.ProblemUpdate) error {
	toUpd, args := problemUpdateQuery(&upd)
	if len(toUpd) == 0 {
		return kilonova.ErrNoUpdates
	}
	where, args1 := problemFilterQuery(&filter)
	args = append(args, args1...)
	query := s.conn.Rebind(fmt.Sprintf(bulkProblemUpdateStatement, strings.Join(toUpd, ", "), strings.Join(where, ", ")))
	_, err := s.conn.ExecContext(ctx, query, args...)
	return err
}

func (s *DB) DeleteProblem(ctx context.Context, id int) error {
	_, err := s.conn.ExecContext(ctx, s.conn.Rebind("DELETE FROM problems WHERE id = ?"), id)
	return err
}

func problemFilterQuery(filter *kilonova.ProblemFilter) ([]string, []interface{}) {
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, v)
	}
	if v := filter.IDs; v != nil && len(v) == 0 {
		where = append(where, "id = -1")
	}
	if v := filter.IDs; len(v) > 0 {
		where = append(where, "id IN (?"+strings.Repeat(",?", len(v)-1)+")")
		for _, el := range v {
			args = append(args, el)
		}
	}
	if v := filter.Name; v != nil {
		where, args = append(where, "lower(name) = lower(?)"), append(args, v)
	}
	if v := filter.ConsoleInput; v != nil {
		where, args = append(where, "console_input = ?"), append(args, v)
	}
	if v := filter.Visible; v != nil {
		where, args = append(where, "visible = ?"), append(args, v)
	}
	if filter.Look {
		var id int
		if filter.LookingUser != nil {
			id = filter.LookingUser.ID
			if filter.LookingUser.Admin {
				id = -1
			}
		}
		if id >= 0 {
			where, args = append(where, "(visible = true OR id IN (SELECT DISTINCT problem_id FROM problem_user_access WHERE user_id = ?))"), append(args, id, id)
		}
	}
	if filter.Unassociated {
		where = append(where, "id NOT IN (SELECT DISTINCT problem_id FROM problem_list_problems)")
	}
	return where, args
}

func problemUpdateQuery(upd *kilonova.ProblemUpdate) ([]string, []interface{}) {
	toUpd, args := []string{}, []interface{}{}
	if v := upd.Name; v != nil {
		toUpd, args = append(toUpd, "name = ?"), append(args, v)
	}
	if v := upd.Description; v != nil {
		toUpd, args = append(toUpd, "description = ?"), append(args, v)
	}
	if v := upd.ShortDesc; v != nil {
		toUpd, args = append(toUpd, "short_description = ?"), append(args, v)
	}

	if v := upd.TestName; v != nil {
		toUpd, args = append(toUpd, "test_name = ?"), append(args, v)
	}

	if v := upd.TimeLimit; v != nil {
		toUpd, args = append(toUpd, "time_limit = ?"), append(args, v)
	}
	if v := upd.MemoryLimit; v != nil {
		toUpd, args = append(toUpd, "memory_limit = ?"), append(args, v)
	}

	if v := upd.DefaultPoints; v != nil {
		toUpd, args = append(toUpd, "default_points = ?"), append(args, v)
	}

	if v := upd.SourceCredits; v != nil {
		toUpd, args = append(toUpd, "source_credits = ?"), append(args, v)
	}
	if v := upd.AuthorCredits; v != nil {
		toUpd, args = append(toUpd, "author_credits = ?"), append(args, v)
	}

	if v := upd.ConsoleInput; v != nil {
		toUpd, args = append(toUpd, "console_input = ?"), append(args, v)
	}
	if v := upd.Visible; v != nil {
		toUpd, args = append(toUpd, "visible = ?"), append(args, v)
	}

	return toUpd, args
}

// TODO: Move to BaseAPI
func (db *DB) SolvedProblems(ctx context.Context, uid int) ([]*kilonova.Problem, error) {
	ids, err := db.solvedProblems(ctx, uid)
	if err != nil {
		return nil, err
	}
	var pbs = make([]*kilonova.Problem, 0, len(ids))
	for _, id := range ids {
		pb, err := db.Problem(ctx, id)
		if err != nil {
			log.Printf("Couldn't get solved problem %d: %s\n", id, err)
		} else {
			pbs = append(pbs, pb)
		}
	}
	return pbs, nil
}

type problemAccess struct {
	ProblemID int          `db:"problem_id"`
	UserID    int          `db:"user_id"`
	Access    pbAccessType `db:"access"`
}

func (s *DB) problemAccessRights(ctx context.Context, pbid int) ([]*problemAccess, error) {
	var rights []*problemAccess
	err := s.conn.SelectContext(ctx, &rights, "SELECT * FROM problem_user_access WHERE problem_id = $1", pbid)
	if errors.Is(err, sql.ErrNoRows) {
		return []*problemAccess{}, nil
	}
	if err != nil {
		return nil, err
	}
	return rights, nil
}

func (s *DB) AddProblemEditor(ctx context.Context, pbid int, uid int) error {
	_, err := s.conn.ExecContext(ctx, "INSERT INTO problem_user_access (problem_id, user_id, access) VALUES ($1, $2, 'editor')", pbid, uid)
	return err
}

func (s *DB) AddProblemViewer(ctx context.Context, pbid int, uid int) error {
	_, err := s.conn.ExecContext(ctx, "INSERT INTO problem_user_access (problem_id, user_id, access) VALUES ($1, $2, 'viewer')", pbid, uid)
	return err
}

func (s *DB) StripProblemAccess(ctx context.Context, pbid int, uid int) error {
	_, err := s.conn.ExecContext(ctx, "DELETE FROM problem_user_access WHERE problem_id = $1 AND user_id = $2", pbid, uid)
	return err
}

type dbProblem struct {
	ID            int       `db:"id"`
	CreatedAt     time.Time `db:"created_at"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	ShortDesc     string    `db:"short_description"`
	TestName      string    `db:"test_name"`
	Visible       bool      `db:"visible"`
	DefaultPoints int       `db:"default_points"`

	// Limit stuff
	TimeLimit   float64 `db:"time_limit"`
	MemoryLimit int     `db:"memory_limit"`
	SourceSize  int     `db:"source_size"`

	SourceCredits string `db:"source_credits"`
	AuthorCredits string `db:"author_credits"`

	// Eval stuff
	ConsoleInput bool `db:"console_input"`
}

func (s *DB) internalToProblems(ctx context.Context, pbs []*dbProblem) []*kilonova.Problem {
	return mapper(pbs, func(pb *dbProblem) *kilonova.Problem {
		pbb, err := s.internalToProblem(ctx, pb)
		if err != nil {
			zap.S().Warn(err)
		}
		return pbb
	})
}

func (s *DB) internalToProblem(ctx context.Context, pb *dbProblem) (*kilonova.Problem, error) {
	if pb == nil {
		return nil, nil
	}

	rights, err := s.problemAccessRights(ctx, pb.ID)
	if err != nil {
		return nil, err
	}

	var editors, viewers []int
	for _, right := range rights {
		switch right.Access {
		case pbAccessEditor:
			editors = append(editors, right.UserID)
		case pbAccessViewer:
			viewers = append(viewers, right.UserID)
		default:
			zap.S().Warn("Unknown access rank", zap.String("right", string(right.Access)))
		}
	}
	if editors == nil {
		editors = []int{}
	}
	if viewers == nil {
		viewers = []int{}
	}

	return &kilonova.Problem{
		ID:          pb.ID,
		CreatedAt:   pb.CreatedAt,
		Name:        pb.Name,
		Description: pb.Description,
		ShortDesc:   pb.ShortDesc,
		TestName:    pb.TestName,

		Visible: pb.Visible,
		Editors: editors,
		Viewers: viewers,

		DefaultPoints: pb.DefaultPoints,

		TimeLimit:   pb.TimeLimit,
		MemoryLimit: pb.MemoryLimit,
		SourceSize:  pb.SourceSize,

		SourceCredits: pb.SourceCredits,
		AuthorCredits: pb.AuthorCredits,

		ConsoleInput: pb.ConsoleInput,
	}, nil
}
