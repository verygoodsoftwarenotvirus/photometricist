from collections import namedtuple
from math import sqrt
import random
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
    return sqrt(sum([(point1.coords[i] - point2.coords[i]) ** 2 for i in range(point1.n)]))


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


def color_is_in_range(to_compare, color_range, margin_of_error=None):
    floor_color, ceiling_color = color_range
    to_compare = hex_to_rgb(to_compare)
    floor_color = hex_to_rgb(floor_color)
    ceiling_color = hex_to_rgb(ceiling_color)

    for color_value in to_compare:
        index = to_compare.index(color_value)
        floor_value = min(floor_color[index], ceiling_color[index])
        ceiling_value = max(floor_color[index], ceiling_color[index])
        if margin_of_error:
            margin_of_error = min(max(0, margin_of_error), 100)
            floor_value -= floor_value * (.01 * margin_of_error)
            ceiling_value += ceiling_value * (.01 * margin_of_error)
        if not floor_value <= color_value <= ceiling_value:
            return False
    return True


def compute_color_matches(config, results, margin_of_error=None, minimum_confidence=None):
    computed_colors = []
    color_relationships = {}
    for color in config["colors"]:
        color_floor = config["colors"][color]["floor"]
        color_ceiling = config["colors"][color]["ceiling"]
        color_range = (color_floor, color_ceiling)
        for result in results:
            if color_is_in_range(result, color_range, margin_of_error):
                if color in color_relationships:
                    color_relationships[result].append(color)
                else:
                    color_relationships[result] = [color]
                computed_colors.append(color)
    return computed_colors, color_relationships