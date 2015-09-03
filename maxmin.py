import json
import math
import statistics
from webcolors import rgb_to_hex, hex_to_rgb

colors = None
with open("colors.json") as colorfile:
    colors = json.load(colorfile)

for color in colors.items():
    hex_color = color[1]['color']
    variance = .01 * color[1]['variance']
    max_color = [0, 0, 0]
    min_color = [0, 0, 0]
    color_value = hex_to_rgb(hex_color)
    for i in range(len(color_value)):
        min_color[i] = statistics.median([0, math.floor(color_value[i] - (color_value[i] * variance)), 255])
        max_color[i] = statistics.median([0, math.ceil(color_value[i] - (color_value[i] * variance)), 255])

    print("\n{0}".format(color[0]))
    print(rgb_to_hex(min_color))
    print(hex_color)
    print(rgb_to_hex(max_color))

    print("\n{0}".format(color[0]))
    print(min_color)
    print(color_value)
    print(max_color)