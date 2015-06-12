var      Mongo = require("mongodb").MongoClient,
       request = require("request"),
         async = require("async");
            fs = require("fs"),
            CT = require("color-thief"),
    colorThief = new CT();

var download = function(uri, callback){
    request.head(uri, function(error){
        var r = request(uri).pipe(fs.createWriteStream("image.jpg"));
        r.on('close', callback);
    });
};

var rgb2hex = function (r, g, b) {
    return "#" + ((1 << 24)
               + (r << 16)
               + (g << 8) + b)
            .toString(16).slice(1);
};

var dbUrl = "mongodb://localhost:27017/local";
Mongo.connect(dbUrl, function(error, db) {
    if( error ) throw error;
    var collection = db.collection("test_images_copy");
    collection.find().toArray(function(error, documents){
        if( error ) throw error;
        async.series([], undefined /*db.close()*/);
        for( var i in documents ){
            var item = documents[i];
            console.dir(item);
            download(item.imageUrl, function(){
                var itemColors = colorThief.getPalette("image.jpg");
                var resultColors = [];
                for( var x in itemColors ) {
                    resultColors.push(
                        rgb2hex(itemColors[x][0],
                                itemColors[x][1],
                                itemColors[x][2])
                    );
                }
                console.dir(itemColors);
                collection.update({ id: item.sku },
                              { $set:
                                  { calculated_colors: resultColors }
                              }, function(error, result){
                    if( error ) {
                        console.log(item.sku);
                        console.log(error);
                    } else {
                        console.log(result);
                    }
                });
            });
        }
        db.close();
    });
});