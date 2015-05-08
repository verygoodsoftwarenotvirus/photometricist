import photo_functions
import color_analysis


def perform_image_color_analysis(conf, photos):
    color_relationships = {}
    for photo in photos:
        image = photo_functions.crop_and_save_photo(photo, conf["crop_percentage"],
                                                    conf["photo_destination_folder"],
                                                    conf["cropped_folder"])
        analysis_results = color_analysis.analyze_color(image, conf["k"])
        relationship_result = color_analysis.compute_color_matches(conf, analysis_results)
        computed_strategy_colors = []
        for relationship in relationship_result:
            for color in relationship_result[relationship]:
                computed_strategy_colors.append(color)
        color_relationships.update(relationship_result)
    return color_relationships