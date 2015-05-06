import os
import csv
import json
import time
import color_analysis
import photo_retriever
import photo_functions
import result_page_builder

FIVE_AMPS = "&&&&&"


def establish_csv_headers(configuration, source_file_headers):
    headers = []
    try:
        photo_column = configuration["file"]["photo_info_column"]
        if "sku" in [column.lower() for column in source_file_headers]:
            headers.append('sku')
        headers += [photo_column, "computed_colors",
                    "computed_strategy_colors", "confidence_index"]
        return headers
    except KeyError:
        print("Photo information column not specified")
        exit()


def get_config_from_file(file_location="config.json"):
    try:
        with open(file_location) as config_file:
            configuration = json.load(config_file)
        return configuration
    except FileNotFoundError:
        print("Configuration file not found!")
        exit()


def main():
    configuration = get_config_from_file("config.json")
    try:
        # required parameters
        source_file = configuration["file"]["source_file"]
        photo_column = configuration["file"]["photo_info_column"]
        colors = configuration["colors"]
    except KeyError as error:
        print("{0} not found in configuration file".format(error))
        exit()

    source_file_dir = configuration["file"].get("source_file_location")
    if source_file_dir[:-1] == "/":
        source_file = "{0}{1}".format(source_file_dir, source_file)
    elif source_file_dir:
        source_file = "{0}/{1}".format(source_file_dir, source_file)
    input_encoding = configuration["file"].get("encoding", "utf-8")
    output_format = configuration["file"].get("output_type", "html").lower()
    save_as = "{0}.{1}".format(configuration["file"].get("save_as", str(time.time())), output_format)

    cropped_folder = configuration["photos"].get("cropped_photo_dir", "cropped_photos").lower()
    photo_destination_folder = configuration["photos"].get("base_photo_dir", "product_photos").lower()
    crop_percentage = configuration["photos"].get("crop_percentage", 100)
    k = configuration["results"].get("desired_analysis_clusters", 3)
    minimum_confidence = configuration["results"].get("confidence_minimum", 0)

    verbose = configuration.get("verbose", False)
    debugging = configuration.get("debug", False)
    cropped_images_are_to_be_saved = configuration["photos"].get("save_cropped_photos", True)
    product_images_are_to_be_saved = configuration["photos"].get("save_full_photos", True)

    try:
        os.mkdir(photo_destination_folder)
        os.mkdir(cropped_folder)
    except FileExistsError:
        pass
    try:
        with open(source_file, encoding=input_encoding) as source, open(save_as, "w") as output:
            reader = csv.DictReader(source)
            color_relationships = {}
            if output_format == "csv":
                fieldnames = establish_csv_headers(configuration, reader.fieldnames)
                writer = csv.DictWriter(output, fieldnames, lineterminator='\n')
                writer.writeheader()
            if output_format == "html":
                results = {}
                crop_widths = []
            for row in reader:
                photo_link = row[photo_column]
                if not photo_link:
                    continue
                photo_path = "{0}{1}".format(photo_destination_folder,
                                             photo_link[photo_link.rfind("/"):photo_link.rfind("?")])
                photo_retriever.retrieve_photos(photo_link, photo_path)

                image = photo_functions.open_image(photo_path)
                image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
                if cropped_images_are_to_be_saved:
                    photo_path = photo_path.replace(photo_destination_folder, cropped_folder)
                    photo_functions.save_image(image, photo_path)

                analysis_results = color_analysis.analyze_color(image, k)
                relationship_result = color_analysis.compute_color_matches(configuration, analysis_results)
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
                    b64_encoded = str(photo_functions.base64_encode_image(photo_path))[2:-1]
                    photo_extension = photo_path[photo_path.rfind(".")+1:]
                    photo_path = "{0}{1}".format("data:image/jpeg;base64,",
                                                 b64_encoded)
                    crop_widths.append(photo_functions.get_image_width(image))
                    results[photo_path] = analysis_results
            if output_format == "html":
                html_output = result_page_builder.build_page(max(crop_widths) + 20, results, color_relationships)
                output.write(html_output)
    except EnvironmentError:
        print("Something went awry.")


if __name__ == "__main__":
    main()
