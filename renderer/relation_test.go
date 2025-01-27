package renderer

import (
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/jsonpb"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestMakeRelationRenderParams(t *testing.T) {
	renderer := getTestRenderer("Anytype.WebPublish.20241217.112212.67")

	t.Run("Default case with no relation and no color", func(t *testing.T) {
		// given
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: "nonexistent-key",
				},
			},
		}

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "test-block-id", params.Id)
		assert.Equal(t, "", params.BackgroundColor)
		assert.Equal(t, "isDeleted", params.IsDeleted)
		assert.Equal(t, "isEmpty", params.IsEmpty)
		assert.Equal(t, defaultName, params.Name)
		assert.Equal(t, "", params.Format)
		assert.Nil(t, params.Value)
	})

	t.Run("Block with background color", func(t *testing.T) {
		block := &model.Block{
			Id:              "test-block-id",
			BackgroundColor: "red",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: "nonexistent-key",
				},
			},
		}

		params := renderer.MakeRelationRenderParams(block)

		assert.Equal(t, "bgColor bgColor-red", params.BackgroundColor)
	})

	t.Run("Block with existing relation metadata", func(t *testing.T) {
		// given
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: bundle.RelationKeyName.String(),
				},
			},
		}

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Name", params.Name)
		assert.Equal(t, "", params.IsDeleted)
		assert.Equal(t, "c-shortText", params.Format)
	})

	t.Run("Block with relation value", func(t *testing.T) {
		// given
		relationKey := "relation-with-value"
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.Int64(1),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("relation-with-value").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Value Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_number)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "relation-with-value.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Value Relation", params.Name)
		assert.Equal(t, "c-number", params.Format)
		assert.NotNil(t, params.Value)
	})

	t.Run("Block with date relation value", func(t *testing.T) {
		// given
		relationKey := "date-relation"

		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		timestamp := float64(time.Now().Unix())
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.Float64(timestamp),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("date-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Date Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_date)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "date-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Date Relation", params.Name)
		assert.Equal(t, "c-date", params.Format)
		assert.NotNil(t, params.Value)
	})
	t.Run("Block with checkbox relation value", func(t *testing.T) {
		// given
		relationKey := "checkbox-relation"

		block := &model.Block{
			Id: "checkbox-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.Bool(true),
					},
				},
			}},
		}

		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("checkbox-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Checkbox Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_checkbox)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "checkbox-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Checkbox Relation", params.Name)
		assert.Equal(t, "c-checkbox", params.Format)
		assert.NotNil(t, params.Value)
	})
	t.Run("Block with c-object relation value", func(t *testing.T) {
		// given
		relationKey := "object-relation"

		block := &model.Block{
			Id: "object-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{relationKey: pbtypes.StringList([]string{"object1", "object2"})},
				}},
			},
		}

		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("object-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Object Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_object)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "object-relation.pb")] = json

		sn = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("object1"),
						bundle.RelationKeyName.String():    pbtypes.String("Object1"),
						bundle.RelationKeyLayout.String():  pbtypes.Int64(int64(model.ObjectType_todo)),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					},
				},
			}},
		}
		json, err = marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("objects", "object1.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Object Relation", params.Name)
		assert.Equal(t, "c-object", params.Format)
		assert.NotNil(t, params.Value)
	})

	t.Run("Block with c-select relation value (status)", func(t *testing.T) {
		// given
		relationKey := "status-relation"

		block := &model.Block{
			Id: "status-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.String("status1"),
					},
				},
			},
			},
		}

		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("status-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Status Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_status)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "status-relation.pb")] = json

		sn = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("status1"),
						bundle.RelationKeyName.String():    pbtypes.String("Status1"),
						bundle.RelationKeyLayout.String():  pbtypes.Int64(int64(model.ObjectType_relationOption)),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					},
				},
			}},
		}
		json, err = marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relationsOptions", "status1.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Status Relation", params.Name)
		assert.Equal(t, "c-select", params.Format)
		assert.NotNil(t, params.Value)
	})

	t.Run("Block with c-select relation value (tag)", func(t *testing.T) {
		// given
		relationKey := "tag-relation"

		block := &model.Block{
			Id: "tag-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.StringList([]string{"tag1", "tag2"}),
					},
				},
			},
			},
		}

		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("tag-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Tag Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_tag)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "tag-relation.pb")] = json

		sn = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("tag1"),
						bundle.RelationKeyName.String():    pbtypes.String("Tag1"),
						bundle.RelationKeyLayout.String():  pbtypes.Int64(int64(model.ObjectType_relationOption)),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					},
				},
			}},
		}
		json, err = marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relationsOptions", "tag1.pb")] = json

		sn = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("tag2"),
						bundle.RelationKeyName.String():    pbtypes.String("Tag2"),
						bundle.RelationKeyLayout.String():  pbtypes.Int64(int64(model.ObjectType_relationOption)),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					},
				},
			}},
		}
		json, err = marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relationsOptions", "tag2.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Tag Relation", params.Name)
		assert.Equal(t, "c-select", params.Format)
		assert.NotNil(t, params.Value)
	})
	t.Run("Block with long text value", func(t *testing.T) {
		// given
		relationKey := "text-relation"
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.String("test"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String(): pbtypes.String(domain.RelationKey("text-relation").URL()),
						bundle.RelationKeyName.String():      pbtypes.String("Text Relation"),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "text-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Text Relation", params.Name)
		assert.Equal(t, "c-longText", params.Format)
		assert.NotNil(t, params.Value)
	})
	t.Run("Block with url text value", func(t *testing.T) {
		// given
		relationKey := "url-relation"
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.String("url"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("url-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Url Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_url)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "url-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Url Relation", params.Name)
		assert.Equal(t, "c-url", params.Format)
		assert.NotNil(t, params.Value)
	})
	t.Run("Block with email text value", func(t *testing.T) {
		// given
		relationKey := "email-relation"
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.String("email"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("email-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Email Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_email)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "email-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Email Relation", params.Name)
		assert.Equal(t, "c-email", params.Format)
		assert.NotNil(t, params.Value)
	})
	t.Run("Block with phone text value", func(t *testing.T) {
		// given
		relationKey := "phone-relation"
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.String("phone"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("phone-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Phone Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_phone)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "phone-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Phone Relation", params.Name)
		assert.Equal(t, "c-phone", params.Format)
		assert.NotNil(t, params.Value)
	})
}
