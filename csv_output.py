import csv


def establish_csv_headers(configuration):
    reader = csv.DictReader(configuration["file"]["source_file"])
    photo_column = configuration["file"].get("photo_info_column")
    headers = []
    if "sku" in [column.lower() for column in reader.fieldnames]:
        headers.append('sku')
    headers += [photo_column, "computed_colors", "computed_strategy_colors", "confidence_index"]
    return headers


def output_csv():
    try:
        with open(source_file, encoding=input_encoding) as source, open(save_as, "w") as output:
            reader = csv.DictReader(source)
            color_relationships = {}
            fieldnames = establish_csv_headers(conf, reader.fieldnames)
            writer = csv.DictWriter(output, fieldnames, lineterminator='\n')
            writer.writeheader()
            for row in reader:
                photo_link = row[photo_column]
                if not photo_link:
                    continue
                photo_path = "{0}{1}".format(photo_destination_folder,
                                             photo_link[photo_link.rfind("/"):photo_link.rfind("?")])
                photo_retriever.retrieve_photo(photo_link, photo_path)

                # crop and save image
                image = photo_functions.open_image(photo_path)
                image = photo_functions.center_crop_image_by_percentage(image, crop_percentage)
                photo_path = photo_path.replace(photo_destination_folder, cropped_folder)
                photo_functions.save_image(image, photo_path)

                analysis_results = color_analysis.analyze_color(image, k)
                relationship_result = color_analysis.compute_color_matches(conf, analysis_results)
                computed_strategy_colors = []
                for relationship in relationship_result:
                    for color in relationship_result[relationship]:
                        computed_strategy_colors.append(color)
                color_relationships.update(relationship_result)
                new_row = {}
                for key, value in row.items():
                    if key in fieldnames:
                        new_row[key] = value
                new_row["computed_colors"] = analysis_results
                computed_strategy_colors = set(computed_strategy_colors)
                if computed_strategy_colors:
                    new_row["computed_strategy_colors"] = computed_strategy_colors
                writer.writerow(new_row)
    except EnvironmentError:
        print("Something went awry.")