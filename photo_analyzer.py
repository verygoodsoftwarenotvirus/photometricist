import os
import json
import photo_retriever
import photo_functions


def get_config_from_file(file_location="config.json"):
    with open(file_location) as config_file:
        config = json.load(config_file)

    return config


def main():
    config = get_config_from_file("config.json")
    source_file = config["file"]["source_file"]
    photo_column = config["file"]["photo_info_column"]
    cropped_folder = config["photos"]["cropped_photo_dir"]
    photo_destination_folder = config["photos"]["base_photo_dir"]
    crop_percentage = config["photos"]["crop_percentage"]

    if not os.path.isdir(photo_destination_folder):
        os.mkdir(photo_destination_folder)
    if not os.path.isdir(cropped_folder):
        os.mkdir(cropped_folder)

    photo_retriever.process_file(source_file, photo_column, photo_destination_folder)
    saved_photos = os.listdir(photo_destination_folder)
    for file in saved_photos:
        file_path = "{0}/{1}".format(photo_destination_folder, file)
        image = photo_functions.open_image(file_path)
        image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
        file_path = "{0}/{1}".format(cropped_folder, file)
        photo_functions.save_image(image, file_path)

if __name__ == "__main__":
    main()