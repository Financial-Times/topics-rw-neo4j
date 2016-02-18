package topics

import (
	"fmt"
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

var topicsDriver baseftrwapp.Service

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"

	topicsDriver = getTopicsCypherDriver(t)

	topicToDelete := Topic{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(topicsDriver.Write(topicToDelete), "Failed to write topic")

	found, err := topicsDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete topic for uuid %", uuid)
	assert.NoError(err, "Error deleting topic for uuid %s", uuid)

	p, found, err := topicsDriver.Read(uuid)

	assert.Equal(Topic{}, p, "Found topic %s who should have been deleted", p)
	assert.False(found, "Found topic for uuid %s who should have been deleted", uuid)
	assert.NoError(err, "Error trying to find topic for uuid %s", uuid)
}

func TestCreateAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	topicsDriver = getTopicsCypherDriver(t)

	topicToWrite := Topic{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")

	readTopicForUUIDAndCheckFieldsMatch(t, uuid, topicToWrite)

	cleanUp(t, uuid)
}

func TestCreateHandlesSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	topicsDriver = getTopicsCypherDriver(t)

	topicToWrite := Topic{UUID: uuid, CanonicalName: "Test 'special chars", TmeIdentifier: "TME_ID"}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")

	readTopicForUUIDAndCheckFieldsMatch(t, uuid, topicToWrite)

	cleanUp(t, uuid)
}

func TestCreateNotAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	topicsDriver = getTopicsCypherDriver(t)

	topicToWrite := Topic{UUID: uuid, CanonicalName: "Test"}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")

	readTopicForUUIDAndCheckFieldsMatch(t, uuid, topicToWrite)

	cleanUp(t, uuid)
}

func TestUpdateWillRemovePropertiesNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	topicsDriver = getTopicsCypherDriver(t)

	topicToWrite := Topic{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")
	readTopicForUUIDAndCheckFieldsMatch(t, uuid, topicToWrite)

	updatedTopic := Topic{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(topicsDriver.Write(updatedTopic), "Failed to write updated topic")
	readTopicForUUIDAndCheckFieldsMatch(t, uuid, updatedTopic)

	cleanUp(t, uuid)
}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	topicsDriver = getTopicsCypherDriver(t)
	err := topicsDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
}

func getTopicsCypherDriver(t *testing.T) service {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(url)
	assert.NoError(err, "Failed to connect to Neo4j")
	return NewCypherTopicsService(neoutils.StringerDb{db}, db)
}

func readTopicForUUIDAndCheckFieldsMatch(t *testing.T, uuid string, expectedTopic Topic) {
	assert := assert.New(t)
	storedTopic, found, err := topicsDriver.Read(uuid)

	assert.NoError(err, "Error finding topic for uuid %s", uuid)
	assert.True(found, "Didn't find topic for uuid %s", uuid)
	assert.Equal(expectedTopic, storedTopic, "topics should be the same")
}

func TestWritePrefLabelIsAlsoWrittenAndIsEqualToName(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)
	uuid := "12345"
	topicToWrite := Topic{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	storedTopic := topicsDriver.Write(topicToWrite)

	fmt.Printf("", storedTopic)

	result := []struct {
		PrefLabel string `json:"t.prefLabel"`
	}{}

	getPrefLabelQuery := &neoism.CypherQuery{
		Statement: `
				MATCH (t:Topic {uuid:"12345"}) RETURN t.prefLabel
				`,
		Result: &result,
	}

	err := topicsDriver.cypherRunner.CypherBatch([]*neoism.CypherQuery{getPrefLabelQuery})
	assert.NoError(err)
	assert.Equal("Test", result[0].PrefLabel, "PrefLabel should be 'Test")
	cleanUp(t, uuid)
}

func cleanUp(t *testing.T, uuid string) {
	assert := assert.New(t)
	found, err := topicsDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete topic for uuid %", uuid)
	assert.NoError(err, "Error deleting topic for uuid %s", uuid)
}
