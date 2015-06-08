var Mongo = require("mongodb").MongoClient ,
    assert = require('assert');

var dbUrl = "mongodb://localhost:27017/local";

var items = [
    {
        "2874842": {
            title: "Cabana Poms Lumbar Pillow - Cobalt",
            itemUrl: "http://www.pier1.com/Cabana-Poms-Lumbar-Pillow---Cobalt/2874842,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2874842/2874842_1.jpg?sw=1200&sh=1200",
        }
    }, {
        "2879481": {
            title: "Calliope Button Lumbar Pillow - Clay",
            itemUrl: "http://www.pier1.com/Calliope-Button-Lumbar-Pillow---Clay/2879481,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2874842/2874842_1.jpg?sw=1200&sh=1200",
        }
    }, {
        "2745465": {
            title: "Cabana Petal Lumbar Pillow - Orange",
            itemUrl: "http://www.pier1.com/Cabana-Petal-Lumbar-Pillow---Orange/2745465,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2745465/2745465_1.jpg?sw=1200&sh=1200",
        }
    }, {
        "2520402": {
            title: "Flounce Pillow - Purple",
            itemUrl: "http://www.pier1.com/Flounce-Pillow---Purple/2520402,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2520402/2520402_1.jpg?sw=1200&sh=1200",
        }
    }, {
        "2714003": {
            title: "Plush Pillow - Red",
            itemUrl: "http://www.pier1.com/Plush-Pillow---Red/2714003,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2741348/2741348_1.jpg?sw=1200&sh=1200",
        }
    }, {
        "2741348": {
            title: "Cabana Pillow - Citrus",
            itemUrl: "http://www.pier1.com/Cabana-Pillow---Citrus/2741348,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2741348/2741348_1.jpg?sw=1200&sh=1200",
        }
    }, {
        "2911526": {
            title: "Herringbone Chenille Pillow - Indigo",
            itemUrl: "http://www.pier1.com/Herringbone-Chenille-Pillow---Indigo/2911526,default,pd.html",
            imageUrl: "http://demandware.edgesuite.net/sits_pod20/dw/image/v2/AAID_PRD/on/demandware.static/Sites-pier1_us-Site/Sites-pier1_master/default/v1430339361406/images/2911526/2911526_1.jpg?sw=1200&sh=1200",
        }
    }
]

Mongo.connect(dbUrl, function(error, db) {
    assert.equal(null, error);
    var collection = db.collection("test_images");
    /*
    for( var item in items ) {
        collection.insert(item, function(error, result) {
            assert.equal(error, null);
        });
    };
    */
    collection.insert(items);

    db.close();
});