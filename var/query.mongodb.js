/**
 * Check substations that have >1 switchboards in the list. To test if my code works
 */
function countSwitchboards() {
    db.substations.aggregate([
        {
            "$project": {
                "switchboardCount": { "$size": "$assets.switchboards" },
                "assets.switchboards": 1,
            },
        },
        {
            "$lookup": {
                "from": "switchboards",
                "localField": "assets.switchboards._id",
                "foreignField": "_id",
                "as": "switchboards",
            }
        },
        {
            "$match": {
                "switchboardCount": { "$gte": 2 },
            }
        },
    ]);
}

/**
 *
 */
function countSubstations() {
    return db.substations.countDocuments();
}

/**
 * Asset Hierarchy query
 */
function getHierarchy() {
    return db.substations.find();
}


function findImportedSubstation() {
    return db.substations.find({
        // "_id": ObjectId("64ed57f5bb16e86671301bc6"),
        // "asset_id": "DFNDXSS20060000170",
        "$or": [
            {"_id": ObjectId("64ed57f5bb16e86671301bc6")},
            {"asset_id": "DFNDXSS20060000170"},
        ]
    });
}


use("gemini");
getHierarchy();
