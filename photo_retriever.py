import csv
import time
import requests


def retrieve_photos(photo_url, target_file_name):
    data = requests.get(photo_url)
    with open(target_file_name, 'wb') as f:
        for chunk in data.iter_content():
            f.write(chunk)


def process_file(filename, photo_column, destination_folder=time.time()):
    with open(filename, "r") as source:
        reader = csv.DictReader(source)
        fieldnames = reader.fieldnames
        fieldnames.append("computed_color")
        for row in reader:
            photo_link = row[photo_column]
            photo_file_name = destination_folder + photo_link[photo_link.rfind("/"):photo_link.rfind("?")]
            retrieve_photos(photo_link, photo_file_name)

