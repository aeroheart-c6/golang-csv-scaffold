use("gemini");
db.substations.find({
    // "_id": ObjectId("64ed57f5bb16e86671301bc6"),
    // "asset_id": "DFNDXSS20060000170",
    "$or": [
        {"_id": ObjectId("64ed57f5bb16e86671301bc6")},
        {"asset_id": "DFNDXSS20060000170"},
    ]
});

/**
 * Upserting an a nested document does not really work for us in MongoDB.
 *      $push does not work the way I want it to be.
 *      The only way is I have to fetch the entire document and update it. This makes bulk writes a huge pain in the butt
 *      I would much rather do $addToSet with just an object ID instead.
 */
use("gemini");
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

// db.getCollection("substations").bulkWrite([
//     {"updateOne": {
//         "upsert": true,
//         "filter": {
//             "_id": ObjectId("000000000000000000000000"),
//         },
//         "arrayFilters": [
//             {"elem._id": ObjectId("000000000000000000000000")}
//         ],
//         "update": {
//             "$set": {
//                 "assets.transformers.$[elem].name": "Transformer #1 22kV",
//                 "assets.transformers.$[elem].status": "commissioned",
//             },
//             "$currentDate": {
//                 "assets.transformers.$[elem].updated_at": true,
//             }
//         }
//     }},
//     {"updateOne": {
//         "upsert": true,
//         "filter": {
//             "_id": ObjectId("000000000000000000000001"),
//         },
//         "arrayFilters": [
//             {"elem._id": ObjectId("000000000000000000000001")}
//         ],
//         "update": {
//             "$set": {
//                 "assets.transformers.$[elem].name": "Transformer #2 66kV",
//                 "assets.transformers.$[elem].status": "commissioned",
//             },
//             "$currentDate": {
//                 "assets.transformers.$[elem].updated_at": true,
//             }
//         }
//     }},
// ]);
