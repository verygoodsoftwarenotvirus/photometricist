import os
import csv
import json
import photo_retriever
import photo_functions
import color_analysis


def get_config_from_file(file_location="config.json"):
    with open(file_location) as config_file:
        config = json.load(config_file)

    return config


def main():
    config = get_config_from_file("config.json")
    source_file = config["file"]["source_file"]
    save_as = config["file"]["save_as"]
    photo_column = config["file"]["photo_info_column"]
    cropped_folder = config["photos"]["cropped_photo_dir"]
    photo_destination_folder = config["photos"]["base_photo_dir"]
    crop_percentage = config["photos"]["crop_percentage"]

    if not os.path.isdir(photo_destination_folder):
        os.mkdir(photo_destination_folder)
    if not os.path.isdir(cropped_folder):
        os.mkdir(cropped_folder)

    with open(source_file) as source, open(save_as, "w") as output:
        reader = csv.DictReader(source)
        fieldnames = reader.fieldnames
        fieldnames.append("computed_colors")
        writer = csv.DictWriter(output, fieldnames)
        writer.writeheader()
        for row in reader:
            new_row = row
            for key, value in row.items():
                new_row.update({key: value})

            photo_link = row[photo_column]
            photo_file_name = photo_destination_folder + photo_link[photo_link.rfind("/"):photo_link.rfind("?")]
            photo_retriever.retrieve_photos(photo_link, photo_file_name)

            image = photo_functions.open_image(photo_file_name)
            image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
            photo_functions.save_image(image, photo_file_name)

            results = color_analysis.analyze_color(image, 3)
            new_row["computed_colors"] = results
            writer.writerow(row)

if __name__ == "__main__":
    main()