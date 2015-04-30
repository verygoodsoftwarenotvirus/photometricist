import csv
import requests


args = {"source_file": "before_color.csv",
        "save_as": "after_color.csv",
        "photo_link_column": "image_url",
        "photo_directory": "saved_files"}


def retrieve_photos(photo_url, target_file_name):
    data = requests.get(photo_url)
    with open(target_file_name, 'wb') as f:
        for chunk in data.iter_content():
            f.write(chunk)


def process_file(filename):
    with open(args["source_file"], "r") as source, open(args["save_as"], "w") as new_file:
        reader = csv.DictReader(source)
        fieldnames = reader.fieldnames
        fieldnames.append("computed_color")
        writer = csv.DictWriter(new_file, fieldnames)
        for row in reader:
            photo_link = row[args["photo_link_column"]]
            photo_file_name = args["photo_directory"] + photo_link[photo_link.rfind("/"):photo_link.rfind("?")]
            retrieve_photos(photo_link, photo_file_name)

