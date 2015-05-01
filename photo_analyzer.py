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
    cropped_images_are_to_be_saved = config["photos"]["save_cropped_photos"]
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
        fieldnames.append("computed_strategy_colors")
        writer = csv.DictWriter(output, fieldnames)
        writer.writeheader()
        for row in reader:
            new_row = row
            for key, value in row.items():
                new_row.update({key: value})

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

            new_row["computed_strategy_colors"] = set(computed_strategy_colors)
            writer.writerow(row)

if __name__ == "__main__":
    main()