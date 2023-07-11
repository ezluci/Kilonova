package test

import (
	"archive/zip"
	"context"
	"io/fs"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/KiloProjects/kilonova"
	"github.com/KiloProjects/kilonova/eval"
	"github.com/KiloProjects/kilonova/sudoapi"
	"go.uber.org/zap"
)

var (
	ErrBadTestFile = kilonova.Statusf(400, "Bad test score file")
	ErrBadArchive  = kilonova.Statusf(400, "Bad archive")
)

type ArchiveCtx struct {
	tests       map[int]archiveTest
	attachments map[string]archiveAttachment
	scoredTests []int
	props       *properties

	submissions []*submissionStub

	params *TestProcessParams
}

type properties struct {
	Subtasks map[int]Subtask
	// seconds
	TimeLimit *float64
	// kbytes
	MemoryLimit *int

	Tags         []*mockTag
	Source       *string
	ConsoleInput *bool
	TestName     *string
	ProblemName  *string

	DefaultPoints *int

	SubtaskedTests []int

	Editors []string

	ScoringStrategy kilonova.ScoringType
}

func NewArchiveCtx(params *TestProcessParams) *ArchiveCtx {
	return &ArchiveCtx{
		tests:       make(map[int]archiveTest),
		attachments: make(map[string]archiveAttachment),
		scoredTests: make([]int, 0, 10),

		params: params,
	}
}

var (
	testInputSuffixes  = []string{".in", ".input"}
	testOutputSuffixes = []string{".out", ".output", ".ok", ".sol"}
)

func ProcessArchiveFile(ctx *ArchiveCtx, file *zip.File) *kilonova.StatusError {
	ext := strings.ToLower(path.Ext(file.Name))
	if slices.Contains(filepath.SplitList(path.Dir(file.Name)), "attachments") { // Is in "attachments" directory
		return ProcessAttachmentFile(ctx, file)
	}

	if slices.Contains(filepath.SplitList(path.Dir(file.Name)), "submissions") { // Is in "submissions" directory
		return ProcessSubmissionFile(ctx, file)
	}

	if ext == ".txt" { // test score file
		return ProcessScoreFile(ctx, file)
	}

	if ext == ".properties" { // test properties file
		return ProcessPropertiesFile(ctx, file)
	}

	if strings.ToLower(file.Name) == "problem.xml" { // Polygon archive format
		r, err := file.Open()
		if err != nil {
			return kilonova.WrapError(err, "Could not open problem.xml")
		}
		defer r.Close()
		return ProcessProblemXMLFile(ctx, r)
	}

	// Polygon-specific handling
	if ctx.params.Polygon {
		if strings.HasPrefix(file.Name, "solutions") {
			return ProcessSubmissionFile(ctx, file)
		}

		if strings.HasPrefix(file.Name, "tests") {
			if ext == ".a" {
				return ProcessTestOutputFile(ctx, file)
			}

			return ProcessTestInputFile(ctx, file)
		}

		if file.Name == "check.cpp" {
			return ProcessPolygonCheckFile(ctx, file)
		}

		return nil
	}

	// if nothing else is detected, it should be a test file
	if slices.Contains(testInputSuffixes, ext) { // test input file (ex: 01.in)
		return ProcessTestInputFile(ctx, file)
	}

	if slices.Contains(testOutputSuffixes, ext) { // test output file (ex: 01.out/01.ok)
		return ProcessTestOutputFile(ctx, file)
	}

	return nil
}

type TestProcessParams struct {
	Requestor *kilonova.UserBrief

	Polygon bool
}

