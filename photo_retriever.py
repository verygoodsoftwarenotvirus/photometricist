import csv
import time
import requests


def retrieve_photo(photo_url, target_file_name):
    data = requests.get(photo_url)
    with open(target_file_name, 'wb') as f:
        for chunk in data.iter_content():
            f.write(chunk)


def retrieve_photos_from_file(source_file, encoding, photo_column, dest_folder=str(time.time())):
    downloaded_photo_paths = []
    with open(source_file, encoding=encoding) as source:
        reader = csv.DictReader(source)
        for row in reader:
            photo_link = row[photo_column]
            if not photo_link:
                continue
            photo_path = "{0}{1}".format(dest_folder,
                                         photo_link[photo_link.rfind("/"):photo_link.rfind("?")])
            retrieve_photo(photo_link, photo_path)
            downloaded_photo_paths.append(photo_path)
    return downloaded_photo_paths