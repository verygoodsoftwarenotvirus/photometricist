import os
import csv
import json
import time
import color_analysis
import photo_retriever
import photo_functions
import result_page_builder
from datetime import datetime

FIVE_AMPS = "&&&&&"


def establish_csv_headers(config, source_file_headers):
    headers = []
    photo_column = config["file"]["photo_info_column"]
    if "sku" in [column.lower() for column in source_file_headers]:
        headers.append('sku')
    headers += [photo_column, "computed_colors",
                "computed_strategy_colors", "confidence_index"]
    return headers


def validate_configuration(config):
    strings = [config["file"]["encoding"],
               config["file"]["source_file"],
               config["file"]["save_as"],
               config["file"]["output_type"],
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

    booleans = [config["debug"],
                config["photos"]["save_cropped_photos"]]
    for boolean in booleans:
        if not isinstance(boolean, bool):
            return False

    # TODO: ensure variables meet standards.
    return True


def get_config_from_file(file_location="config.json"):
    with open(file_location) as config_file:
        config = json.load(config_file)

    return config


def main():
    config = get_config_from_file("config.json")
    config_valid = validate_configuration(config)
    if not config_valid:
        raise ValueError("Configuration values are invalid.")

    source_file = config["file"]["source_file"]
    source_file_encoding = config["file"]["encoding"]
    photo_column = config["file"]["photo_info_column"]
    cropped_images_are_to_be_saved = config["photos"]["save_cropped_photos"]
    cropped_folder = config["photos"]["cropped_photo_dir"]
    photo_destination_folder = config["photos"]["base_photo_dir"]
    crop_percentage = config["photos"]["crop_percentage"]
    minimum_confidence = config["results"]["confidence_minimum"]
    k = config["results"]["desired_analysis_clusters"]
    output_format = config["file"]["output_type"].lower()
    if "save_as" in config["file"]:
        save_as = "{0}.{1}".format(config["file"]["save_as"], output_format)
    else:
        save_as = "{0}.{1}".format(time.time(), output_format)
    debugging = config.get("debug")

    if not os.path.isdir(photo_destination_folder):
        os.mkdir(photo_destination_folder)
    if not os.path.isdir(cropped_folder):
        os.mkdir(cropped_folder)

    if debugging:
        debug_folder = "test_images"
        test_files = os.listdir(debug_folder)
        results = {}
        for test_file in test_files:
            photo_path = "{0}/{1}".format(debug_folder, test_file)
            image = photo_functions.open_image(photo_path)
            image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
            if cropped_images_are_to_be_saved:
                photo_path = photo_path.replace(debug_folder, cropped_folder)
                photo_functions.save_image(image, photo_path)
            analysis_result = color_analysis.analyze_color(image, k)
            results[photo_path] = analysis_result
    else:
        with open(source_file, encoding=source_file_encoding) as source, open(save_as, "w") as output:
            reader = csv.DictReader(source)
            color_relationships = {}
            if output_format == "csv":
                fieldnames = establish_csv_headers(config, reader.fieldnames)
                writer = csv.DictWriter(output, fieldnames, lineterminator='\n')
                writer.writeheader()
            if output_format == "html":
                results = {}
                crop_widths = []
            for row in reader:
                photo_link = row[photo_column]
                if not photo_link:
                    continue
                photo_path = photo_destination_folder + photo_link[photo_link.rfind("/"):photo_link.rfind("?")]
                photo_retriever.retrieve_photos(photo_link, photo_path)

                image = photo_functions.open_image(photo_path)
                image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
                if cropped_images_are_to_be_saved:
                    photo_path = photo_path.replace(photo_destination_folder, cropped_folder)
                    photo_functions.save_image(image, photo_path)

                analysis_results = color_analysis.analyze_color(image, k)
                relationship_result = color_analysis.compute_color_matches(config, analysis_results)
                computed_strategy_colors = []
                for relationship in relationship_result:
                    for color in relationship_result[relationship]:
                        computed_strategy_colors.append(color)
                color_relationships.update(relationship_result)
                if output_format == "csv":
                    new_row = {}
                    for key, value in row.items():
                        if key in fieldnames:
                            new_row[key] = value
                    new_row["computed_colors"] = analysis_results
                    computed_strategy_colors = set(computed_strategy_colors)
                    if computed_strategy_colors:
                        new_row["computed_strategy_colors"] = computed_strategy_colors
                    writer.writerow(new_row)
                elif output_format == "html":
                    crop_widths.append(photo_functions.get_image_width(image))
                    results[photo_path] = analysis_results
            if output_format == "html":
                html_output = result_page_builder.build_page(max(crop_widths)+20, results, color_relationships)
                output.write(html_output)

if __name__ == "__main__":
    start_time = datetime.now()
    print("Started at: {0}".format(start_time))
    main()
    finish_time = datetime.now() - start_time
    print("End time: {0}\n".format(datetime.now()))
    print("Script complete, elapsed time: {0}".format(finish_time))