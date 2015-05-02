import os
import csv
import json
import photo_retriever
import photo_functions
import color_analysis

FIVE_AMPS = "&&&&&"


def establish_csv_headers(config, source_file_headers):
    headers = []
    photo_column = config["file"]["photo_info_column"]
    if "sku" in [column.lower() for column in source_file_headers]:
        headers.append('sku')
    headers += [photo_column, "computed_colors",
                "computed_strategy_colors","confidence_index"]
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
               config["results"]["confidence_minimum"],
               config["results"]["desired_analysis_clusters"]]
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
    # minimum_confidence = config["results"]["confidence_minimum"]
    k = config["results"]["desired_analysis_clusters"]
    debugging = config.get("debug")

    if not os.path.isdir(photo_destination_folder):
        os.mkdir(photo_destination_folder)
    if not os.path.isdir(cropped_folder):
        os.mkdir(cropped_folder)

    if debugging:
        debug_folder = "test_images"
        test_files = os.listdir(debug_folder)
        with open("debug_results.csv", "w") as output:
            writer = csv.DictWriter(output, ["photo_file", "computed_colors"])
            for test_file in test_files:
                new_row = {"photo_file": test_file}
                photo_path = "{0}/{1}".format(debug_folder, test_file)
                image = photo_functions.open_image(photo_path)
                image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
                if cropped_images_are_to_be_saved:
                    photo_path = photo_path.replace(debug_folder, cropped_folder)
                    photo_functions.save_image(image, photo_path)
                results = color_analysis.analyze_color(image, k)
                new_row["computed_colors"] = results
                writer.writerow(new_row)
        # TODO: web-page-ish output.
    else:
        with open(source_file) as source, open(save_as, "w") as output:
            reader = csv.DictReader(source)
            fieldnames = establish_csv_headers(config, reader.fieldnames)
            writer = csv.DictWriter(output, fieldnames, lineterminator='\n')
            writer.writeheader()
            for row in reader:
                new_row = {}
                photo_link = row[photo_column]
                photo_path = photo_destination_folder + photo_link[photo_link.rfind("/"):photo_link.rfind("?")]
                photo_retriever.retrieve_photos(photo_link, photo_path)

                image = photo_functions.open_image(photo_path)
                image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
                if cropped_images_are_to_be_saved:
                    photo_path = photo_path.replace(photo_destination_folder, cropped_folder)
                    photo_functions.save_image(image, photo_path)

                results = color_analysis.analyze_color(image, k)
                new_row["computed_colors"] = results

                computed_strategy_colors = set(color_analysis.compute_color_matches(config, results))
                if computed_strategy_colors:
                    new_row["computed_strategy_colors"] = computed_strategy_colors
                writer.writerow(new_row)

if __name__ == "__main__":
    main()