import argparse
import time
import json


def get_config_from_file(file_location="config.json"):
    with open(file_location) as config_file:
        config = json.load(config_file)

    return config


def set_up_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("-d", "--data",
                        dest="source_file",
                        type=str,
                        default=None,
                        help="The file that contains your source data")

    parser.add_argument("-o", "--output",
                        dest="save_as",
                        type=str,
                        default=None,
                        help="The filename you want your results saved as")

    parser.add_argument("-i", "--photo_column",
                        dest="photo_link_column",
                        type=str,
                        default=None,
                        help="The column in your source CSV that contains" +
                             " the URL for that product's image.")

    parser.add_argument("-f", "--folder",
                        dest="photo_directory",
                        default=int(time.time()),
                        help="Where your saved photos will be stored.")

    return parser.parse_args()


if __name__ == "__main__":
    config = get_config_from_file("config.json")