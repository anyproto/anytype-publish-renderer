package renderer

import (
	"strings"
	"time"

	"github.com/gogo/protobuf/types"

	"go.uber.org/zap"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/localstore/addr"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

const (
	yesterday = "Yesterday"
	today     = "Today"
	tomorrow  = "Tomorrow"
)

func (r *Renderer) getDateSnapshot(objectId string) *pb.SnapshotWithType {
	date := strings.TrimPrefix(objectId, addr.DatePrefix)
	dateObject, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Error("failed to parse date", zap.Error(err))
		return nil
	}
	spaceId := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeySpaceId.String()]
	name := getDateName(dateObject)
	return &pb.SnapshotWithType{
		SbType: model.SmartBlockType_Date,
		Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyName.String():           pbtypes.String(name),
			bundle.RelationKeyId.String():             pbtypes.String(objectId),
			bundle.RelationKeyResolvedLayout.String(): pbtypes.Int64(int64(model.ObjectType_date)),
			bundle.RelationKeySpaceId.String():        spaceId,
		}}}},
	}
}

func getDateName(input time.Time) string {
	now := time.Now()
	currLocation := now.Location()
	todayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, currLocation)

	yesterdayTime := todayTime.AddDate(0, 0, -1)
	tomorrowTime := todayTime.AddDate(0, 0, 1)

	inputDate := time.Date(input.Year(), input.Month(), input.Day(), 0, 0, 0, 0, currLocation)

	switch inputDate {
	case todayTime:
		return today
	case yesterdayTime:
		return yesterday
	case tomorrowTime:
		return tomorrow
	default:
		return input.Format("January 2, 2006")
	}
}
