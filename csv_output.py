import csv
import photo_retriever
import photo_functions
import color_analysis


def output_csv(conf):
    with open(conf["source_file"], encoding=conf["input_encoding"]) as source, \
            open(conf["save_as"], "w") as output:
        color_relationships = {}
        reader = csv.DictReader(source)
        photo_column = conf["photo_column"]
        fieldnames = []
        if "sku" in [column.lower() for column in reader.fieldnames]:
            fieldnames.append('sku')
        fieldnames += [photo_column, "computed_colors", "computed_strategy_colors", "confidence_index"]
        writer = csv.DictWriter(output, fieldnames, lineterminator='\n')
        writer.writeheader()
        for row in reader:
            photo_link = row[conf["photo_column"]]
            if not photo_link:
                continue
            photo_path = "{0}{1}".format(conf["photo_destination_folder"],
                                         photo_link[photo_link.rfind("/"):photo_link.rfind("?")])
            photo_retriever.retrieve_photo(photo_link, photo_path)
            image = photo_functions.crop_and_save_photo(photo_path, conf["crop_percentage"],
                                                        conf["photo_destination_folder"],
                                                        conf["cropped_folder"])
            analysis_results = color_analysis.analyze_color(image, conf["k"])
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
