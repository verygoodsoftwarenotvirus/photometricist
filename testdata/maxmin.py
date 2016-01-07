import json
import math
import statistics
import colour
import colorsys
from webcolors import rgb_to_hex

colors = None
with open("colordefs.json") as colorfile:
    colors = json.load(colorfile)

for color in colors.items():
    minHue = color[1]["minHue"]/360
    maxHue = color[1]["maxHue"]/360
    minSaturation = color[1]["minSaturation"]/100
    maxSaturation = color[1]["maxSaturation"]/100
    minLightness = color[1]["minLightness"]/100
    maxLightness = color[1]["maxLightness"]/100

    print("\n", color[0])

    minColor = colour.Color(hsl=(minHue, minSaturation, minLightness))
    maxColor = colour.Color(hsl=(maxHue, maxSaturation, maxLightness))

    minHSV = colorsys.rgb_to_hsv(minColor.rgb[0], minColor.rgb[1], minColor.rgb[2])
    maxHSV = colorsys.rgb_to_hsv(maxColor.rgb[0], maxColor.rgb[1], maxColor.rgb[2])

    # print("[{0}, {1}, {2}]".format(int(minHSV[0]*360), minHSV[1], minHSV[2]))
    # print("[{0}, {1}, {2}]".format(int(maxHSV[0]*360), maxHSV[1], maxHSV[2]))
