import csv
import requests


def retrieve_photo(photo_url, target_file_name):
    data = requests.get(photo_url)
    with open(target_file_name, 'wb') as f:
        for chunk in data.iter_content():
            f.write(chunk)


def retrieve_photos_from_file(conf):
    downloaded_photo_paths = []
    with open(conf["source_file"], encoding=conf["input_encoding"]) as source:
        reader = csv.DictReader(source)
        for row in reader:
            photo_link = row[conf["photo_column"]]
            if not photo_link:
                continue
            photo_path = "{0}{1}".format(conf["photo_destination_folder"],
                                         photo_link[photo_link.rfind("/"):photo_link.rfind("?")])
            retrieve_photo(photo_link, photo_path)
            downloaded_photo_paths.append(photo_path)
    return downloaded_photo_paths