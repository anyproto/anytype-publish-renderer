package renderer

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							"shortTextKey": pbtypes.String("test"),
						},
					},
				}},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "shortTextKey.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "shortTextKey"}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("long text", func(t *testing.T) {
		// given
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							"longTextKey": pbtypes.String("test"),
						},
					},
				}},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "longTextKey.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "longTextKey"}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("number", func(t *testing.T) {
		// given
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							"numberKey": pbtypes.Int64(4),
						},
					},
				},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "numberKey.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "numberKey"}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("phone", func(t *testing.T) {
		// given
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							"phone": pbtypes.Int64(12345),
						},
					}}},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "phone.pb"),
				&pb.SnapshotWithType{
					SbType: model.SmartBlockType_STRelation,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{
							Fields: map[string]*types.Value{
								bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("phone").URL()),
								bundle.RelationKeyName.String():           pbtypes.String("Phone Relation"),
								bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_phone)),
							}},
					},
					},
				},
			),
		)
		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "phone"}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("email", func(t *testing.T) {
		// given
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							"email": pbtypes.String("email"),
						},
					},
				}},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "email.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "email"}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, component)
	})
	t.Run("empty key", func(t *testing.T) {
		// given
		r := NewTestRenderer()
		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: ""}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.Nil(t, component)
	})
	t.Run("unknown key", func(t *testing.T) {
		// given
		r := NewTestRenderer()
		block := &model.Block{Content: &model.BlockContentOfRelation{Relation: &model.BlockContentRelation{Key: "unknown key"}}}

		// when
		component := r.makeRelationTemplate(block)

		// then
		assert.Nil(t, component)
	})
	t.Run("date value", func(t *testing.T) {
		// given
		relationKey := "date-relation"
		timestamp := float64(time.Now().Unix())

		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							relationKey: pbtypes.Float64(timestamp),
						},
					},
				}},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "date-relation.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
	})
	t.Run("Block with checkbox relation value", func(t *testing.T) {
		// given
		relationKey := "checkbox-relation"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							relationKey: pbtypes.Bool(true),
						},
					},
				}},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "checkbox-relation.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "checkbox-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
	})
	t.Run("Block with c-object relation value", func(t *testing.T) {
		// given
		relationKey := "object-relation"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{relationKey: pbtypes.StringList([]string{"object1", "object2"})},
					}},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "object-relation.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "object-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
	})

	t.Run("Block with c-select relation value (status)", func(t *testing.T) {
		// given
		relationKey := "status-relation"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							relationKey: pbtypes.String("status1"),
						},
					},
				},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "status-relation.pb"), &pb.SnapshotWithType{
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
			}),
			WithLinkedSnapshot(t, filepath.Join("relationsOptions", "status1.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "status-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
	})

	t.Run("Block with c-select relation value (tag)", func(t *testing.T) {
		// given
		relationKey := "tag-relation"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							relationKey: pbtypes.StringList([]string{"tag1", "tag2"}),
						},
					},
				},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "tag-relation.pb"), &pb.SnapshotWithType{
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
			}),
			WithLinkedSnapshot(t, filepath.Join("relationsOptions", "tag1.pb"), &pb.SnapshotWithType{
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
			}),
			WithLinkedSnapshot(t, filepath.Join("relationsOptions", "tag2.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "tag-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
	})
	t.Run("empty list relation", func(t *testing.T) {
		// given
		relationKey := "tag-relation"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							relationKey: pbtypes.StringList([]string{"tag1", "tag2"}),
						},
					},
				},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "tag-relation.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "tag-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
		builder := strings.Builder{}
		err := params.Render(context.Background(), &builder)
		assert.NoError(t, err)
		assert.Equal(t, `<div class="sides"><div class="info"><div class="name">Tag Relation</div></div></div>`, builder.String())
	})
	t.Run("empty relation", func(t *testing.T) {
		// given
		relationKey := "text-relation"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							relationKey: pbtypes.String(""),
						},
					},
				},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("relations", "text-relation.pb"), &pb.SnapshotWithType{
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
			}),
		)

		block := &model.Block{
			Id: "text-value",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		// when
		params := r.makeRelationTemplate(block)

		// then
		assert.NotNil(t, params)
		builder := strings.Builder{}
		err := params.Render(context.Background(), &builder)
		assert.NoError(t, err)
		assert.Equal(t, `<div class="sides"><div class="info"><div class="name">Text Relation</div></div></div>`, builder.String())
	})
}
