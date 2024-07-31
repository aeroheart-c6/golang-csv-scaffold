use("gemini");
db.dropDatabase("gemini");

var collection = db.getCollection("substations");

collection.createIndex(
    {"assets.transformers._id": 1 },
    {"unique": true},
);

collection.insertMany([
    {
        "_id":  ObjectId("000000000000000000000000"),
        "name": "Substation #1",
        "assets": {
            "transformers": [
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
            "transformers": [
                {
                    "_id": ObjectId("000000000000000000000001"),
                    "status": "planned",
                }
            ]
        }
    }
]);
