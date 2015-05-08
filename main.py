import os
import json
import time
import csv_output
import html_output
import photo_retriever

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
    try:
        os.mkdir(conf["photo_destination_folder"])
        os.mkdir(conf["cropped_folder"])
    except FileExistsError:
        pass

    if conf["app_mode"].lower() == "color":
        if conf["output_format"].lower() == "csv":
            csv_output.output_csv(conf)
        elif conf["output_format"].lower() == "html":
            photos = photo_retriever.retrieve_photos_from_file(conf)
            html_output.output_html(conf, photos)
    elif conf["app_mode"].lower() == "shape":
        pass
    else:
        raise ValueError("Application mode '{0}' invalid.".format(conf["app_mode"]))

if __name__ == "__main__":
    main()
