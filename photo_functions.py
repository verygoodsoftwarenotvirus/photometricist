from PIL import Image


def open_image(filename):
    return Image.open(filename)


def save_image(image, filename):
    image.save(filename)


def center_crop_image_by_percentage(image, percentage=0):
    """
    TODO: handle different percentages

    :param image: a PIL Image object
    :param percentage: the percentage you'd like to crop away.
    :return: the cropped image object.
    """

    percentage = int(percentage)
    if percentage > 100:
        raise ValueError("you cannot crop an image to be larger than the original image.")

    modifier = 0
    if percentage:
        modifier = percentage * .01

    width, height = image.size
    left = abs(int(width - (width * modifier)))
    right = abs(int(height - (height * modifier)))
    top = abs(int(width * modifier))
    bottom = abs(int(height * modifier))

    crop_box = (left, right, top, bottom)
    return image.crop(crop_box)


def create_photo_thumbnail(image, size=200):
    return image.thumbnail((size, size))