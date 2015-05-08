import math
import random
import photo_functions
from collections import namedtuple
from webcolors import rgb_to_hex, hex_to_rgb

"""
    This code is either stolen and modified or inpired by:
     http://charlesleifer.com/blog/using-python-and-k-means-to-find-the-dominant-colors-in-images/

    ...depending on your perspective.
"""

Point = namedtuple('Point', ('coords', 'n', 'ct'))
Cluster = namedtuple('Cluster', ('points', 'center', 'n'))


def get_points(image):
    points = []
    w, h = image.size
    for count, color in image.getcolors(w * h):
        points.append(Point(color, 3, count))
    return points


def analyze_color(image, k=3):
    points = get_points(image)
    clusters = k_means(points, k, 1)
    colors = [map(int, c.center.coords) for c in clusters]

    result = []
    for color in colors:
        result.append(rgb_to_hex(color))
    return result


def euclidean(point1, point2):
    return math.sqrt(sum([(point1.coords[i] - point2.coords[i]) ** 2 for i in range(point1.n)]))


def calculate_center(points, n):
    values = [0.0 for i in range(n)]
    p_length = 0
    for p in points:
        p_length += p.ct
        for i in range(n):
            values[i] += (p.coords[i] * p.ct)
    return Point([(v / p_length) for v in values], n, 1)


def k_means(points, k, min_diff):
    clusters = [Cluster([p], p, p.n) for p in random.sample(points, k)]

    while 1:
        point_lists = [[] for i in range(k)]
        for p in points:
            smallest_distance = float('Inf')
            for i in range(k):
                distance = euclidean(p, clusters[i].center)
                if distance < smallest_distance:
                    smallest_distance = distance
                    idx = i
            point_lists[idx].append(p)

        diff = 0
        for i in range(k):
            old = clusters[i]
            center = calculate_center(point_lists[i], old.n)
            new = Cluster(point_lists[i], center, old.n)
            clusters[i] = new
            diff = max(diff, euclidean(old.center, new.center))

        if diff < min_diff:
            break
    return clusters

"""
    Everything beyond here is my own code
"""


def color_is_in_range(to_compare, floor_color, ceiling_color):
    to_compare = hex_to_rgb(to_compare)
    for x in range(len(to_compare)):
        if not floor_color[x] <= to_compare[x] <= ceiling_color[x]:
            return False
    return True


def compute_color_matches(config, results, minimum_confidence=None):
    color_relationships = {}
    for color in config["colors"]:
        percentage = .01 * config["colors"][color]["variance"]
        color_value = hex_to_rgb(config["colors"][color]["color"])
        floor_color = []
        ceiling_color = []
        for value in color_value:
            floor_color.append(math.floor(value - (value * percentage)))
            ceiling_color.append(math.ceil(value + (value * percentage)))
        for result in results:
            if color_is_in_range(result, floor_color, ceiling_color):
                if color in color_relationships:
                    color_relationships[result].append(color)
                else:
                    color_relationships[result] = [color]
    return color_relationships


def analyze_colors(conf, photos):
        results = {}
        crop_widths = []
        for photo_path in photos:
            crop_percentage = conf["crop_percentage"]
            photo_folder = conf["photo_destination_folder"]
            cropped_folder = conf["cropped_folder"]
            image = photo_functions.crop_and_save_photo(photo_path, crop_percentage,
                                                        photo_folder, cropped_folder)

            cropped_photo_path = photo_path.replace(conf["photo_destination_folder"], conf["cropped_folder"])
            computed_colors = analyze_color(image, conf["k"])
            color_relationship = compute_color_matches(conf, computed_colors)
            # encoding adds some junk on the sides.
            b64_encoded = str(photo_functions.base64_encode_image(cropped_photo_path))[2:-1]
            photo_extension = photo_path[photo_path.rfind(".") + 1:]
            if photo_extension == "jpg":
                photo_extension = "jpeg"
            b64_prefix = "data:image/{0};base64,".format(photo_extension)
            image_encoding = "{0}{1}".format(b64_prefix, b64_encoded)
            crop_widths.append(photo_functions.get_image_width(image))

            result = {"computed_colors": computed_colors,
                      "color_relationship": color_relationship,
                      "image_encoding": image_encoding}
            results[photo_path] = result
        results["crop_width"] = max(crop_widths) + 20
        return results
