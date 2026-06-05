package store

import (
	"testing"
	"time"
)

func TestStore_GetAllItemsOrder(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name          string
		ordering      string
		items         []Item
		expectedOrder []string
	}{
		{
			name:     "test descending_by_new",
			ordering: "desc",
			items: []Item{
				{
					PublishedAt: now.Add(-2 * time.Hour),
					CreatedAt:   now.Add(-2 * time.Hour),
					Title:       "First Item",
					Content:     "A short string",
					FeedURL:     "example.com",
					Link:        "example.com/endpoint",
				},
				{
					PublishedAt: now.Add(-23 * time.Hour),
					CreatedAt:   now.Add(-24 * time.Hour),
					Title:       "Second Item",
					Content:     "A really long string",
					FeedURL:     "example.com",
					Link:        "example.com/api",
				},
				{
					PublishedAt: now.Add(-20 * time.Hour),
					CreatedAt:   now.Add(-20 * time.Hour),
					Title:       "Third Item",
					Content:     "Some boring string. Use lorem ipsum",
					FeedURL:     "example.com",
					Link:        "example.com/anotherExample",
				},
			},
			expectedOrder: []string{"First Item", "Second Item", "Third Item"},
		},
		{
			name:     "test_order_by_content_size",
			ordering: "length",
			items: []Item{
				{
					PublishedAt: now.Add(-2 * time.Hour),
					CreatedAt:   now.Add(-2 * time.Hour),
					Title:       "First Item",
					Content:     "AA",
					FeedURL:     "example.com",
					Link:        "example.com/endpoint",
				},
				{
					PublishedAt: now.Add(-23 * time.Hour),
					CreatedAt:   now.Add(-24 * time.Hour),
					Title:       "Second Item",
					Content:     "AAA",
					FeedURL:     "example.com",
					Link:        "example.com/api",
				},
				{
					PublishedAt: now.Add(-20 * time.Hour),
					CreatedAt:   now.Add(-20 * time.Hour),
					Title:       "Third Item",
					Content:     "A",
					FeedURL:     "example.com",
					Link:        "example.com/anotherExample",
				},
			},
			expectedOrder: []string{"First Item", "Second Item", "Third Item"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sls, err := NewInMemorySQLiteStore()
			if err != nil {
				t.Fatalf("Error initailizing db: %v", err)
			}
			for _, item := range tc.items {
				err = sls.UpsertItem(&item)
				if item.Title == "Third Item" {
					sls.ToggleRead(item.ID)
				}
				if err != nil {
					t.Errorf("Error insert %v: %v", item, err)
				}
			}
			allItems, err := sls.GetAllItems(tc.ordering)
			if err != nil {
				t.Fatalf("Error fetching all items: %v", err)
			}
			for i := range allItems {
				if allItems[i].Title != tc.expectedOrder[i] {
					t.Fatalf("Unexpected item order. Want: %v Got: %v", tc.expectedOrder[i], allItems[i].Title)
				}
			}

		})
	}
}
