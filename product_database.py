import psycopg2


def write_product_to_db():
    connection = psycopg2.connect(database="clients", user='curatorbot', host='localhost')
    cursor = connection.cursor()

    creation_phrase = "INSERT INTO products.products(sku, text_sku, curated_colors) VALUES('123456', 'fart', 'butts');"
    cursor.execute(creation_phrase)
    # thing = cursor.execute("SELECT * FROM pier_one.products WHERE sku = {0}".format(42069))
    # print(thing)


write_product_to_db()