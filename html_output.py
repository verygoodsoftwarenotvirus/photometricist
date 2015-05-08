import photo_functions
import color_analysis
import result_page_builder


def output_html(conf, photos):
    with open(conf["save_as"], "w") as output:
        color_relationships = {}
        results = {}
        crop_widths = []
        for photo in photos:
            photo_path = "{0}{1}".format(conf["photo_destination_folder"],
                                         photo[photo.rfind("/"):])
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
            photo_path = photo_path.replace(conf["photo_destination_folder"], conf["cropped_folder"])
            b64_encoded = str(photo_functions.base64_encode_image(photo_path))[2:-1]
            photo_extension = photo_path[photo_path.rfind(".") + 1:]
            if photo_extension == "jpg":
                photo_extension = "jpeg"
            b64_prefix = "data:image/{0};base64,".format(photo_extension)
            photo_path = "{0}{1}".format(b64_prefix, b64_encoded)
            crop_widths.append(photo_functions.get_image_width(image))
            results[photo_path] = analysis_results
        html_output = result_page_builder.build_page(max(crop_widths) + 20,
                                                     results, color_relationships)
        output.write(html_output)
