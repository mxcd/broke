package outline

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserClient(t *testing.T) {
	client, err := NewOutlineClient(ClientConfig{
		Url:   "https://wiki.fsintra.net/api",
		Token: "ol_api_uNF5qNh4fo9rvkxcZwoqbuUaqCCusoWTIdf17C",
	})

	assert.NoError(t, err, "error creating outline client")

	userId, err := client.GetUserIdByMail("malte.meiners@rennteam-stuttgart.de")
	assert.NoError(t, err, "error getting user id")

	println(userId)
}
