import os
import csv
import json
import time
import shutil
import logging
import argparse
import color_analysis
import product_database
import analysis_objects
import html_results_page


def establish_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", dest="config_file", default="config.json",
                        help="Configuration file, in JSON format.")
    return parser.parse_args()


def establish_logger(arguments):
    conf = get_config_from_file(arguments.config_file)
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
        logging.info("Product photos will be saved in {0}.".format(conf["photo_destination_folder"]))
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
        conf["sku_column"] = cfg["file"]["sku_column"]
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
    conf["color_mode"] = cfg["colors"].get("mode", "RGB").upper()

    valid_color_modes = ["RGBA", "RGB", "HSL"]
    if conf["color_mode"].upper() == "RGBA":
        conf["color_mode"] = "RGB"
    if conf["color_mode"].upper() not in valid_color_modes:
        raise ValueError("Invalid color mode specified, the acceptable values are RGB and HSL.")

    # things to be implemented later.
    conf["minimum_confidence"] = cfg["results"].get("confidence_minimum", 0)
    conf["debugging"] = cfg.get("debug", False)

    if conf["verbose"]:
        log_configuration_values(conf)

    return conf


def retrieve_photos_from_file(conf):
    products = []
    with open(conf["source_file"], encoding=conf["input_encoding"]) as source:
        reader = csv.DictReader(source)
        for row in reader:
            product = analysis_objects.Product()
            product.sku = row[conf["sku_column"]]
            # product_database.write_sku_to_db(product.sku)
            product.photo_url = row[conf["photo_column"]]
            if not product.photo_url:
                continue
            product.retrieve_photo(conf["photo_destination_folder"])
            products.append(product)
    return products


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
        try:
            os.remove("image_curator.log")
        except FileNotFoundError:
            pass


def main():
    args = establish_arguments()
    establish_logger(args)
    start_time = time.time()
    logging.info("Script started.")

    conf = ensure_valid_configuration(get_config_from_file("config.json"))
    logging.info("Configuration file successfully loaded.")

    if conf["debugging"]:
        with open(conf["source_file"], encoding=conf["input_encoding"]) as source:
            reader = csv.DictReader(source)
            for row in reader:
                product = analysis_objects.Product()
                product.sku = row[conf["sku_column"]]
                product.title = row["title"]
                product.photo_url = row[conf["photo_column"]]
                product_database.write_product_to_db(product.__dict__)
    else:
        folders = [conf["photo_destination_folder"], conf["cropped_folder"]]
        for folder in folders:
            try:
                os.mkdir(folder)
                logging.info("Created folder: {0}".format(folder))
            except FileExistsError:
                pass

        if conf["app_mode"].lower() == "color":
            logging.info("Analyzing color")
            products = retrieve_photos_from_file(conf)
            logging.info("Collected {0} products".format(len(products)))
            # TODO: investigate HSL comparison vs. RGB comparison.
            analysis_results = color_analysis.analyze_image_colors(conf, products)
            with open(conf["save_as"], "w") as output:
                if conf["output_format"].lower() == "html":
                    html_output = html_results_page.builder(analysis_results)
                    output.write(html_output)
                    logging.info("HTML file created.")
                elif conf["output_format"].lower() == "csv":
                    # TODO: implement CSV output...again.
                    pass
        elif conf["app_mode"].lower() == "shape":
            raise ValueError("Shape analysis not yet supported.")
        else:
            raise ValueError("Application mode '{0}' invalid.".format(conf["app_mode"]))
        tidy_up(conf, folders)
    end_time = time.time()
    run_time = end_time - start_time
    logging.info("Script completed in {0}".format(run_time))

if __name__ == "__main__":
    main()
