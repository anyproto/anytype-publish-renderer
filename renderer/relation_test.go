package renderer

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func TestMakeRelationRenderParams_ShortText(t *testing.T) {
	t.Run("short text", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						"shortTextKey": pbtypes.String("test"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("shortTextKey").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Short text Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_shorttext)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "shortTextKey.pb")] = json

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "shortTextKey"}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("long text", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						"longTextKey": pbtypes.String("test"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("longTextKey").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Long text Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_longtext)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "longTextKey.pb")] = json

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "longTextKey"}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("number", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						"numberKey": pbtypes.Int64(4),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("numberKey").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Number Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_number)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "numberKey.pb")] = json

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "numberKey"}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("phone", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						"phone": pbtypes.Int64(12345),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("phone").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Phone Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_phone)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "phone.pb")] = json

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "phone"}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("email", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						"email": pbtypes.String("email"),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("email").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Email Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_email)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "email.pb")] = json

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "email"}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("empty key", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: ""}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.Nil(t, component)
	})
	t.Run("unknown key", func(t *testing.T) {
		// given
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "unknown key"}}}

		// when
		component := r.MakeRelationRenderParams(block)

		// then
		assert.Nil(t, component)
	})
	t.Run("date value", func(t *testing.T) {
		// given
		relationKey := "date-relation"
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		timestamp := float64(time.Now().Unix())
		r.Sp = &pb.SnapshotWithType{
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
		r.UberSp.PbFiles[filepath.Join("relations", "date-relation.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
	})
	t.Run("Block with checkbox relation value", func(t *testing.T) {
		// given
		relationKey := "checkbox-relation"
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "checkbox-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		r.Sp = &pb.SnapshotWithType{
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
		r.UberSp.PbFiles[filepath.Join("relations", "checkbox-relation.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
	})
	t.Run("Block with c-object relation value", func(t *testing.T) {
		// given
		relationKey := "object-relation"
		r := &Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "object-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		r.Sp = &pb.SnapshotWithType{
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
		r.UberSp.PbFiles[filepath.Join("relations", "object-relation.pb")] = json

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
		r.UberSp.PbFiles[filepath.Join("objects", "object1.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
	})

	t.Run("Block with c-select relation value (status)", func(t *testing.T) {
		// given
		relationKey := "status-relation"
		r := &Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "status-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		r.Sp = &pb.SnapshotWithType{
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
		r.UberSp.PbFiles[filepath.Join("relations", "status-relation.pb")] = json

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
		r.UberSp.PbFiles[filepath.Join("relationsOptions", "status1.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
	})

	t.Run("Block with c-select relation value (tag)", func(t *testing.T) {
		// given
		relationKey := "tag-relation"
		r := &Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "tag-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		r.Sp = &pb.SnapshotWithType{
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
		r.UberSp.PbFiles[filepath.Join("relations", "tag-relation.pb")] = json

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
		r.UberSp.PbFiles[filepath.Join("relationsOptions", "tag1.pb")] = json

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
		r.UberSp.PbFiles[filepath.Join("relationsOptions", "tag2.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
	})
	t.Run("empty list relation", func(t *testing.T) {
		// given
		relationKey := "tag-relation"
		r := &Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "tag-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		r.Sp = &pb.SnapshotWithType{
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
		r.UberSp.PbFiles[filepath.Join("relations", "tag-relation.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
		builder := strings.Builder{}
		err = params.Render(context.Background(), &builder)
		assert.NoError(t, err)
		assert.Equal(t, `<div class="sides"><div class="info"><div class="name">Tag Relation</div></div></div>`, builder.String())
	})
	t.Run("empty relation", func(t *testing.T) {
		// given
		relationKey := "text-relation"
		r := &Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		block := &model.Block{
			Id: "text-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.String(""),
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
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("text-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Text Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_shorttext)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "text-relation.pb")] = json

		// when
		params := r.MakeRelationRenderParams(block)

		// then
		assert.NotNil(t, params)
		builder := strings.Builder{}
		err = params.Render(context.Background(), &builder)
		assert.NoError(t, err)
		assert.Equal(t, `<div class="sides"><div class="info"><div class="name">Text Relation</div></div></div>`, builder.String())
	})
}
