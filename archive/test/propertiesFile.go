package test

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

type PropertiesRaw struct {
	Groups       string  `props:"groups"`
	Weights      string  `props:"weights"`
	Dependencies string  `props:"dependencies"`
	Time         float64 `props:"time"`
	Memory       float64 `props:"memory"`
}

func ParsePropertiesFile(r io.Reader) (*PropertiesRaw, bool, error) {
	vals := map[string][]string{}
	buf := bufio.NewScanner(r)
	for buf.Scan() {
		line := buf.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			return nil, false, nil
		}
		vals[kv[0]] = kv[1:]
	}
	if buf.Err() != nil {
		return nil, false, buf.Err()
	}

	dec := schema.NewDecoder()
	dec.SetAliasTag("props")

	rawProps := PropertiesRaw{}
	if err := dec.Decode(&rawProps, vals); err != nil {
		return nil, false, err
	}

	spew.Dump(rawProps)

	return &rawProps, true, nil
}

func ProcessPropertiesFile(ctx *ArchiveCtx, file *zip.File) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	rawProps, ok, err := ParsePropertiesFile(f)
	if err != nil {
		return err
	}
	if ok == false {
		return errors.New("Invalid properties file")
	}

	props := &Properties{
		TimeLimit:   rawProps.Time,
		MemoryLimit: int(rawProps.Memory * 1024.0),
	}

	// handle subtasks
	if rawProps.Groups != "" {
		subtaskedTests := map[int]bool{}

		type group struct{ start, end int }
		groups := map[int]group{}
		subTaskGroups := map[int][]group{}

		groupStrings := strings.Split(rawProps.Groups, ",")
		for i, g := range groupStrings {
			start, end := -1, -1
			if _, err := fmt.Sscanf(g, "%d-%d", &start, &end); err != nil {
				zap.S().Info(err)
				return errors.New("Invalid `group` string in properties")
			}
			groups[i+1] = group{start, end}
		}

		weights := map[int]int{}
		weightStrings := strings.Split(rawProps.Weights, ",")
		if len(groupStrings) != len(weightStrings) {
			return errors.New("Number of weights must match number of groups")
		}
		for i, w := range weightStrings {
			val, err := strconv.Atoi(w)
			if err != nil {
				return errors.New("Invalid `weight` string in properties")
			}
			weights[i+1] = val
		}

		if rawProps.Dependencies != "" {
			depStrings := strings.Split(rawProps.Dependencies, ",")
			if len(depStrings) != len(weightStrings) {
				return errors.New("Number of dependencies must match number of groups")
			}

			for i, d := range depStrings {
				subTaskGroups[i+1] = []group{groups[i+1]}
				if d == "" {
					continue
				}
				depGroups := strings.Split(d, ";")
				for _, dg := range depGroups {
					val, err := strconv.Atoi(dg)
					if err != nil {
						return fmt.Errorf("Invalid `dependencies` string in properties: %q is not a number", dg)
					}
					if val <= 0 || val > len(groupStrings) {
						return errors.New("Dependency number out of range")
					}
					subTaskGroups[i+1] = append(subTaskGroups[i+1], groups[val])
				}
			}
		} else {
			for i := range groupStrings {
				subTaskGroups[i+1] = []group{groups[i+1]}
			}
		}

		// coalesce maps into a single data type
		stks := map[int]Subtask{}

		for id, groups := range subTaskGroups {
			stk := Subtask{}
			stk.Score = weights[id]

			for _, group := range groups {
				for i := group.start; i <= group.end; i++ {
					subtaskedTests[i] = true
					stk.Tests = append(stk.Tests, i)
				}
			}
			sort.Ints(stk.Tests)

			stks[id] = stk
		}

		tests := []int{}
		for k := range subtaskedTests {
			tests = append(tests, k)
		}

		props.SubtaskedTests = tests
		props.Subtasks = stks
	}

	ctx.props = props
	return nil
}