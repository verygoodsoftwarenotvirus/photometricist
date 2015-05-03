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
              height: 2.5rem;
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


def build_swatches(colors):
    swatches = ""
    for color in colors:
        swatches += '<div class="swatch" style="background: {};"></div>\n'.format(color)
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


def build_page(crop_width, analysis_results):
    result = ""
    for image in analysis_results:
        swatches = build_swatches(analysis_results[image])
        result += build_result(image, swatches)
    return page_template.format(crop_width=crop_width, products=result)
