var      Mongo = require("mongodb").MongoClient,
       phantom = require("phantomjs");

var dbUrl = "mongodb://localhost:27017/local";
Mongo.connect(dbUrl, function(error, db) {
    if( error ) throw error;
    var collection = db.collection("test_images");
    collection.find().toArray(function(error, documents){
        if( error ) throw error;
        for( var i = 0; i < documents.length; i++ ){
            var item = documents[i];
            var imageExtension = ".jpg";
            var filename = item.sku + imageExtension;
            var filePath = __dirname;
        }
        db.close();
    });
});