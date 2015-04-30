from collections import namedtuple
from math import sqrt
import random
from webcolors import rgb_to_hex, hex_to_rgb

"""
    Modified from
     http://charlesleifer.com/blog/using-python-and-k-means-to-find-the-dominant-colors-in-images/
"""


Point = namedtuple('Point', ('coords', 'n', 'ct'))
Cluster = namedtuple('Cluster', ('points', 'center', 'n'))


def get_points(image):
    points = []
    w, h = image.size
    for count, color in image.getcolors(w * h):
        points.append(Point(color, 3, count))
    return points


def analyze_color(image, n=3):
    points = get_points(image)
    clusters = k_means(points, n, 1)
    colors = [map(int, c.center.coords) for c in clusters]

    result = []
    for color in colors:
        result.append(rgb_to_hex(color))
    return result


def euclidean(p1, p2):
    return sqrt(sum([(p1.coords[i] - p2.coords[i]) ** 2 for i in range(p1.n)]))


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