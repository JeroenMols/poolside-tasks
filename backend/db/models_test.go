package db

import (
	"backend/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTodoItem_ChangeStatus(t *testing.T) {
	tests := []struct {
		oldStatus string
		newStatus string
		error     bool
	}{
		// Allowed transitions
		{
			oldStatus: "todo",
			newStatus: "ongoing",
			error:     false,
		},
		{
			oldStatus: "ongoing",
			newStatus: "done",
			error:     false,
		},
		{
			oldStatus: "done",
			newStatus: "ongoing",
			error:     false,
		},
		{
			oldStatus: "ongoing",
			newStatus: "todo",
			error:     false,
		},
		// Disallowed transitions
		{
			oldStatus: "todo",
			newStatus: "done",
			error:     true,
		},
		{
			oldStatus: "done",
			newStatus: "todo",
			error:     true,
		},
		{
			oldStatus: "todo",
			newStatus: "invalid",
			error:     true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("from %s to %s", tt.oldStatus, tt.newStatus), func(t *testing.T) {
			item := TodoDatabaseItem{
				Id:          "fake_id",
				UpdatedAt:   util.FakeTime(2021, 1, 1),
				Description: "fake description",
				Status:      tt.oldStatus,
				User:        "fake_user",
			}

			err := item.ChangeStatus(tt.newStatus)
			if err != nil || tt.error {
				expected := fmt.Sprintf("invalid status transition from %s to %s", tt.oldStatus, tt.newStatus)
				assert.Equal(t, expected, err.Error())
			} else {
				assert.Equal(t, tt.newStatus, item.Status)
			}
		})
	}
}
