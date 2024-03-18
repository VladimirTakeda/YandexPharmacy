package tests_go

import (
	"encoding/json"
	"github.com/antelman107/net-wait-go/wait"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	Endpoint = "http://localhost:8080/check"
)

var (
	Problems = []string{
		"WRONG_CATEGORY",
		"INCORRECT_ITEM_ID",
		"ITEM_NOT_FOUND",
		"NO_USER",
		"NO_USER_NO_RECEIPT",
		"NO_USER_SPECIAL_ITEM",
		"NO_RECEIPT",
		"ITEM_IS_SPECIAL",
		"ITEM_SPECIAL_WRONG_SPECIFIC",
	}
)

func hasUppercaseLetters(text string) bool {
	for _, ch := range text {
		if ch >= 'A' && ch <= 'Z' {
			return true
		}
	}
	return false
}

func problemExists(problem string) bool {
	for _, p := range Problems {
		if p == problem {
			return true
		}
	}
	return false
}

func TestCheck(t *testing.T) {
	testCases := []struct {
		UserID int
		Items  map[string]string
	}{
		{
			UserID: 100500,
			Items: map[string]string{
				"blah_123": "WRONG_CATEGORY",
			},
		},
		{
			UserID: 100500,
			Items: map[string]string{
				"common_blah": "INCORRECT_ITEM_ID",
			},
		},
		{
			UserID: 100500,
			Items: map[string]string{
				"common_8": "NO_USER",
			},
		},
		{
			UserID: 100500,
			Items: map[string]string{
				"receipt_8": "NO_USER_NO_RECEIPT",
			},
		},
		{
			UserID: 100500,
			Items: map[string]string{
				"special_8": "NO_USER_SPECIAL_ITEM",
			},
		},
		{
			UserID: 4,
			Items: map[string]string{
				"receipt_11": "NO_RECEIPT",
			},
		},
		{
			UserID: 5,
			Items: map[string]string{
				"special_8": "ITEM_IS_SPECIAL",
			},
		},
		{
			UserID: 63,
			Items: map[string]string{
				"special_24": "ITEM_SPECIAL_WRONG_SPECIFIC",
			},
		},
		{
			UserID: 63,
			Items: map[string]string{
				"special_52":  "",
				"common_1234": "ITEM_NOT_FOUND",
			},
		},
		{
			UserID: 4,
			Items: map[string]string{
				"common_5":     "",
				"special_1234": "ITEM_NOT_FOUND",
			},
		},
		{
			UserID: 4,
			Items: map[string]string{
				"receipt_68":   "",
				"special_1234": "ITEM_NOT_FOUND",
			},
		},
		{
			UserID: 63,
			Items: map[string]string{
				"common_5":   "",
				"receipt_68": "",
				"special_52": "",
				"special_24": "ITEM_SPECIAL_WRONG_SPECIFIC",
			},
		},
		{
			UserID: 4,
			Items: map[string]string{
				"common_5":   "",
				"receipt_68": "",
				"receipt_11": "NO_RECEIPT",
				"special_8":  "ITEM_IS_SPECIAL",
			},
		},
		{
			UserID: 4,
			Items: map[string]string{
				"Common_4":     "",
				"Special_1111": "ITEM_NOT_FOUND",
				"Common_1235":  "ITEM_NOT_FOUND",
			},
		},
	}

	for _, tc := range testCases {

		if !wait.New(
			wait.WithProto("tcp"),
			wait.WithWait(200*time.Millisecond),
			wait.WithBreak(50*time.Millisecond),
			wait.WithDeadline(15*time.Second),
			wait.WithDebug(true),
		).Do([]string{"localhost:8080"}) {
			logrus.Fatalf("service is not available")
			return
		}

		t.Run("", func(t *testing.T) {
			params := url.Values{
				"user_id": {strconv.Itoa(tc.UserID)},
			}
			lowerCaseItems := make(map[string]string, 0)
			for itemID, value := range tc.Items {
				params.Add("item_id", itemID)
				lowerCaseItems[strings.ToLower(itemID)] = value
			}

			resp, err := http.Get(Endpoint + "?" + params.Encode())
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var responseItems []map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&responseItems)
			assert.NoError(t, err)

			responseItemIDs := make(map[string]bool)
			for _, respItem := range responseItems {
				respItemID, ok := respItem["item_id"].(string)
				assert.True(t, ok, "Response items has no \"item_id\" field")
				respProblem, ok := respItem["problem"].(string)
				assert.True(t, ok, "Response item has no \"problem\" field")
				assert.False(t, hasUppercaseLetters(respItemID), "item_id in response has uppercase letters")
				assert.True(t, problemExists(respProblem), "Unexpected problem code")

				assert.Contains(t, lowerCaseItems, respItemID, "Got item_id that was not in the request")

				if expectedProblem := lowerCaseItems[respItemID]; expectedProblem != "" {
					assert.Equal(t, expectedProblem, respProblem, "Incorrect problem for item_id")
				}

				responseItemIDs[respItemID] = true
			}

			expectedInResponse := make(map[string]bool)
			for itemID, expectedProblem := range lowerCaseItems {
				if expectedProblem != "" {
					expectedInResponse[itemID] = true
				}
			}

			missingItemIDs := make([]string, 0)
			for itemID := range expectedInResponse {
				if !responseItemIDs[itemID] {
					missingItemIDs = append(missingItemIDs, itemID)
				}
			}
			assert.Empty(t, missingItemIDs, "Some item_ids are missing in response")

			expectedNotInResponse := make(map[string]bool)
			for itemID, expectedProblem := range lowerCaseItems {
				if expectedProblem == "" {
					expectedNotInResponse[itemID] = true
				}
			}

			extraItemIDs := make([]string, 0)
			for itemID := range responseItemIDs {
				if expectedNotInResponse[itemID] {
					extraItemIDs = append(extraItemIDs, itemID)
				}
			}
			assert.Empty(t, extraItemIDs, "Got unexpected item_ids in response")
		})
	}
}
