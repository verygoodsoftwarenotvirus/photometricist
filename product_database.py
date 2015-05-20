import pymongo


def write_product_to_db(product):
    client = pymongo.MongoClient('localhost', 27017)
    db = client.local
    entries = db["test"]
    try:
        entries.insert_one(product)
    except pymongo.errors.DuplicateKeyError:
        pass


def update_product_in_db(product):
    # TODO
    return product