package sudoapi

import (
	"context"
	"errors"

	"github.com/KiloProjects/kilonova"
	"go.uber.org/zap"
)

func (s *BaseAPI) Tags(ctx context.Context) ([]*kilonova.Tag, *StatusError) {
	tags, err := s.db.Tags(ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			zap.S().Warn(err)
		}
		return nil, WrapError(err, "Couldn't get tags")
	}
	return tags, nil
}

func (s *BaseAPI) TagsByType(ctx context.Context, tagType kilonova.TagType) ([]*kilonova.Tag, *StatusError) {
	tags, err := s.db.TagsByType(ctx, tagType)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			zap.S().Warn(err)
		}
		return nil, WrapError(err, "Couldn't get tags")
	}
	return tags, nil
}

func (s *BaseAPI) RelevantTags(ctx context.Context, tagID int, max int) ([]*kilonova.Tag, *StatusError) {
	tags, err := s.db.RelevantTags(ctx, tagID, max)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			zap.S().Warn(err)
		}
		return nil, WrapError(err, "Couldn't get relevant tags")
	}
	return tags, nil
}

func (s *BaseAPI) TagByID(ctx context.Context, id int) (*kilonova.Tag, *StatusError) {
	tag, err := s.db.Tag(ctx, id)
	if err != nil || tag == nil {
		return nil, WrapError(err, "Tag not found")
	}
	return tag, nil
}

func (s *BaseAPI) TagByName(ctx context.Context, name string) (*kilonova.Tag, *StatusError) {
	tag, err := s.db.TagByName(ctx, name)
	if err != nil || tag == nil {
		return nil, WrapError(err, "Tag not found")
	}
	return tag, nil
}

func (s *BaseAPI) UpdateTagName(ctx context.Context, id int, newName string) *StatusError {
	if err := s.db.UpdateTagName(ctx, id, newName); err != nil {
		return WrapError(err, "Couldn't update tag")
	}
	return nil
}

func (s *BaseAPI) UpdateTagType(ctx context.Context, id int, newType kilonova.TagType) *StatusError {
	if err := s.db.UpdateTagType(ctx, id, newType); err != nil {
		return WrapError(err, "Couldn't update tag")
	}
	return nil
}

func (s *BaseAPI) DeleteTag(ctx context.Context, id int) *StatusError {
	if err := s.db.DeleteTag(ctx, id); err != nil {
		return WrapError(err, "Couldn't delete tag")
	}
	return nil
}

func (s *BaseAPI) CreateTag(ctx context.Context, name string, tagType kilonova.TagType) (int, *StatusError) {
	if name == "" || tagType == kilonova.TagTypeNone {
		return -1, ErrMissingRequired
	}
	id, err := s.db.CreateTag(ctx, name, tagType)
	if err != nil {
		return -1, WrapError(err, "Couldn't create tag")
	}
	return id, nil
}

// master - the OG that will remain after the merge
// clone - the one that will be replaced
func (s *BaseAPI) MergeTags(ctx context.Context, master int, clone int) *StatusError {
	if err := s.db.MergeTags(ctx, master, clone); err != nil {
		return WrapError(err, "Couldn't merge tags")
	}
	return nil
}

func (s *BaseAPI) ProblemTags(ctx context.Context, problemID int) ([]*kilonova.Tag, *StatusError) {
	tags, err := s.db.ProblemTags(ctx, problemID)
	if err != nil {
		return nil, WrapError(err, "Couldn't get problem tags")
	}
	return tags, nil
}

func (s *BaseAPI) UpdateProblemTags(ctx context.Context, problemID int, tagIDs []int) *StatusError {
	if err := s.db.UpdateProblemTags(ctx, problemID, tagIDs); err != nil {
		return WrapError(err, "Couldn't update problem tags")
	}
	return nil
}
