import os
import csv
import json
import photo_retriever
import photo_functions
import color_analysis


def establish_csv_headers(config, source_file_headers):
    headers = []
    photo_column = config["file"]["photo_info_column"]
    if "sku" in [column.lower() for column in source_file_headers]:
        headers.append('sku')
    headers += [photo_column, "computed_color", "confidence_index"]
    return headers


def validate_configuration(config):
    strings = [config["file"]["source_file"],
               config["file"]["save_as"],
               config["file"]["photo_info_column"],
               config["photos"]["cropped_photo_dir"],
               config["photos"]["base_photo_dir"]]
    for string in strings:
        if not isinstance(string, str):
            return False

    numbers = [config["photos"]["crop_percentage"],
               config["results"]["confidence_minimum"]]
    for number in numbers:
        if not isinstance(number, (float, int)):
            return False

    booleans = [config["photos"]["save_cropped_photos"]]
    for boolean in booleans:
        if not isinstance(boolean, bool):
            return False

    return True


def get_config_from_file(file_location="config.json"):
    with open(file_location) as config_file:
        config = json.load(config_file)

    return config


def main():
    config = get_config_from_file("config.json")
    config_valid = validate_configuration(config)
    if not config_valid:
        raise ValueError("Configuration values are invalid, ")
    source_file = config["file"]["source_file"]
    save_as = config["file"]["save_as"]
    photo_column = config["file"]["photo_info_column"]
    cropped_images_are_to_be_saved = config["photos"]["save_cropped_photos"]
    cropped_folder = config["photos"]["cropped_photo_dir"]
    photo_destination_folder = config["photos"]["base_photo_dir"]
    crop_percentage = config["photos"]["crop_percentage"]
    minimum_confidence = config["results"]["confidence_minimum"]

    if not os.path.isdir(photo_destination_folder):
        os.mkdir(photo_destination_folder)
    if not os.path.isdir(cropped_folder):
        os.mkdir(cropped_folder)

    with open(source_file) as source, open(save_as, "w") as output:
        reader = csv.DictReader(source)
        fieldnames = establish_csv_headers(config, reader.fieldnames)
        writer = csv.DictWriter(output, fieldnames)
        writer.writeheader()
        for row in reader:
            new_row = row
            photo_link = row[photo_column]
            photo_path = photo_destination_folder + photo_link[photo_link.rfind("/"):photo_link.rfind("?")]
            photo_retriever.retrieve_photos(photo_link, photo_path)

            image = photo_functions.open_image(photo_path)
            image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
            if cropped_images_are_to_be_saved:
                photo_path = photo_path.replace(photo_destination_folder, cropped_folder)
                photo_functions.save_image(image, photo_path)

            results = color_analysis.analyze_color(image, 3)
            new_row["computed_colors"] = results

            computed_strategy_colors = []
            for color in config["colors"]:
                color_floor = config["colors"][color]["floor"]
                color_ceiling = config["colors"][color]["ceiling"]
                color_range = (color_floor, color_ceiling)
                for result in results:
                    if color_analysis.color_is_in_range(result, color_range):
                        computed_strategy_colors.append(color)

            computed_strategy_colors = set(computed_strategy_colors)
            if computed_strategy_colors:
                new_row["computed_strategy_colors"] = computed_strategy_colors
            writer.writerow(row)

if __name__ == "__main__":
    main()