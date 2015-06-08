var Mongo = require("mongodb").MongoClient ,
    assert = require('assert');

var dbUrl = "mongodb://localhost:27017/local";

Mongo.connect(dbUrl, function(error, db) {
    assert.equal(null, error);
    var collection = db.collection("test_images");

    collection.remove();

    db.close();
});