import os
import json
import time
import shutil
import logging
import color_analysis
import photo_retriever
from html_results_page_builder import build_results_page


def establish_logger():
    conf = get_config_from_file("config.json")
    if conf.get("verbose", False):
        logging.info("Logging established.")
        logging.basicConfig(level=logging.INFO,
                            format='%(asctime)s %(levelname)s: %(message)s',
                            datefmt='%m/%d/%y %I:%M:%S %p',
                            filename='image_curator.log',
                            filemode='w')
        console = logging.StreamHandler()
        console.setLevel(logging.INFO)
        formatter = logging.Formatter('%(asctime)s %(levelname)s: %(message)s',
                                      datefmt='%m/%d/%y %I:%M:%S %p',)
        console.setFormatter(formatter)
        logging.getLogger().addHandler(console)


def log_configuration_values(conf):
    logging.info("Source file path set to: {0}".format(conf["source_file"]))
    logging.info("Input file encoding set to: {0}".format(conf["input_encoding"]))
    logging.info("Output file name set to: {0}".format(conf["save_as"]))
    logging.info("Downloaded photos will be cropped to {0}%".format(conf["crop_percentage"]))
    logging.info("Photos analyzed will be broken into {0} clusters".format(conf["k"]))

    if conf["save_cropped_photos"]:
        logging.info("Cropped photos will be saved in {0}.".format(conf["cropped_folder"]))
    else:
        logging.info("Cropped photos will be discarded")

    if conf["save_product_photos"]:
        logging.info("Product photos will be saved in {0}.".format(conf["product_folder"]))
    else:
        logging.info("Product photos will be discarded")



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
    try:  # required parameters
        conf["app_mode"] = cfg["mode"]
        conf["source_file"] = cfg["file"]["source_file"]
        conf["photo_column"] = cfg["file"]["photo_info_column"]
        conf["colors"] = cfg["colors"]
    except KeyError as error:
        logging.getLogger().setLevel(logging.CRITICAL)
        logging.critical("{0} not found in configuration file".format(error))
        exit()

    conf["source_file_dir"] = cfg["file"].get("source_file_location")
    if conf["source_file_dir"][:-1] == "/":
        conf["source_file"] = "{0}{1}".format(conf["source_file_dir"], conf["source_file"])
    elif conf["source_file_dir"]:
        conf["source_file"] = "{0}/{1}".format(conf["source_file_dir"], conf["source_file"])

    conf["verbose"] = cfg.get("verbose", False)
    conf["input_encoding"] = cfg["file"].get("encoding", "utf-8")
    conf["output_format"] = cfg["file"].get("output_type", "html").lower()
    conf["save_as"] = "{0}.{1}".format(cfg["file"].get("save_as", str(time.time())), conf["output_format"])
    conf["crop_percentage"] = cfg["photos"].get("crop_percentage", 100)
    conf["k"] = cfg["results"].get("desired_analysis_clusters", 3)
    conf["cropped_folder"] = cfg["photos"].get("cropped_photo_dir", "cropped_photos").lower()
    conf["photo_destination_folder"] = cfg["photos"].get("base_photo_dir", "product_photos").lower()
    conf["save_product_photos"] = cfg["photos"].get("save_product_photos", False)
    conf["save_cropped_photos"] = cfg["photos"].get("save_cropped_photos", False)

    # things to be implemented later.
    conf["minimum_confidence"] = cfg["results"].get("confidence_minimum", 0)
    conf["debugging"] = cfg.get("debug", False)

    if conf["verbose"]:
        log_configuration_values(conf)

    return conf


def tidy_up(conf, folders):
    if conf["save_product_photos"]:
        logging.info("Removing downloaded files")
        folders.remove(conf["photo_destination_folder"])
    if conf["save_cropped_photos"]:
        logging.info("Removing cropped files")
        folders.remove(conf["cropped_folder"])
    for folder in folders:
        try:
            shutil.rmtree(folder)
            logging.info("Removed folder: {0}".format(folder))
        except OSError:
            logging.info("Unable to delete folder: {0}".format(folder))
    if not conf["verbose"]:
        os.remove("image_curator.log")


def main():
    conf = ensure_valid_configuration(get_config_from_file("config.json"))
    logging.info("Configuration file successfully loaded.")
    folders = [conf["photo_destination_folder"], conf["cropped_folder"]]
    for folder in folders:
        try:
            os.mkdir(folder)
            logging.info("Created folder: {0}".format(folder))
        except FileExistsError:
            pass

    if conf["app_mode"].lower() == "color":
        logging.info("Analyzing color")
        photos = photo_retriever.retrieve_photos_from_file(conf)
        logging.info("Collected {0} photos".format(len(photos)))
        analysis_results = color_analysis.analyze_image_colors(conf, photos)
        with open(conf["save_as"], "w") as output:
            html_output = build_results_page(analysis_results)
            output.write(html_output)
        logging.info("HTML file created.")
        tidy_up(conf, folders)
    elif conf["app_mode"].lower() == "shape":
        # TODO: implement this feature somewhere/how.
        pass
    else:
        raise ValueError("Application mode '{0}' invalid.".format(conf["app_mode"]))

if __name__ == "__main__":
    establish_logger()
    start_time = time.time()
    logging.info("Script started.")
    main()
    end_time = time.time()
    run_time = end_time - start_time
    logging.info("Script completed in {0}".format(run_time))
