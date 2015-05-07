import os
import csv
import json
import time
import csv_output
import color_analysis
import photo_retriever
import photo_functions
import result_page_builder

FIVE_AMPS = "&&&&&"


def get_config_from_file(file_location="config.json"):
    try:
        with open(file_location) as config_file:
            configuration = json.load(config_file)
        return configuration
    except FileNotFoundError:
        print("Configuration file not found!")
        exit()


def ensure_valid_configuration(cfg):
    conf = {}
    try:
        # required parameters
        conf["app_mode"] = cfg["mode"]
        conf["source_file"] = cfg["file"]["source_file"]
        conf["photo_column"] = cfg["file"]["photo_info_column"]
        conf["colors"] = cfg["colors"]
    except KeyError as error:
        print("{0} not found in configuration file".format(error))
        exit()

    conf["source_file_dir"] = cfg["file"].get("source_file_location")
    if conf["source_file_dir"][:-1] == "/":
        conf["source_file"] = "{0}{1}".format(conf["source_file_dir"], conf["source_file"])
    elif conf["source_file_dir"]:
        conf["source_file"] = "{0}/{1}".format(conf["source_file_dir"], conf["source_file"])
    conf["input_encoding"] = cfg["file"].get("encoding", "utf-8")
    conf["output_format"] = cfg["file"].get("output_type", "html").lower()
    conf["save_as"] = "{0}.{1}".format(cfg["file"].get("save_as", str(time.time())), conf["output_format"])

    conf["cropped_folder"] = cfg["photos"].get("cropped_photo_dir", "cropped_photos").lower()
    conf["photo_destination_folder"] = cfg["photos"].get("base_photo_dir", "product_photos").lower()
    conf["crop_percentage"] = cfg["photos"].get("crop_percentage", 100)
    conf["k"] = cfg["results"].get("desired_analysis_clusters", 3)
    conf["cropped_images_are_to_be_saved"] = cfg["photos"].get("save_cropped_photos", True)

    # things to be implemented later.
    conf["minimum_confidence"] = cfg["results"].get("confidence_minimum", 0)
    conf["verbose"] = cfg.get("verbose", False)
    conf["debugging"] = cfg.get("debug", False)

    return conf


def main():
    conf = ensure_valid_configuration(get_config_from_file("config.json"))

    if conf["app_mode"].lower() == "color":
        pass
    elif conf["app_mode"].lower() == "shape":
        pass
    else:
        raise ValueError("Application mode '{0}' invalid.".format(conf["app_mode"]))

    try:
        os.mkdir(conf["photo_destination_folder"])
        os.mkdir(conf["cropped_folder"])
    except FileExistsError:
        pass

    try:
        with open(conf["source_file"], encoding=conf["input_encoding"]) as source, \
                open(conf["save_as"], "w") as output:
            reader = csv.DictReader(source)
            color_relationships = {}
            if conf["output_format"] == "csv":
                fieldnames = csv_output.establish_csv_headers(conf)
                writer = csv.DictWriter(output, fieldnames, lineterminator='\n')
                writer.writeheader()
            if conf["output_format"] == "html":
                results = {}
                crop_widths = []
            for row in reader:
                photo_link = row[conf["photo_column"]]
                if not photo_link:
                    continue
                photo_path = "{0}{1}".format(conf["photo_destination_folder"],
                                             photo_link[photo_link.rfind("/"):photo_link.rfind("?")])
                photo_retriever.retrieve_photo(photo_link, photo_path)

                # crop and save image
                image = photo_functions.open_image(photo_path)
                image = photo_functions.center_crop_image_by_percentage(image, conf["crop_percentage"])
                photo_path = photo_path.replace(conf["photo_destination_folder"], conf["cropped_folder"])
                photo_functions.save_image(image, photo_path)

                analysis_results = color_analysis.analyze_color(image, conf["k"])
                relationship_result = color_analysis.compute_color_matches(conf, analysis_results)
                computed_strategy_colors = []
                for relationship in relationship_result:
                    for color in relationship_result[relationship]:
                        computed_strategy_colors.append(color)
                color_relationships.update(relationship_result)
                if conf["output_format"] == "csv":
                    new_row = {}
                    for key, value in row.items():
                        if key in fieldnames:
                            new_row[key] = value
                    new_row["computed_colors"] = analysis_results
                    computed_strategy_colors = set(computed_strategy_colors)
                    if computed_strategy_colors:
                        new_row["computed_strategy_colors"] = computed_strategy_colors
                    writer.writerow(new_row)
                elif conf["output_format"] == "html":
                    b64_encoded = str(photo_functions.base64_encode_image(photo_path))[2:-1]
                    photo_extension = photo_path[photo_path.rfind(".") + 1:]
                    if photo_extension == "jpg":
                        photo_extension = "jpeg"
                    b64_prefix = "data:image/{0};base64,".format(photo_extension)
                    photo_path = "{0}{1}".format(b64_prefix, b64_encoded)
                    crop_widths.append(photo_functions.get_image_width(image))
                    results[photo_path] = analysis_results
            if conf["output_format"] == "html":
                html_output = result_page_builder.build_page(max(crop_widths) + 20,
                                                             results, color_relationships)
                output.write(html_output)
    except EnvironmentError:
        print("Something went awry.")


if __name__ == "__main__":
    main()
