/**
 * Upserting an a nested document does not really work for us in MongoDB.
 *      $push does not work the way I want it to be.
 *      The only way is I have to fetch the entire document and update it. This makes bulk writes a huge pain in the butt
 *      I would much rather do $addToSet with just an object ID instead.
 */
function upsertSeededSwitchboard() {
    return db.substations.updateOne(
        {
            "_id": ObjectId("000000000000000000000000"),
        },
        {
            // "$set": {
            //     "assets.switchboards.$[swb]._id": ObjectId("000000000000000000000002"),
            //     "assets.switchboards.$[swb].status": "seed update"
            // },
            "$addToSet": {
                "assets.switchboards": {
                    "_id": ObjectId("000000000000000000000003"),
                    "name": "oh wow2",
                    "status": "planned"
                }
            }
        },
        {
            "upsert": true,
            // "arrayFilters": [
            //     {"swb._id": ObjectId("000000000000000000000002")},
            // ],
        }
    );
}


function upsertImportedSwitchboard() {
    db.getCollection("switchboards").bulkWrite([
        {"updateOne": {
            "upsert": true,
            "filter": {
                "$or": [
                    {"_id": ObjectId("64ed57f5bb16e86671301bc6")},
                    {"asset_id": "DFNDXSS20060000170"},
                ]
            },
            "arrayFilters": [
                {"elem._id": ObjectId("000000000000000000000000")}
            ],
            "update": {
                "$setOnInsert": {
                    "_id": ObjectId("000000000000000000000000"),
                    "asset_id": "DNSWB01"
                },
                "$set": {
                    "assets.switchboards.$[elem].name": "Switchboard #1 22kV",
                    "assets.switchboards.$[elem].status": "commissioned",
                },
                "$currentDate": {
                    "assets.switchboards.$[elem].created_at": true,
                    "assets.switchboards.$[elem].updated_at": true,
                }
            }
        }},
    ]);
}

use("gemini");
upsertSeededSwitchboard();
