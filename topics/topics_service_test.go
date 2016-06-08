package topics

import (
	"os"
	"testing"

	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

const (
	topicUUID            = "12345"
	newTopicUUID         = "123456"
	tmeID                = "TME_ID"
	newTmeID             = "NEW_TME_ID"
	fsetID               = "fset_ID"
	leiCodeID            = "leiCode"
	prefLabel            = "Test"
	specialCharPrefLabel = "Test 'special chars"
)

var defaultTypes = []string{"Thing", "Concept", "Topic"}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)
	err := topicsDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
}

func TestPrefLabelIsCorrectlyWritten(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{UUIDS: []string{topicUUID}}
	topicToWrite := Topic{UUID: topicUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	err := topicsDriver.Write(topicToWrite)
	assert.NoError(err, "ERROR happened during write time")

	storedTopic, found, err := topicsDriver.Read(topicUUID)
	assert.NoError(err, "ERROR happened during read time")
	assert.Equal(true, found)
	assert.NotEmpty(storedTopic)

	assert.Equal(prefLabel, storedTopic.(Topic).PrefLabel, "PrefLabel should be "+prefLabel)
	cleanUp(assert, topicUUID, topicsDriver)
}

func TestPrefLabelSpecialCharactersAreHandledByCreate(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{topicUUID}}
	topicToWrite := Topic{UUID: topicUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")

	//add default types that will be automatically added by the writer
	topicToWrite.Types = defaultTypes
	//check if topicToWrite is the same with the one inside the DB
	readTopicForUUIDAndCheckFieldsMatch(assert, topicsDriver, topicUUID, topicToWrite)
	cleanUp(assert, topicUUID, topicsDriver)
}

func TestCreateCompleteTopicWithPropsAndIdentifiers(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{topicUUID}, FactsetIdentifier: fsetID, LeiCode: leiCodeID}
	topicToWrite := Topic{UUID: topicUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")

	//add default types that will be automatically added by the writer
	topicToWrite.Types = defaultTypes
	//check if topicToWrite is the same with the one inside the DB
	readTopicForUUIDAndCheckFieldsMatch(assert, topicsDriver, topicUUID, topicToWrite)
	cleanUp(assert, topicUUID, topicsDriver)
}

func TestUpdateWillRemovePropertiesAndIdentifiersNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)

	allAlternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{topicUUID}, FactsetIdentifier: fsetID, LeiCode: leiCodeID}
	topicToWrite := Topic{UUID: topicUUID, PrefLabel: prefLabel, AlternativeIdentifiers: allAlternativeIdentifiers}

	assert.NoError(topicsDriver.Write(topicToWrite), "Failed to write topic")
	//add default types that will be automatically added by the writer
	topicToWrite.Types = defaultTypes
	readTopicForUUIDAndCheckFieldsMatch(assert, topicsDriver, topicUUID, topicToWrite)

	tmeAlternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{topicUUID}}
	updatedTopic := Topic{UUID: topicUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: tmeAlternativeIdentifiers}

	assert.NoError(topicsDriver.Write(updatedTopic), "Failed to write updated topic")
	//add default types that will be automatically added by the writer
	updatedTopic.Types = defaultTypes
	readTopicForUUIDAndCheckFieldsMatch(assert, topicsDriver, topicUUID, updatedTopic)

	cleanUp(assert, topicUUID, topicsDriver)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{topicUUID}}
	topicToDelete := Topic{UUID: topicUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(topicsDriver.Write(topicToDelete), "Failed to write topic")

	found, err := topicsDriver.Delete(topicUUID)
	assert.True(found, "Didn't manage to delete topic for uuid %", topicUUID)
	assert.NoError(err, "Error deleting topic for uuid %s", topicUUID)

	p, found, err := topicsDriver.Read(topicUUID)

	assert.Equal(Topic{}, p, "Found topic %s who should have been deleted", p)
	assert.False(found, "Found topic for uuid %s who should have been deleted", topicUUID)
	assert.NoError(err, "Error trying to find topic for uuid %s", topicUUID)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	topicsDriver := getTopicsCypherDriver(t)

	alternativeIds := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{topicUUID}}
	topicOneToCount := Topic{UUID: topicUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIds}

	assert.NoError(topicsDriver.Write(topicOneToCount), "Failed to write topic")

	nr, err := topicsDriver.Count()
	assert.Equal(1, nr, "Should be 1 topics in DB - count differs")
	assert.NoError(err, "An unexpected error occurred during count")

	newAlternativeIds := alternativeIdentifiers{TME: []string{newTmeID}, UUIDS: []string{newTopicUUID}}
	topicTwoToCount := Topic{UUID: newTopicUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: newAlternativeIds}

	assert.NoError(topicsDriver.Write(topicTwoToCount), "Failed to write topic")

	nr, err = topicsDriver.Count()
	assert.Equal(2, nr, "Should be 2 topics in DB - count differs")
	assert.NoError(err, "An unexpected error occurred during count")

	cleanUp(assert, topicUUID, topicsDriver)
	cleanUp(assert, newTopicUUID, topicsDriver)
}

func readTopicForUUIDAndCheckFieldsMatch(assert *assert.Assertions, topicsDriver service, uuid string, expectedTopic Topic) {

	storedTopic, found, err := topicsDriver.Read(uuid)

	assert.NoError(err, "Error finding topic for uuid %s", uuid)
	assert.True(found, "Didn't find topic for uuid %s", uuid)
	assert.Equal(expectedTopic, storedTopic, "topics should be the same")
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

func cleanUp(assert *assert.Assertions, uuid string, topicsDriver service) {
	found, err := topicsDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete topic for uuid %", uuid)
	assert.NoError(err, "Error deleting topic for uuid %s", uuid)
}
