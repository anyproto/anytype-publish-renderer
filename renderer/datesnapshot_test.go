package renderer

import (
	"testing"
	"time"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/localstore/addr"
	"github.com/stretchr/testify/require"
)

func TestGetDateSnapshot(t *testing.T) {
	r := &Renderer{}

	cases := []struct {
		name      string
		objectId  string
		expectErr bool
		expectVal string
	}{
		{
			name:      "Valid Today",
			objectId:  addr.DatePrefix + time.Now().Format("2006-01-02"),
			expectErr: false,
			expectVal: "Today",
		},
		{
			name:      "Valid Yesterday",
			objectId:  addr.DatePrefix + time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			expectErr: false,
			expectVal: "Yesterday",
		},
		{
			name:      "Valid Tomorrow",
			objectId:  addr.DatePrefix + time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			expectErr: false,
			expectVal: "Tomorrow",
		},
		{
			name:      "Valid Past Date",
			objectId:  addr.DatePrefix + "2020-05-15",
			expectErr: false,
			expectVal: "May 15, 2020",
		},
		{
			name:      "Invalid Date Format",
			objectId:  addr.DatePrefix + "abc-xyz",
			expectErr: true,
			expectVal: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := r.getDateSnapshot(tc.objectId)
			if tc.expectErr {
				require.Nil(t, result)
			} else {
				require.NotNil(t, result)
				require.Equal(t, tc.expectVal, result.Snapshot.Data.Details.Fields[bundle.RelationKeyName.String()].GetStringValue())
			}
		})
	}
}
