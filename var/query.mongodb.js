/**
 * Check substations that have >1 switchboards in the list. To test if my code works
 */
use("gemini");
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
])

/**
 * Asset Hierarchy query
 */
use("gemini");
db.substations.aggregate([])

use("gemini");
db.switchboards.countDocuments();

use("gemini");
db.dropDatabase("gemini");
