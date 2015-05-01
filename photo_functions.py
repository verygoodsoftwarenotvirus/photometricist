import base64
from PIL import Image


def open_image(filename):
    return Image.open(filename)


def save_image(image, filename):
    image.save(filename)


def base645_encode_image(filename):
    with open(filename, "rb") as image:
        return base64.b64encode(image.read())


def create_photo_thumbnail(image, size=300):
    return image.thumbnail((size, size))


def center_crop_image_by_percentage(image, percentage=0):
    """
    TODO: handle different percentages

    :param image: a PIL Image object
    :param percentage: the percentage you'd like to crop away.
    :return: the cropped image object.
    """

    percentage = int(min(percentage, 100))
    modifier = 0
    if percentage:
        modifier = percentage * .01

    width, height = image.size
    left = abs(int(width - (width * modifier)))
    right = abs(int(height - (height * modifier)))
    top = abs(int(width * modifier))
    bottom = abs(int(height * modifier))

    return image.crop((left,
                       right,
                       top,
                       bottom))
