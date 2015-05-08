import logging


page_template = """
<html>
    <head>
        <style>
            html {{
                background-color: #222;
            }}

            .wrapper {{
                display: flex;
                width: {crop_width}px;
                margin: 0 auto;
                flex-flow: row wrap;
                font-weight: bold;
                text-align: center;
                padding: 1rem;
            }}

            .wrapper > * {{
                padding: 10px;
                flex: 1 100%;
            }}

            .swatch{{
                flex: 1 auto;
                height: 1.5rem;
                font-family: sans-serif;
                color: #FFFFFF;
                text-shadow: -1px 0 black, 0 1px black, 1px 0 black, 0 -1px black;
            }}
        </style>
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js"></script>
        <script>
            $(document).ready(function(){{
                $(".swatch").click(function(){{
                    var color = $(this).css("background-color");
                    $(this).remove();
                }});
            }});
        </script>
    </head>
    <body>
        {products}
    </body>
</html>
"""


def build_swatches(colors, color_relationships):
    swatches = ""
    for color in colors:
        matched_colors = ""
        if color in color_relationships:
            for match in color_relationships[color]:
                matched_colors += "{0}, ".format(match)
            matched_colors = matched_colors[:-2]
        swatches += '<div class="swatch" style="background: {0};" title="{0}">{1}</div>\n                   '.format(color, matched_colors)
    return swatches


def build_result(image_link, swatches):
    return """
           <div class="wrapper">
               <div style="border: 1px solid white;">
                   <img src="{0}">
               </div>
               {1}
           </div>
           """.format(image_link, swatches)


def build_results_page(analysis_results):
    logging.info("Begun results page construction")
    output = ""
    crop_width = analysis_results.pop("crop_width")
    for result in analysis_results:
        logging.info("Building elements for {0}".format(result))
        image = analysis_results[result]
        swatches = build_swatches(image["computed_colors"], image["color_relationship"])
        output += build_result(image["image_encoding"], swatches)
    return page_template.format(crop_width=crop_width, products=output)
