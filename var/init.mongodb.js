function reset() {
    return db.dropDatabase("gemini");
}

/**
 * Cannot create a unique index because if we do not store a value for these indexed fields, Mongo will insert them as
 * `null` -- which eventually breaks the uniqueness constraint.
 *
 * This also means, we are unable to enforce the list items to be unique of each other as a constraint. More as an
 * application logic, then.
 *
 * https://www.mongodb.com/docs/manual/core/index-unique/#:~:text=%22%20%7D%20%5D%20%7D%20)-,Unique%20Index%20and%20Missing%20Field,that%20lacks%20the%20indexed%20field.
 */

function createIndexes() {
    db.substations.createIndex(
        {"assets.transformers._id": 1},
        {"unique": true},
    );
    db.substations.createIndex(
        {"assets.switchboards._id": 1},
        {"unique": true},
    );
}

function load() {
    return db.substations.insertMany([
        {
            "_id":  ObjectId("000000000000000000000000"),
            "name": "Substation #1",
            "assets": {
                "switchboards": [
                    {
                        "_id": ObjectId("000000000000000000000000"),
                        "status": "planned",
                    },
                    {
                        "_id": ObjectId("000000000000000000000002"),
                        "status": "planned",
                    }
                ]
            }
        },
        {
            "_id":  ObjectId("000000000000000000000001"),
            "name": "Substation #2",
            "assets": {
                "switchboards": [
                    {
                        "_id": ObjectId("000000000000000000000001"),
                        "status": "planned",
                    }
                ]
            }
        }
    ]);
}

function getSize(name) {
    const collection = db.getCollection(name);

    if (!collection)
        return { "ok": 0 };

    return collection.stats();
}

function loadSensors(count) {
    const
        configs = [],
        thresholds = [];

    for (var loop = 0; loop < count; loop++) {
        const configId = ObjectId()

        configs.push({
            "_id": configId,
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z",
            "deleted_at": "2023-01-01T00:00:00Z",
            "substation_id": ObjectId("000000000000000000000001"),
            "asset_id": ObjectId("000000000000000000000001"),
            "asset_type": "LOW_VOLTAGE_BOARD",
            "sensor_id": "12345678900000",
            "property": "LVB_RED_TEMPERATURE",
            "thresholds": {
                "min": 9999,
                "max": 9999,
                "started_at": "2023-01-01T00:00:00Z",
            }
        });

        thresholds.push({
            "_id": ObjectId(),
            "created_at": "2023-01-01T00:00:00Z",
            "updated_at": "2023-01-01T00:00:00Z",
            "ended_at": "2023-01-01T00:00:00Z",
            "sensor_config_id": configId,
            "min": 9999,
            "max": 9999,
        });
    }

    db.sensorConfigs.insertMany(configs);
    db.sensorThresholds.insertMany(thresholds);
}

use("gemini");
reset();
// loadSensors(100);
// getSize("sensorConfigs");
// getSize("sensorThresholds");
