import base64
import logging
from PIL import Image


def open_image(filename):
    logging.info("Opening image: {0}".format(filename))
    return Image.open(filename)


def save_image(image, filename):
    logging.info("Saving image: {0}".format(filename))
    image.save(filename)


def get_image_width(image):
    return image.size[0]


def get_image_height(image):
    return image.size[1]


def get_image_size(image):
    return image.size


def base64_encode_image(filename):
    with open(filename, "rb") as image:
        return base64.b64encode(image.read())


def create_photo_thumbnail(image, size=300):
    logging.info("Creating thumbnail {0} x {0}".format(size))
    return image.thumbnail((size, size))


def crop_and_save_photo(photo_path, crop_percentage, photo_folder, cropped_folder):
    logging.info("Cropping {0}".format(photo_path))
    image = open_image(photo_path)
    image = center_crop_image_by_percentage(image, crop_percentage)
    photo_path = photo_path.replace(photo_folder, cropped_folder)
    save_image(image, photo_path)
    logging.info("{0} saved".format(photo_path))
    return image


def center_crop_image_by_percentage(image, percentage=0):
    logging.info("Cropping image from the center by {1}%".format(image, percentage))
    percentage = int(min(percentage, 100))
    modifier = 0
    if percentage:
        modifier = percentage * .01

    width, height = image.size
    left = abs(int(width - (width * modifier)))
    right = abs(int(height - (height * modifier)))
    top = abs(int(width * modifier))
    bottom = abs(int(height * modifier))

    return image.crop((left, right, top, bottom))