func ProcessZipTestArchive(ctx context.Context, pb *kilonova.Problem, ar *zip.Reader, base *sudoapi.BaseAPI, params *TestProcessParams) *kilonova.StatusError {
	aCtx := NewArchiveCtx(params)

	// Try to autodetect polygon archive
	if _, err := fs.Stat(ar, "problem.xml"); err == nil {
		aCtx.params.Polygon = true
	}

	if params.Requestor == nil {
		return kilonova.Statusf(400, "There must be a requestor")
	}

	for _, file := range ar.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if err := ProcessArchiveFile(aCtx, file); err != nil {
			return err
		}
	}

	if aCtx.props != nil && aCtx.props.Subtasks != nil && len(aCtx.props.SubtaskedTests) != len(aCtx.tests) {
		zap.S().Info(len(aCtx.props.SubtaskedTests), len(aCtx.tests))
		return kilonova.Statusf(400, "Mismatched number of tests in archive and tests that correspond to at least one subtask")
	}

	for k, v := range aCtx.tests {
		if v.InFile == nil || v.OutFile == nil {
			return kilonova.Statusf(400, "Missing input or output file for test %d", k)
		}
	}

	if len(aCtx.scoredTests) != len(aCtx.tests) {
		// Try to deduce scoring remaining tests
		// zap.S().Info("Automatically inserting scores...")
		totalScore := 100
		for _, test := range aCtx.scoredTests {
			totalScore -= aCtx.tests[test].Score
		}

		// Since map order is ambiguous, get an ordered list of test IDs.
		// Regrettably, there is not easy way to do the set difference of the keys of the map and the scoredTests
		// so we'll do an O(N^2) operation for clarity's sake.
		testIDs := []int{}
		for id := range aCtx.tests {
			ok := true
			for _, scID := range aCtx.scoredTests {
				if id == scID {
					ok = false
					break
				}
			}
			if ok {
				testIDs = append(testIDs, id)
			}
		}
		sort.Ints(testIDs)

		n := len(aCtx.tests) - len(aCtx.scoredTests)
		perTest := totalScore/n + 1
		toSub := n - totalScore%n
		k := 0
		for _, i := range testIDs {
			if aCtx.tests[i].Score > 0 {
				continue
			}
			tst := aCtx.tests[i]
			tst.Score = perTest
			if k < toSub {
				tst.Score--
			}
			aCtx.tests[i] = tst
			k++
		}
	}

	// If we are loading an archive, the user might want to remove all tests first
	// So let's do it for them
	if err := base.DeleteTests(ctx, pb.ID); err != nil {
		zap.S().Warn(err)
		return err
	}

	createdTests := map[int]kilonova.Test{}

	for testID, v := range aCtx.tests {
		var test kilonova.Test
		test.ProblemID = pb.ID
		test.VisibleID = testID
		test.Score = v.Score
		if err := base.CreateTest(ctx, &test); err != nil {
			zap.S().Warn(err)
			return err
		}

		createdTests[testID] = test

		f, err := v.InFile.Open()
		if err != nil {
			return kilonova.WrapError(err, "Couldn't open() input file")
		}
		if err := base.SaveTestInput(test.ID, f); err != nil {
			zap.S().Warn("Couldn't create test input", err)
			f.Close()
			return kilonova.WrapError(err, "Couldn't create test input")
		}
		f.Close()
		f, err = v.OutFile.Open()
		if err != nil {
			return kilonova.WrapError(err, "Couldn't open() output file")
		}
		if err := base.SaveTestOutput(test.ID, f); err != nil {
			zap.S().Warn("Couldn't create test output", err)
			f.Close()
			return kilonova.WrapError(err, "Couldn't create test output")
		}
		f.Close()
	}

	if len(aCtx.attachments) > 0 {
		atts, err := base.ProblemAttachments(ctx, pb.ID)
		if err != nil {
			zap.S().Warn("Couldn't get problem attachments")
			return kilonova.WrapError(err, "Couldn't get attachments")
		}
		attIDs := []int{}
		for _, att := range atts {
			attIDs = append(attIDs, att.ID)
		}
		// TODO: In the future maybe opt in to a "merging" strategy instead of delete and add?
		if len(attIDs) > 0 {
			if _, err := base.DeleteProblemAtts(ctx, pb.ID, attIDs); err != nil {
				zap.S().Warn("Couldn't remove attachments")
				return kilonova.WrapError(err, "Couldn't delete attachments")
			}
		}
		for _, att := range aCtx.attachments {
			if att.File == nil {
				zap.S().Infof("Skipping attachment %s since it only has props", att.Name)
				continue
			}

			f, err := att.File.Open()
			if err != nil {
				zap.S().Warn("Couldn't open attachment zip file", err)
				continue
			}

			var userID *int
			if params.Requestor != nil {
				userID = &params.Requestor.ID
			}

			if err := base.CreateProblemAttachment(ctx, &kilonova.Attachment{
				Name:    att.Name,
				Private: att.Private,
				Visible: att.Visible,
				Exec:    att.Exec,
			}, pb.ID, f, userID); err != nil {
				zap.S().Warn("Couldn't create attachment", err)
				f.Close()
				continue
			}
			f.Close()
		}
	}

	if aCtx.props != nil {
		shouldUpd := false
		upd := kilonova.ProblemUpdate{}
		if aCtx.props.MemoryLimit != nil {
			shouldUpd = true
			upd.MemoryLimit = aCtx.props.MemoryLimit
		}
		if aCtx.props.TimeLimit != nil {
			shouldUpd = true
			upd.TimeLimit = aCtx.props.TimeLimit
		}
		if aCtx.props.DefaultPoints != nil {
			shouldUpd = true
			upd.DefaultPoints = aCtx.props.DefaultPoints
		}
		if aCtx.props.Source != nil {
			shouldUpd = true
			upd.SourceCredits = aCtx.props.Source
		}
		if aCtx.props.ConsoleInput != nil {
			shouldUpd = true
			upd.ConsoleInput = aCtx.props.ConsoleInput
		}
		if aCtx.props.ScoringStrategy != kilonova.ScoringTypeNone {
			shouldUpd = true
			upd.ScoringStrategy = aCtx.props.ScoringStrategy
		}
		if aCtx.props.TestName != nil {
			shouldUpd = true
			upd.TestName = aCtx.props.TestName
		}

		if aCtx.props.ProblemName != nil && *aCtx.props.ProblemName != "" {
			shouldUpd = true
			upd.Name = aCtx.props.ProblemName
		}

		if shouldUpd {
			if err := base.UpdateProblem(ctx, pb.ID, upd, nil); err != nil {
				zap.S().Warn(err)
				return kilonova.WrapError(err, "Couldn't update problem medatada")
			}
		}

		if aCtx.props.Subtasks != nil {
			if err := base.DeleteSubTasks(ctx, pb.ID); err != nil {
				zap.S().Warn(err)
				return kilonova.WrapError(err, "Couldn't delete existing subtasks")
			}
			for stkId, stk := range aCtx.props.Subtasks {
				outStk := kilonova.SubTask{
					ProblemID: pb.ID,
					VisibleID: stkId,
					Score:     stk.Score,
					Tests:     []int{},
				}
				for _, test := range stk.Tests {
					if tt, exists := createdTests[test]; !exists {
						return kilonova.Statusf(400, "Test %d not found in added tests. Aborting subtask creation", test)
					} else {
						outStk.Tests = append(outStk.Tests, tt.ID)
					}
				}

				if err := base.CreateSubTask(ctx, &outStk); err != nil {
					zap.S().Warn(err)
					return kilonova.WrapError(err, "Couldn't create subtask")
				}
			}
		}

		if len(aCtx.props.Tags) > 0 {
			realTagIDs := []int{}
			for _, mTag := range aCtx.props.Tags {
				tag, err := base.TagByLooseName(ctx, mTag.Name)
				if err != nil || tag == nil {
					id, err := base.CreateTag(ctx, mTag.Name, mTag.Type)
					if err != nil {
						zap.S().Warn("Couldn't create tag: ", err)
						continue
					}
					realTagIDs = append(realTagIDs, id)
					continue
				}
				realTagIDs = append(realTagIDs, tag.ID)
			}
			if err := base.UpdateProblemTags(ctx, pb.ID, realTagIDs); err != nil {
				zap.S().Warn(err)
				return kilonova.WrapError(err, "Couldn't update tags")
			}
		}

		if len(aCtx.props.Editors) > 0 && params.Requestor.Admin {
			var newEditors []*kilonova.UserBrief
			// First, get the new editors to make sure they are valid
			for _, editor := range aCtx.props.Editors {
				user, err := base.UserBriefByName(ctx, editor)
				if err == nil && user != nil {
					newEditors = append(newEditors, user)
				}
			}

			if len(newEditors) > 0 {
				// Then, remove existing editors
				cEditors, err := base.ProblemEditors(ctx, pb.ID)
				if err != nil {
					zap.S().Warn(err)
					return err
				}
				for _, ed := range cEditors {
					if err := base.StripProblemAccess(ctx, pb.ID, ed.ID); err != nil {
						zap.S().Warn(err)
					}
				}

				// Lastly, add the new editors
				for _, editor := range newEditors {
					if err := base.AddProblemEditor(ctx, pb.ID, editor.ID); err != nil {
						zap.S().Warn(err)
					}
				}
			}

		}

	}

	// Do submissions at the end after all changes have been merged
	if len(aCtx.submissions) > 0 {
		for _, sub := range aCtx.submissions {
			lang, ok := eval.Langs[sub.lang]
			if !ok {
				zap.S().Warn("Skipping submission")
				continue
			}
			if _, err := base.CreateSubmission(ctx, params.Requestor, pb, sub.code, lang, nil, true); err != nil {
				zap.S().Warn(err)
			}
		}
	}

	return nil
}
